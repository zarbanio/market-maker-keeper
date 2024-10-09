package run

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"github.com/zarbanio/market-maker-keeper/configs"
	blockptr "github.com/zarbanio/market-maker-keeper/internal/block-ptr"
	"github.com/zarbanio/market-maker-keeper/internal/chain"
	"github.com/zarbanio/market-maker-keeper/internal/dextrader"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/pair"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/erc20"
	"github.com/zarbanio/market-maker-keeper/internal/executor"
	"github.com/zarbanio/market-maker-keeper/internal/keystore"
	"github.com/zarbanio/market-maker-keeper/internal/mstore"
	"github.com/zarbanio/market-maker-keeper/internal/nobitex"
	"github.com/zarbanio/market-maker-keeper/internal/strategy"
	"github.com/zarbanio/market-maker-keeper/internal/uniswapv3"
	"github.com/zarbanio/market-maker-keeper/store"
)

func main(cfg configs.Config) {
	lvl, err := zerolog.ParseLevel(cfg.General.LogLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Level(lvl)

	postgresStore := store.NewPostgres(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB)
	err = postgresStore.Migrate(cfg.Postgres.MigrationsPath)
	if err != nil {
		log.Panic(err)
	}
	logger.Info().Msg("database migrated")

	blockPtr := blockptr.NewDBBlockPointer(postgresStore, cfg.Indexer.StartBlock)
	if !blockPtr.Exists() {
		logger.Debug().Msg("block pointer doest not exits. creating a new one")
		err := blockPtr.Create()
		if err != nil {
			logger.Fatal().Err(err).Msg("error creating block pointer")

		}
		logger.Debug().Uint64("start block", cfg.Indexer.StartBlock).Msg("new block pointer created.")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		logger.Panic().Msg("PRIVATE_KEY environment variable is not set")
	}

	executorWallet, err := keystore.New(privateKey)
	if err != nil {
		logger.Panic().Err(err).Msg("error while initializing new executor wallet")
	}
	eth, err := ethclient.Dial(cfg.Chain.Url)
	if err != nil {
		logger.Panic().Err(err).Msg("error while dialing eth client")
	}

	tokenStore := mstore.NewMemoryTokenStore()

	indexer := chain.NewIndexer(eth, chain.NewBlockCache(eth), cfg.Chain.BlockInterval, blockPtr, nil, nil)
	dexTrader := dextrader.New(
		eth,
		common.HexToAddress(cfg.Contracts.DexTrader),
		executorWallet,
		tokenStore,
	)

	for _, t := range cfg.Tokens {
		sym, err := symbol.FromString(t.Symbol)
		if err != nil {
			logger.Panic().Err(err).Msg("error while converting symbol type")
		}
		token := erc20.NewToken(common.HexToAddress(t.Address), sym, int64(t.Decimals))
		erc20Client := erc20.New(eth, token)
		err = tokenStore.AddToken(token)
		if err != nil {
			logger.Panic().Err(err).Msg("error while adding new token in token store")
		}
		dexTrader.AddERC20Client(erc20Client)
	}

	uniswapV3Factory := uniswapv3.NewFactory(eth, common.HexToAddress(cfg.Contracts.UniswapV3Factory))

	dai, err := tokenStore.GetTokenBySymbol(symbol.DAI)
	if err != nil {
		logger.Panic().Err(err).Msg("error while getting token by symbol")
	}
	zar, err := tokenStore.GetTokenBySymbol(symbol.ZAR)
	if err != nil {
		logger.Panic().Err(err).Msg("error while getting token by symbol")
	}

	// crate pair in database if not exist
	botPair := pair.Pair{QuoteAsset: dai.Symbol(), BaseAsset: zar.Symbol()}
	pairId, err := postgresStore.CreatePairIfNotExist(context.Background(), &botPair)
	if err != nil {
		logger.Panic().Err(err).Msg("error while creating pair")
	}
	botPair.Id = pairId

	poolFee := domain.ParseUniswapFee(cfg.Uniswap.PoolFee)
	_, err = uniswapV3Factory.GetPool(context.Background(), dai.Address(), zar.Address(), poolFee)
	if err != nil {
		logger.Panic().Err(err).Msg("error while getting pool from uniswapV3")
	}
	quoter := uniswapv3.NewQuoter(eth, common.HexToAddress(cfg.Contracts.UniswapV3Quoter))

	exchangeMarkets := make(map[string]symbol.Symbol)
	exchangeMarkets["USDT"] = symbol.USDT
	exchangeMarkets["IRT"] = symbol.IRT

	dexMarkets := make(map[string]symbol.Symbol)
	dexMarkets["DAI"] = symbol.DAI
	dexMarkets["ZAR"] = symbol.ZAR

	var nobitexExchange domain.Exchange
	nobitexExchange = nobitex.New(
		cfg.Nobitex.Url,
		cfg.Nobitex.Key,
		cfg.Nobitex.Timeout,
		decimal.NewFromInt(cfg.Nobitex.MinimumOrderToman),
		cfg.Nobitex.OrderStatusInterval,
	)

	if cfg.General.Environment == configs.EnvironmentTestnet {
		nobitexExchange = nobitex.NewMockExchange(
			cfg.Nobitex.Url,
			cfg.Nobitex.Timeout,
			decimal.NewFromInt(cfg.Nobitex.MinimumOrderToman),
			cfg.Nobitex.OrderStatusInterval,
			30,
		)
		nobitexExchange.SetBalance(symbol.TMN, decimal.NewFromInt(60000000))
		nobitexExchange.SetBalance(symbol.USDT, decimal.NewFromInt(1000))
		nobitexExchange.SetBalance(symbol.IRT, decimal.NewFromInt(60000000))
	}

	tokens := make(map[symbol.Symbol]domain.Token)
	tokens[symbol.DAI] = dai
	tokens[symbol.ZAR] = zar

	strategyConfig := strategy.Config{
		StartQty:        decimal.NewFromFloat(cfg.MarketMaker.StartQty),
		StepQty:         decimal.NewFromFloat(cfg.MarketMaker.StepQty),
		ProfitThreshold: decimal.NewFromInt(cfg.MarketMaker.ProfitThreshold),
		Slippage:        decimal.NewFromFloat(cfg.MarketMaker.Slippage),
	}

	logger.Info().
		Str("startQty", strategyConfig.StartQty.String()).
		Str("stepQty", strategyConfig.StepQty.String()).
		Str("profitThreshold", strategyConfig.ProfitThreshold.String()).
		Str("slippage", strategyConfig.Slippage.String()).
		Msg("market maker started")

	buyDaiInUniswapSellTetherInNobitex := &strategy.BuyDaiUniswapSellTetherNobitex{
		Store:      postgresStore,
		Nobitex:    nobitexExchange,
		DexQuoter:  quoter,
		DexTrader:  dexTrader,
		Tokens:     tokens,
		UniswapFee: domain.UniswapFeeFee01,
		Marketsdata: map[strategy.Market]strategy.MarketData{
			strategy.Nobitex:   strategy.NewMarketData(),
			strategy.UniswapV3: strategy.NewMarketData(),
		},
		Config: strategyConfig,
		Logger: logger,
	}
	buyTetherInNobitexSellDaiInUniswap := &strategy.SellDaiUniswapBuyTetherNobitex{
		Store:         postgresStore,
		Nobitex:       nobitexExchange,
		UniswapQuoter: quoter,
		DexTrader:     dexTrader,
		Tokens:        tokens,
		UniswapFee:    domain.UniswapFeeFee01,
		Marketsdata: map[strategy.Market]strategy.MarketData{
			strategy.Nobitex:   strategy.NewMarketData(),
			strategy.UniswapV3: strategy.NewMarketData(),
		},
		Config: strategyConfig,
		Logger: logger,
	}

	ctx := context.Background()

	strategies := []strategy.ArbitrageStrategy{buyTetherInNobitexSellDaiInUniswap, buyDaiInUniswapSellTetherInNobitex}

	exec := &executor.Executor{
		Store:                postgresStore,
		Strategies:           strategies,
		PairId:               pairId,
		Nobitex:              nobitexExchange,
		DexTrader:            *dexTrader,
		Indxer:               *indexer,
		Logger:               logger,
		NobitexRetryTimeOut:  cfg.Nobitex.RetryTimeOut,
		NobitexSleepDuration: cfg.Nobitex.RetrySleepDuration,
	}

	ticker := time.NewTicker(cfg.MarketMaker.Interval)
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			lastCycleId, err := postgresStore.GetLastCycleId(ctx)
			if err != nil {
				logger.Fatal().Err(err).Msg("error while getting last cycle Id")
			}
			cycleId := lastCycleId + 1
			err = postgresStore.CreateCycle(ctx, time.Now(), domain.CycleStatusRunning)
			if err != nil {
				logger.Fatal().Err(err).Msg("error while creating new cycle")
			}
			exec.SetCycleId(cycleId)
			logger.Info().Int64("cycleId", cycleId).Msg("new cycle started")
			exec.RunAll()
			status := domain.CycleStatusSuccess

			err = postgresStore.UpdateCycle(ctx, cycleId, time.Now(), status)
			if err != nil {
				logger.Fatal().Err(err).Msg("error while updating cycle")
			}
			logger.Info().Int64("cycleId", cycleId).Msg("cycle finished")
			select {
			case <-done:
				return
			case <-ticker.C:
				continue
			}
		}
	}()

	<-quit
	ticker.Stop()
}

func Register(root *cobra.Command) {
	root.PersistentFlags().String("config", "config.yaml", "read config file")
	root.AddCommand(
		&cobra.Command{
			Use:   "run",
			Short: "Run market maker bot",
			Run: func(cmd *cobra.Command, args []string) {
				configPath, _ := cmd.Flags().GetString("config")
				cfg := configs.ReadConfig(configPath)
				main(cfg)
			},
		},
	)
}
