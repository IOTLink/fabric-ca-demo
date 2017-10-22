package main

import (
. "fabric-ca-demo/modelv2/fabric"
"log"
"fmt"
)

func main() {
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





