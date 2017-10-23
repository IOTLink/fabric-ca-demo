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

func ca_register_enroll() {
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("peerorg1", "enroll_user_peerorg2")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	userNameA := "lhy9"
	secret := "passwd"

	secret,err = myca.Register(userNameA,"passwd","org2.department1") //只能一次
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//return
	for i:=0;i<10;i++ {
		_, _, err = myca.EnrollUser(userNameA,secret)//可以多次
		if err != nil {
			fmt.Println("RegisterAndEnrollUser FAILT",err)
			return
		} else {
			fmt.Println("RegisterAndEnrollUser User",userNameA)
		}
	}
}

/*
同一个用户名　和密码：
2017/10/23 21:04:58 Error from Register: Error Registering User: Error response from server was: Identity 'lhy2' is already registered

同一个用户名　不同密码：
2017/10/23 21:06:23 Error from Register: Error Registering User: Error response from server was: Identity 'lhy2' is already registered

ca-server 配置不存在org2.department2
2017/10/23 21:09:16 Error from Register: Error Registering User: Error response from server was: Failed getting affiliation 'org2.department2': sql: no rows in result set


设置ca server 配置文件 registry --> maxenrollments = 0 则限制注册
2017/10/23 21:16:38 Enroll return error: %v Enroll failed: Error response from server was: Enrollments are disabled; user 'admin' cannot enroll
2017/10/23 21:16:38 Init CA FAILT: %!(EXTRA string=Enroll failed: Error response from server was: Enrollments are disabled; user 'admin' cannot enroll)

设置ca server 配置文件 registry --> maxenrollments = 3 则限制注册3次，及使用同一个用户名的情况下
数据库user表　state记录申请了多少次：
Error enroling user: Enroll failed: Error response from server was: Authorization failure


 */
func main() {
	//ca_org1()
	//ca_org2()
	//ca_map()
	ca_register_enroll()

}

