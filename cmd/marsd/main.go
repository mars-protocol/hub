package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	marsapp "github.com/mars-protocol/hub/app"
)

func main() {
	setAddressPrefixes(marsapp.AccountAddressPrefix)
	rootCmd := NewRootCmd(marsapp.MakeEncodingConfig())
	if err := svrcmd.Execute(rootCmd, marsapp.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
