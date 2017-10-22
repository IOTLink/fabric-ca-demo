package main

import (
	. "fabric-ca-demo/modelv4/fabric"
	"log"
	"fmt"
	"math/rand"
	"time"
	"strconv"
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

func testca1() {
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

func testca2() {
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("caorg1", "enroll_user_peerorg1")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
		return
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	var count uint64 = 0
	for {
		count++
		userNameA := GenerateRandomID() + strconv.FormatUint(count,10)
		_,_,err = myca.RegisterAndEnrollUser(userNameA,"userAW", "org1.department1")
		if err != nil {
			fmt.Println("RegisterAndEnrollUser FAILT",err)
			return
		} else {
			fmt.Println("RegisterAndEnrollUser User",userNameA)
		}
	}
}


func main() {
	//testca1
	testca2()
}
