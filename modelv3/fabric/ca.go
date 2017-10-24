package fabric

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"errors"
	"os"
	"log"
	ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	config "github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	client "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/identity"
	kvs "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/keyvaluestore"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	fabricCAClient "github.com/hyperledger/fabric-sdk-go/pkg/fabric-ca-client"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
)

var allFabricCAConfig config.Config

type CA struct{
	CaService *fabricCAClient.FabricCA
	CaConfig  *config.CAConfig
	Client    *client.Client
	AdminUser fab.User
	MspID     string
	OrgId     string
}


func InitCA(configFile string) {
	fmt.Println("Init CA ConfigFile:",configFile)
	testSetup := BaseSetupImpl{
		ConfigFile: configFile,
	}

	fabricCAConfig, err := testSetup.InitConfig()
	if err != nil {
		fmt.Printf("Failed InitConfig [%s]\n", err)
		os.Exit(1)
	}
	allFabricCAConfig = fabricCAConfig
}


func (c *CA)InitCaServer(orgId string, enroll_user_dir string) error {
	var err error
	c.OrgId = orgId
	c.MspID, err = allFabricCAConfig.MspID(orgId)
	if err != nil {
		log.Println("MspID() returned error: %v,%s", err, orgId)
	}
	log.Println("InitCaServer MspID:", c.OrgId)

	c.CaConfig, err = allFabricCAConfig.CAConfig(orgId)
	if err != nil {
		log.Println("GetCAConfig returned error: %s", err)
		return err
	}
	c.Client = client.NewClient(allFabricCAConfig)

	err = bccspFactory.InitFactories(allFabricCAConfig.CSPConfig())
	if err != nil {

		log.Println("Failed getting ephemeral software-based BCCSP [%s]", err)
		return err
	}

	cryptoSuite := bccspFactory.GetDefault()

	c.Client.SetCryptoSuite(cryptoSuite)
	//stateStore, err := kvs.CreateNewFileKeyValueStore("enroll_user") //Path: "enroll_user"
	stateStore, err := kvs.CreateNewFileKeyValueStore(enroll_user_dir)
	if err != nil {
		log.Println("CreateNewFileKeyValueStore return error[%s]", err)
		return err
	}
	c.Client.SetStateStore(stateStore)

	c.CaService, err = fabricCAClient.NewFabricCAClient(allFabricCAConfig, orgId)
	if err != nil {
		log.Println("NewFabricCAClient return error: %v", err)
		return err
	}

	// Admin user is used to register, enrol and revoke a test user
	c.AdminUser, err = c.Client.LoadUserFromStateStore("admin")
	if err != nil {
		log.Println("client.LoadUserFromStateStore return error: %v", err)
		return err
	}
	if c.AdminUser == nil {
		key, cert, err := c.CaService.Enroll("admin", "adminpw")
		if err != nil {
			log.Println("Enroll return error: %v", err)
			return err
		}
		if key == nil {
			log.Println("private key return from Enroll is nil")
			return err
		}
		if cert == nil {
			log.Println("cert return from Enroll is nil")
			return errors.New("cert return from Enroll is nil")
		}

		certPem, _ := pem.Decode(cert)
		if err != nil {
			log.Println("pem Decode return error: %v", err)
			return err
		}

		cert509, err := x509.ParseCertificate(certPem.Bytes)
		if err != nil {
			log.Println("x509 ParseCertificate return error: %v", err)
			return err
		}
		if cert509.Subject.CommonName != "admin" {
			log.Println("CommonName in x509 cert is not the enrollmentID")
			return errors.New("CommonName in x509 cert is not the enrollmentID")
		}
		adminUser2 := identity.NewUser("admin", c.MspID)
		log.Println("InitCaServer Save admin MspID:", c.MspID)
		adminUser2.SetPrivateKey(key)
		adminUser2.SetEnrollmentCertificate(cert)
		err = c.Client.SaveUserToStateStore(adminUser2, false)
		if err != nil {
			log.Println("client.SaveUserToStateStore return error: %v", err)
			return err
		}
		c.AdminUser, err = c.Client.LoadUserFromStateStore("admin")
		if err != nil {
			log.Println("client.LoadUserFromStateStore return error: %v", err)
			return err
		}
		if c.AdminUser == nil {
			log.Println("client.LoadUserFromStateStore return nil")
			return errors.New("client.LoadUserFromStateStore return nil")
		}
	}
	return nil
}



