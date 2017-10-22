package fabric

import (
	"strings"
	"fmt"
	"path"
	"github.com/spf13/viper"
	"github.com/op/go-logging"
	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"


	//"crypto/x509"
	//"encoding/pem"
	//"fmt"
	//"errors"
	"os"
	//"log"
	//ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	//config "github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	//client "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/identity"
	//kvs "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/keyvaluestore"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	//fabricCAClient "github.com/hyperledger/fabric-sdk-go/pkg/fabric-ca-client"
	//fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"

)

type DefaultConfig struct {
	TcertBatch int `yaml:"tcertbatch"`
	LoggingLevel string `yaml:"logginglevel"`
	KeystorePath string `yaml:"keystorepath"`
}


type SecurityConfig struct {
	Enabled         bool    `yaml:"enabled"`
	HashAlgorithm   string  `yaml:"hashAlgorithm"`
	Level           int     `yaml:"level"`
}

type Organizations struct {
	MspID         string  `yaml:"mspID"`
	TlsEnabled    bool     `yaml:"hashAlgorithm"`
	CaName        string  `yaml:"caname"`
	ServerURL     string  `yaml:"serverURL"`
	Tlscertfiles  string  `yaml:"tlscertfiles"`
	Tlskeyfile    string  `yaml:"tlskeyfile"`
	Tlscertfile   string  `yaml:"tlscertfile"`
}


// Config represents the configuration for the client
type Config struct {
	defaultConfig     *DefaultConfig
	securityConfig    map[string]SecurityConfig
	organizations     map[string]Organizations
}


