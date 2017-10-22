package fabric

import (
	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	//ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	//fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	//deffab "github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
)

// BaseSetupImpl implementation of BaseTestSetup
// BaseSetupImpl implementation of BaseTestSetup
type BaseSetupImpl struct {
	ConfigFile      string
}

// InitConfig ...
func (setup *BaseSetupImpl) InitConfig() (apiconfig.Config, error) {
	configImpl, err := config.InitConfig(setup.ConfigFile)
	if err != nil {
		return nil, err
	}
	return configImpl, nil
}


