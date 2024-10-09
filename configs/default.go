package configs

import (
	"time"
)

func DefaultConfig() Config {
	return Config{
		General: General{
			Environment: EnvironmentMainnet,
			LogLevel:    "info",
		},
		MarketMaker: MarketMaker{
			StartQty:        1,
			StepQty:         1,
			ProfitThreshold: 5_000, // 50_000 TMN
			Interval:        time.Minute * 10,
			Slippage:        0.001,
		},
		Chain: Chain{
			BlockInterval: time.Millisecond * 500,
		},
		Tokens: []Token{
			{
				Address:  "0xd946188a614a0d9d0685a60f541bba1e8cc421ae",
				Decimals: 18,
				Symbol:   "ZAR",
			},
			{
				Address:  "0xda10009cbd5d07dd0cecc66161fc93d7c9000da1",
				Decimals: 18,
				Symbol:   "DAI",
			},
		},
		Uniswap: Uniswap{
			PoolFee: 0.01,
		},
		Nobitex: Nobitex{
			Url:                 "https://api.nobitex.ir",
			Key:                 "", // Assuming no default value for Key
			MinimumOrderToman:   300_000,
			Timeout:             time.Second * 60,  // 60s
			OrderStatusInterval: time.Second * 2,   // 2s
			RetryTimeOut:        time.Second * 360, // 360s
			RetrySleepDuration:  time.Second * 5,   // 5s
		},
		Contracts: Contracts{
			UniswapV3Factory: "0x1F98431c8aD98523631AE4a59f267346ea31F984",
			UniswapV3Quoter:  "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6",
		},
		Indexer: Indexer{
			StartBlock: 247010149,
		},
		Postgres: Postgres{
			Host:           "localhost",
			Port:           5432,
			User:           "postgres",
			Password:       "postgres",
			DB:             "postgres",
			MigrationsPath: "/migrations",
		},
	}
}
