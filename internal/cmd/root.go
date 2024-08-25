package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zarbanio/market-maker-keeper/internal/cmd/run"
)

func Execute() {
	root := &cobra.Command{
		Use:     "Market Maker",
		Short:   "A market market bot for  ZAR/DAI",
		Version: "0.1",
	}

	run.Register(root)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
