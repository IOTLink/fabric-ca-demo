package main

import (
	. "fabric-ca-demo/modelv4/fabric"
	"log"
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomID generates random ID
func GenerateRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return randomString(10)
}

// Utility to create random string of strlen length
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}


func main() {
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("caorg1", "enroll_user_peerorg1")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
		return
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
