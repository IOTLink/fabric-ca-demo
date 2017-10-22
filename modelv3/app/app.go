package main

import (
. "fabric-ca-demo/modelv3/fabric"
"log"
"fmt"
)


func ca_org1() {
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("peerorg1", "enroll_user_peerorg1")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	userNameA := GenerateRandomID()
	_,_,err = myca.RegisterAndEnrollUser(userNameA,"userAW", "org1.department1")
	if err != nil {
		fmt.Println("RegisterAndEnrollUser FAILT",err)
		return
	} else {
		fmt.Println("RegisterAndEnrollUser User",userNameA)
	}
}


func ca_org2() {
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("peerorg2", "enroll_user_peerorg2")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	userNameA := GenerateRandomID()
	_,_,err = myca.RegisterAndEnrollUser(userNameA,"userAW", "org2.department1")
	if err != nil {
		fmt.Println("RegisterAndEnrollUser FAILT",err)
		return
	} else {
		fmt.Println("RegisterAndEnrollUser User",userNameA)
	}
}



func ca_map() {
	/*
	2017/10/20 16:20:54 Error from Register: Error Registering User: Error response from server was: Authorization failure

	报这个错误的原因是：
	   org2用户申请注册使用的是 enroll_user 中org1的admin证书;

	和org2.department1 或者 org1.department1 没有关系
	*/

	InitCA("config.yaml")
	var caMap  map[string]*CA
	caMap = make(map[string]*CA)

	caMap["peerorg1"] = new(CA)
	caMap["peerorg2"] = new(CA)

	// 遍历map
	for k, v := range caMap {
		err := v.InitCaServer(k,"enroll_user")
		if err != nil {
			log.Fatalf("Init CA FAILT: ",err.Error())
		} else {
			fmt.Println("Init CA SUCCESS")
		}
	}

	for k, v := range caMap {
		userNameA := GenerateRandomID()
		_,_,err := v.RegisterAndEnrollUser(userNameA,"userAW", "org2.department1")
		if err != nil {
			fmt.Println("RegisterAndEnrollUser FAILT",k, err)
			return
		} else {
			fmt.Println("RegisterAndEnrollUser User",k,userNameA)
		}
	}
}

func main() {
	ca_org1()
	ca_org2()
	//ca_map()

}

