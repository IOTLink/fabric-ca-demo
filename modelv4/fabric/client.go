package fabric

import (
	"encoding/json"
	"fmt"
	//"time"

	//"github.com/golang/protobuf/proto"
	//google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	//config "github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	//"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	//channel "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/identity"
	//fc "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/internal"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/internal/txnproc"
	//packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/packager"
	//peer "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/peer"

	"github.com/op/go-logging"

	"github.com/hyperledger/fabric/bccsp"
	//"github.com/hyperledger/fabric/common/crypto"
	//fcutils "github.com/hyperledger/fabric/common/util"
	//"github.com/hyperledger/fabric/protos/common"
	//pb "github.com/hyperledger/fabric/protos/peer"
	//protos_utils "github.com/hyperledger/fabric/protos/utils"
	//"log"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

// Client enables access to a Fabric network.
type Client struct {
	//channels    map[string]fab.Channel
	cryptoSuite bccsp.BCCSP
	stateStore  fab.KeyValueStore
	userContext fab.User
	config      Config//config.Config
}

// NewClient returns a Client instance.
func NewClient(config Config) *Client {
	//channels := make(map[string]fab.Channel)
	c := Client{cryptoSuite: nil, stateStore: nil, userContext: nil, config: config}
	return &c
}



// Config returns the configuration of the client.
func (c *Client) Config() Config {
	return c.config
}

// SetStateStore ...
/*
 * The SDK should have a built-in key value store implementation (suggest a file-based implementation to allow easy setup during
 * development). But production systems would want a store backed by database for more robust storage and clustering,
 * so that multiple app instances can share app state via the database (note that this doesn’t necessarily make the app stateful).
 * This API makes this pluggable so that different store implementations can be selected by the application.
 */
func (c *Client) SetStateStore(stateStore fab.KeyValueStore) {
	c.stateStore = stateStore
}

// StateStore is a convenience method for obtaining the state store object in use for this client.
func (c *Client) StateStore() fab.KeyValueStore {
	return c.stateStore
}

// SetCryptoSuite is a convenience method for obtaining the state store object in use for this client.
func (c *Client) SetCryptoSuite(cryptoSuite bccsp.BCCSP) {
	c.cryptoSuite = cryptoSuite
}

// CryptoSuite is a convenience method for obtaining the CryptoSuite object in use for this client.
func (c *Client) CryptoSuite() bccsp.BCCSP {
	return c.cryptoSuite
}

// SaveUserToStateStore ...
/*
 * Sets an instance of the User class as the security context of this client instance. This user’s credentials (ECert) will be
 * used to conduct transactions and queries with the blockchain network. Upon setting the user context, the SDK saves the object
 * in a persistence cache if the “state store” has been set on the Client instance. If no state store has been set,
 * this cache will not be established and the application is responsible for setting the user context again when the application
 * crashed and is recovered.
 */
func (c *Client) SaveUserToStateStore(user fab.User, skipPersistence bool) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	if user.Name() == "" {
		return fmt.Errorf("user name is empty")
	}
	c.userContext = user
	if !skipPersistence {
		if c.stateStore == nil {
			return fmt.Errorf("stateStore is nil")
		}
		userJSON := &identity.JSON{
			MspID:                 user.MspID(),
			Roles:                 user.Roles(),
			PrivateKeySKI:         user.PrivateKey().SKI(),
			EnrollmentCertificate: user.EnrollmentCertificate(),
		}
		data, err := json.Marshal(userJSON)
		if err != nil {
			return fmt.Errorf("Marshal json return error: %v", err)
		}
		err = c.stateStore.SetValue(user.Name(), data)
		if err != nil {
			return fmt.Errorf("stateStore SaveUserToStateStore return error: %v", err)
		}
	}
	return nil

}

// LoadUserFromStateStore ...
/**
 * Restore the state of this member from the key value store (if found).  If not found, do nothing.
 * @returns {Promise} A Promise for a {User} object upon successful restore, or if the user by the name
 * does not exist in the state store, returns null without rejecting the promise
 */
func (c *Client) LoadUserFromStateStore(name string) (fab.User, error) {
	if c.userContext != nil {
		return c.userContext, nil
	}
	if name == "" {
		return nil, nil
	}
	if c.stateStore == nil {
		return nil, nil
	}
	if c.cryptoSuite == nil {
		return nil, fmt.Errorf("cryptoSuite is nil")
	}
	value, err := c.stateStore.Value(name)
	if err != nil {
		//return nil, fmt.Errorf("name is not exist")
		return nil, nil
	}

	var userJSON identity.JSON
	err = json.Unmarshal(value, &userJSON)
	if err != nil {
		return nil, fmt.Errorf("stateStore GetValue return error: %v", err)
	}


	user := identity.NewUser(name, userJSON.MspID)
	user.SetRoles(userJSON.Roles)
	user.SetEnrollmentCertificate(userJSON.EnrollmentCertificate)
	fmt.Println("SPI----------->",string(userJSON.PrivateKeySKI))
	key, err := c.cryptoSuite.GetKey(userJSON.PrivateKeySKI)
	if err != nil {
		return nil, fmt.Errorf("cryptoSuite GetKey return error: %v", err)
	}

	user.SetPrivateKey(key)
	c.userContext = user
	//test
	/*
	ski := userJSON.PrivateKeySKI
	prikey := user.PrivateKey()
	keyi, err := json.Marshal(prikey)
	if err != nil {
		return nil, err
	}
	log.Println("LoadUserFromStateStore","name:",user.Name(),"prikey:",keyi, "SKI:", hex.EncodeToString(ski[:]))
	*/
	return c.userContext, nil
}


// UserContext returns the current User.
func (c *Client) UserContext() fab.User {
	return c.userContext
}

// SetUserContext ...
func (c *Client) SetUserContext(user fab.User) {
	c.userContext = user
}
