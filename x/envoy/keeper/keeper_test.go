package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	marsapp "github.com/mars-protocol/hub/app"
	marsapptesting "github.com/mars-protocol/hub/app/testing"
)

func init() {
	ibctesting.DefaultTestingAppInit = SetupEnvoyTestingApp
}

func SetupEnvoyTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encCfg := marsapp.MakeEncodingConfig()
	app := marsapptesting.MakeSimpleMockApp()
	return app, marsapp.DefaultGenesisState(encCfg.Codec)
}

// KeeperTestSuite is a testing suite
type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
}