// affiliation string
func (c *CA)RegisterAndEnrollUser(appid string, appkey string, affiliation string)  ([]byte, []byte, error) {
	if appid == "" || appkey == "" ||  c.CaConfig == nil ||  c.CaConfig.Name == "" {
		return nil, nil, errors.New("Parameter can not be empty")
	}
	// Register a random user
	registerRequest := ca.RegistrationRequest{
		Name:        appid,
		Secret:      appkey,
		Type:        "user",
		//Affiliation: "org2.department1"
		//Affiliation: "org1.department1",
		Affiliation: affiliation,
		CAName:      c.CaConfig.Name,
	}
	enrolmentSecret, err := c.CaService.Register(c.AdminUser, &registerRequest)
	if err != nil {
		log.Fatalf("Error from Register: %s", err)
		return nil, nil, err
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", appid, enrolmentSecret)

	// Enrol the previously registered user
	ekey, ecert, err := c.CaService.Enroll(appid, enrolmentSecret)
	if err != nil {
		log.Fatalf("Error enroling user: %s", err.Error())
		return nil, nil, err
	}
	//enroll
	fmt.Printf("** Attempt to enrolled user:  '%s'\n", appid)
	//create new user object and set certificate and private key of the previously enrolled user
	enrolleduser := identity.NewUser(appid, c.MspID)
	enrolleduser.SetEnrollmentCertificate(ecert)
	enrolleduser.SetPrivateKey(ekey)

	err = c.Client.SaveUserToStateStore(enrolleduser, false)
	if err != nil {
		log.Fatalf("client.SaveUserToStateStore return error: %v", err)
		return nil, nil, err
	}

	return ekey.SKI(), ecert, nil
}


func (c *CA)Register(appid string, appkey string, affiliation string) (string, error){

	if appid == "" || appkey == "" ||  c.CaConfig == nil ||  c.CaConfig.Name == "" {
		return "", errors.New("Parameter can not be empty")
	}
	// Register a random user
	registerRequest := ca.RegistrationRequest{
		Name:        appid,
		Secret:      appkey,
		Type:        "user",
		//Affiliation: "org2.department1"
		//Affiliation: "org1.department1",
		Affiliation: affiliation,
		CAName:      c.CaConfig.Name,
	}
	enrolmentSecret, err := c.CaService.Register(c.AdminUser, &registerRequest)
	if err != nil {
		log.Fatalf("Error from Register: %s", err)
		return "", err
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", appid, enrolmentSecret)
	return enrolmentSecret, nil
}

func (c *CA)EnrollUser(appid string, enrolmentSecret string) ([]byte, []byte, error){
	// Enrol the previously registered user
	ekey, ecert, err := c.CaService.Enroll(appid, enrolmentSecret)
	if err != nil {
		log.Fatalf("Error enroling user: %s", err.Error())
		return nil, nil, err
	}
	//enroll
	fmt.Printf("** Attempt to enrolled user:  '%s'\n", appid)
	//create new user object and set certificate and private key of the previously enrolled user
	enrolleduser := identity.NewUser(appid, c.MspID)
	enrolleduser.SetEnrollmentCertificate(ecert)
	enrolleduser.SetPrivateKey(ekey)

	err = c.Client.SaveUserToStateStore(enrolleduser, false)
	if err != nil {
		log.Fatalf("client.SaveUserToStateStore return error: %v", err)
		return nil, nil, err
	}

	return ekey.SKI(), ecert, nil
}



func (c *CA)RegisterClient(appid string, appkey string, affiliation string) (string, error){

	if appid == "" || appkey == "" ||  c.CaConfig == nil ||  c.CaConfig.Name == "" {
		return "", errors.New("Parameter can not be empty")
	}
	// Register a random user
	registerRequest := ca.RegistrationRequest{
		Name:        appid,
		Secret:      appkey,
		Type:        "user",//"client", //"user" // 修改为clietn或者user不影响
		//Affiliation: "org2.department1"
		//Affiliation: "org1.department1",
		Attributes:  []ca.Attribute{{"hf.Registrar.Roles","client,user,peer,validator,auditor"}}, //只有client权限的用户，才可以作为amdin用户给洽谈用户申请证书
		Affiliation: affiliation,

		CAName:      c.CaConfig.Name,
	}
	enrolmentSecret, err := c.CaService.Register(c.AdminUser, &registerRequest)
	if err != nil {
		log.Fatalf("Error from Register: %s", err)
		return "", err
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", appid, enrolmentSecret)
	return enrolmentSecret, nil
}


func (c *CA)InitCaServerOtherUser(preuser string, orgId string, enroll_user_dir string) error {
	var err error
	c.OrgId = orgId
	c.MspID, err = allFabricCAConfig.MspID(orgId)
	if err != nil {
		log.Println("MspID() returned error: %v,%s", err, orgId)
	}
	log.Println("InitCaServer MspID:", c.OrgId)

	c.CaConfig, err = allFabricCAConfig.CAConfig(orgId)
	if err != nil {
		log.Println("GetCAConfig returned error: %s", err)
		return err
	}
	c.Client = client.NewClient(allFabricCAConfig)

	err = bccspFactory.InitFactories(allFabricCAConfig.CSPConfig())
	if err != nil {

		log.Println("Failed getting ephemeral software-based BCCSP [%s]", err)
		return err
	}

	cryptoSuite := bccspFactory.GetDefault()

	c.Client.SetCryptoSuite(cryptoSuite)
	//stateStore, err := kvs.CreateNewFileKeyValueStore("enroll_user") //Path: "enroll_user"
	stateStore, err := kvs.CreateNewFileKeyValueStore(enroll_user_dir)
	if err != nil {
		log.Println("CreateNewFileKeyValueStore return error[%s]", err)
		return err
	}
	c.Client.SetStateStore(stateStore)

	c.CaService, err = fabricCAClient.NewFabricCAClient(allFabricCAConfig, orgId)
	if err != nil {
		log.Println("NewFabricCAClient return error: %v", err)
		return err
	}

	// Admin user is used to register, enrol and revoke a test user
	c.AdminUser, err = c.Client.LoadUserFromStateStore(preuser)
	if err != nil {
		log.Println("client.LoadUserFromStateStore return error: %v", err)
		return err
	}

	return nil
}


