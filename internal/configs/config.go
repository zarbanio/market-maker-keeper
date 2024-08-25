package configs

import (
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	General struct {
		Environment string `yaml:"Environment"`
		LogLevel    string `yaml:"LogLevel"`
	} `yaml:"General"`
	MarketMaker struct {
		StartQty        float64       `yaml:"StartQty"`
		StepQty         float64       `yaml:"StepQty"`
		EndQty          int64         `yaml:"EndQty"`
		ProfitThreshold int64         `yaml:"ProfitThreshold"`
		Interval        time.Duration `yaml:"Interval"`
		Slippage        float64       `yaml:"Slippage"`
	} `yaml:"MarketMaker"`
	Chain struct {
		Url           string        `yaml:"Url"`
		BlockInterval time.Duration `yaml:"BlockInterval"`
	} `yaml:"Chain"`
	Tokens []struct {
		Address  string `yaml:"Address"`
		Decimals int    `yaml:"Decimals"`
		Symbol   string `yaml:"Symbol"`
	} `yaml:"Tokens"`
	Uniswap struct {
		PoolFee float64 `yaml:"PoolFee"`
	}
	Nobitex struct {
		Url                 string        `yaml:"Url"`
		Key                 string        `yaml:"Key"`
		MinimumOrderToman   int64         `yaml:"MinimumOrderToman"`
		Timeout             time.Duration `yaml:"Timeout"`
		OrderStatusInterval time.Duration `yaml:"OrderStatusInterval"`
		RetryTimeOut        time.Duration `yaml:"RetryTimeOut"`
		RetrySleepDuration  time.Duration `yml:"RetrySleepDuration"`
	} `yaml:"nobitex"`
	Contracts struct {
		DexTrader        string `yaml:"DexTrader"`
		UniswapV3Factory string `yaml:"UniswapV3Factory"`
		UniswapV3Quoter  string `yaml:"UniswapV3Quoter"`
	} `yaml:"Contracts"`
	Indexer struct {
		StartBlock uint64 `yaml:"StartBlock"`
	}
	Postgres struct {
		Host           string `yaml:"Host"`
		Port           int    `yaml:"Port"`
		User           string `yaml:"User"`
		Password       string `yaml:"Password"`
		DB             string `yaml:"DB"`
		MigrationsPath string `yaml:"MigrationsPath"`
	} `yaml:"Postgres"`
}

func ReadConfig(configFile string) Config {
	c := &Config{}
	err := c.Unmarshal(c, configFile)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return *c
}

func (c *Config) Unmarshal(rawVal interface{}, fileName string) error {
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	var input interface{} = viper.AllSettings()
	config := defaultDecoderConfig(rawVal)
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	return c
}
