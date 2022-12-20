package main

import serverconfig "github.com/cosmos/cosmos-sdk/server/config"

// initAppConfig generates contents for `app.toml`. Take the default template
// and config, append custom parameters.

func initAppConfig() (string, interface{}) {
	template := serverconfig.DefaultConfigTemplate
	cfg := serverconfig.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	//   - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak
	//     their own app.toml config,
	//   - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//     own app.toml to override, or use this default value.
	cfg.MinGasPrices = "0umars"

	return template, cfg
}