var myViper = viper.New()
var mylog = logging.MustGetLogger("fabric_sdk_go")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} [%{module}] %{level:.4s} : %{color:reset} %{message}`,
)

// BaseSetupImpl implementation of BaseTestSetup
type BaseSetupImpl struct {
	ConfigFile      string
}

// InitConfig ...
func (setup *BaseSetupImpl) InitConfig() (*Config, error) {
	configImpl, err := MyInitConfig(setup.ConfigFile)
	if err != nil {
		return nil, err
	}
	return configImpl, nil
}


func MyInitConfig(configFile string) (*Config, error) {
	return InitConfigWithCmdRoot(configFile, "mydefine api")
}

// InitConfigWithCmdRoot reads in a config file and allows the
// environment variable prefixed to be specified
func InitConfigWithCmdRoot(configFile string, cmdRootPrefix string) (*Config, error) {
	myViper.SetEnvPrefix(cmdRootPrefix)
	myViper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	myViper.SetEnvKeyReplacer(replacer)

	if configFile != "" {
		// create new viper
		myViper.SetConfigFile(configFile)
		// If a config file is found, read it in.
		err := myViper.ReadInConfig()

		if err == nil {
			mylog.Infof("Using config file: %s", myViper.ConfigFileUsed())
		} else {
			return nil, fmt.Errorf("Fatal error config file: %v", err)
		}
	}

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	defaultConfig := new(DefaultConfig)
	defaultConfig.TcertBatch   = myViper.GetInt("default.tcertbatch")
	defaultConfig.LoggingLevel = myViper.GetString("default.logginglevel")
	defaultConfig.KeystorePath = myViper.GetString("default.keystorepath")

	loggingLevelString := defaultConfig.LoggingLevel
	logLevel := logging.INFO
	if loggingLevelString != "" {
		mylog.Infof("fabric_sdk_go Logging level: %v", loggingLevelString)
		var err error
		logLevel, err = logging.LogLevel(loggingLevelString)
		if err != nil {
			panic(err)
		}
	}
	logging.SetBackend(backendFormatter).SetLevel(logging.Level(logLevel), "mydefine api")

	var securityConfig map[string]SecurityConfig
	securityConfig = make(map[string]SecurityConfig)
	securityOpt := myViper.InConfig("security")
	if securityOpt == true {
		configMap := myViper.GetStringMap("security")
		for key, value := range configMap {
			//fmt.Println("key:", key, "value:", value)
			var security SecurityConfig
			m := value.(map[string]interface{})
			for key1, value1 := range m {
				fmt.Println("key1:", key1, "value1:", value1)
				if strings.ToLower(key1) == "enabled" {
					security.Enabled = value1.(bool)
				} else if strings.ToLower(key1) == "hashalgorithm" {
					security.HashAlgorithm = value1.(string)
				} else if strings.ToLower(key1) == "level" {
					security.Level = value1.(int)
				}
			}
			securityConfig[key] = security
		}
	}
	//fmt.Println(securityConfig)

	var orgConfig map[string]Organizations
	orgConfig = make(map[string]Organizations)
	orgConfigOpt := myViper.InConfig("organizations")
	if orgConfigOpt == true {
		configMap := myViper.GetStringMap("organizations")
		for key, value := range configMap {
			//fmt.Println("key:", key, "value:", value)
			var org Organizations
			m := value.(map[string]interface{})
			for key1, value1 := range m {
				//fmt.Println("key1:", key1, "value1:", value1)
				if strings.ToLower(key1) == "mspid" {
					org.MspID = value1.(string)
				} else if strings.ToLower(key1) == "tlsenabled" {
					org.TlsEnabled = value1.(bool)
				} else if strings.ToLower(key1) == "name" {
					org.CaName = value1.(string)
				} else if strings.ToLower(key1) == "serverurl" {
					org.ServerURL = value1.(string)
				} else if strings.ToLower(key1) == "tlscertfiles" {
					org.Tlscertfiles = value1.(string)
				} else if strings.ToLower(key1) == "tlskeyfile" {
					org.Tlskeyfile = value1.(string)
				} else if strings.ToLower(key1) == "tlscertfile" {
					org.Tlscertfile = value1.(string)
				}
			}
			orgConfig[key] = org
		}
	}
	//fmt.Println(orgConfig)
	return &Config{defaultConfig,securityConfig,orgConfig}, nil
}

// MspID returns the MSP ID for the requested organization
func (c *Config) MspID(org string) (string, error) {
	if c == nil || c.organizations == nil {
		return "", fmt.Errorf("param is abort")
	}
	mspID := c.organizations[org].MspID
	if mspID == "" {
		fmt.Println(c.organizations)
		return "", fmt.Errorf("MSP ID is empty for org: %s", org)
	}
	return mspID, nil
}

// CAConfig returns the CA configuration.
func (c *Config) CAConfig(org string) (*apiconfig.CAConfig, error) {
	if c == nil || c.organizations == nil {
		return nil, fmt.Errorf("param is abort")
	}
	config := c.organizations[org]

	var tls apiconfig.MutualTLSConfig //{config.Tlscertfiles,{config.Tlskeyfile,config.Tlscertfile}}}
    tls.Certfiles = config.Tlscertfiles
    tls.Client.Keyfile = config.Tlskeyfile
    tls.Client.Certfile = config.Tlscertfile
	caConfig := apiconfig.CAConfig{config.TlsEnabled,config.CaName, config.ServerURL,tls}
	return &caConfig, nil
}

// CAServerCertFiles Read configuration option for the server certificate files
func (c *Config) CAServerCertFiles(org string) ([]string, error) {
	if c == nil || c.organizations == nil {
		return nil, fmt.Errorf("param is abort")
	}
	config := c.organizations[org]
	certFiles := strings.Split(config.Tlscertfiles, ",")
	certFileModPath := make([]string, len(certFiles))
	for i, v := range certFiles {
		certFileModPath[i] = strings.Replace(v, "$GOPATH", os.Getenv("GOPATH"), -1)
	}
	return certFileModPath, nil
}

// CAClientKeyFile Read configuration option for the fabric CA client key file
func (c *Config) CAClientKeyFile(org string) (string, error) {
	if c == nil || c.organizations == nil {
		return "", fmt.Errorf("param is abort")
	}
	config := c.organizations[org]
	return strings.Replace(config.Tlskeyfile,"$GOPATH", os.Getenv("GOPATH"), -1), nil
}

// CAClientCertFile Read configuration option for the fabric CA client cert file
func (c *Config) CAClientCertFile(org string) (string, error) {
	if c == nil || c.organizations == nil {
		return "", fmt.Errorf("param is abort")
	}
	config := c.organizations[org]
	return strings.Replace(config.Tlscertfile,"$GOPATH", os.Getenv("GOPATH"), -1), nil
}


// SecurityAlgorithm ...
func (c *Config) SecurityAlgorithm(crypto string) (string, error) {
	if c == nil || c.securityConfig == nil {
		return "", fmt.Errorf("param is abort")
	}
	return c.securityConfig[crypto].HashAlgorithm,nil
}

// SecurityLevel ...
func (c *Config) SecurityLevel(crypto string) (int,error) {
	if c == nil || c.securityConfig == nil {
		return 0, fmt.Errorf("param is abort")
	}
	return c.securityConfig[crypto].Level,nil
}

// IsSecurityEnabled ...
func (c *Config) IsSecurityEnabled(crypto string) (bool,error) {
	if c == nil || c.securityConfig == nil {
		return false, fmt.Errorf("param is abort")
	}
	return c.securityConfig[crypto].Enabled, nil
}


// KeyStorePath returns the keystore path used by BCCSP
func (c *Config) KeyStorePath() (string , error) {
	if c == nil || c.defaultConfig == nil {
		return "", fmt.Errorf("param is abort")
	}
	return path.Join(c.defaultConfig.KeystorePath, "keystore"),nil
	//return c.defaultConfig.KeystorePath,nil
}

// CAKeyStorePath returns the same path as KeyStorePath() without the
// 'keystore' directory added. This is done because the fabric-ca-client
// adds this to the path
func (c *Config) CAKeyStorePath() string {
	return c.defaultConfig.KeystorePath
}


// TcertBatchSize ...
func (c *Config) TcertBatchSize() (int,error) {
	if c == nil || c.defaultConfig == nil {
		return 0, fmt.Errorf("param is abort")
	}
	return c.defaultConfig.TcertBatch,nil
}

// SecurityLevel ...
func (c *Config) DefaultConfigLevel() (string, error) {
	if c == nil || c.defaultConfig == nil {
		return "", fmt.Errorf("param is abort")
	}
	return c.defaultConfig.LoggingLevel, nil
}





// CSPConfig ...
func (c *Config) CSPConfig(crypto string) *bccspFactory.FactoryOpts {
	var err error
	algorHash, err := c.SecurityAlgorithm(crypto)
	if err != nil {
		mylog.Errorf("%s",err.Error())
		return nil
	}

	securityLevel, err :=  c.SecurityLevel(crypto)
	if err != nil {
		mylog.Errorf("%s",err.Error())
		return nil
	}

	keyStorePath, err := c.KeyStorePath()
	if err != nil {
		mylog.Errorf("%s",err.Error())
		return nil
	}

	return &bccspFactory.FactoryOpts{
		ProviderName: "SW",
		SwOpts: &bccspFactory.SwOpts{
			HashFamily: algorHash,
			SecLevel:  securityLevel,
			FileKeystore: &bccspFactory.FileKeystoreOpts{
				KeyStorePath: keyStorePath,
			},
			Ephemeral: false,
		},
	}
}
