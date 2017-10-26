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

func ca_register_enroll2() {
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
	//for i:=0;i<10;i++

	secret = "passwd222"
	_, _, err = myca.EnrollUser(userNameA,secret)//可以多次
	if err != nil {
		fmt.Println("RegisterAndEnrollUser FAILT",err)
		return
	} else {
		fmt.Println("RegisterAndEnrollUser User",userNameA)
	}
}

const g_admin = "admin19"

func ca_register_enroll3() {
	fmt.Println("_____________________ca_register_enroll3___________________")
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("peerorg1", "enroll_user_peerorg2")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	userNameA := g_admin
	secret := "passwd"

	secret,err = myca.RegisterClient(userNameA,secret,"org2.department1") //只能一次
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//return
	//for i:=0;i<10;i++
	_, _, err = myca.EnrollUser(userNameA,secret)//可以多次
	if err != nil {
		fmt.Println("RegisterAndEnrollUser FAILT",err)
		return
	} else {
		fmt.Println("RegisterAndEnrollUser User",userNameA)
	}
}


func ca_register_enroll4() {
	fmt.Println("_____________________ca_register_enroll4__________________")
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServerOtherUser(g_admin,"peerorg1", "enroll_use" +
		"" +
			"" +
				"r_peerorg2")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	userNameA := "otheruser"
	secret := "passwd"

	secret,err = myca.Register(userNameA,"passwd","org2.department1") //只能一次
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//return
	//for i:=0;i<10;i++
	_, _, err = myca.EnrollUser(userNameA,secret)//可以多次
	if err != nil {
		fmt.Println("RegisterAndEnrollUser FAILT",err)
		return
	} else {
		fmt.Println("RegisterAndEnrollUser User",userNameA)
	}

}


func ca_register_enroll_reenroll_revoke() {
	InitCA("config.yaml")
	myca := new(CA)
	err := myca.InitCaServer("peerorg1", "enroll_user_peerorg2")
	if err != nil {
		log.Fatalf("Init CA FAILT: ",err.Error())
	} else {
		fmt.Println("Init CA SUCCESS")
	}

	userNameA := "alice"
	/*
	secret := "passwd"

	secret,err = myca.Register(userNameA,secret,"org2.department1") //只能一次
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	_, _, err = myca.EnrollUser(userNameA,secret)//可以多次
	if err != nil {
		fmt.Println("RegisterAndEnrollUser FAILT",err)
		return
	} else {
		fmt.Println("RegisterAndEnrollUser User",userNameA)
	}

	return

	err = myca.ReenrollUser(userNameA)
	if err != nil {
		fmt.Println("ReenrollUser FAILT",err)
		return
	} else {
		fmt.Println("ReenrollUser User",userNameA)
	}
	/*
	只在certificates表增加证书文件
	users 表没有变化
	*/


	//return

	err = myca.RevokeUser(userNameA,"ca-org1")
	if err != nil {
		fmt.Println("RevokeUser FAILT",err)
		return
	} else {
		fmt.Println("RevokeUser User",userNameA)
	}

	/*
	-------
	id              | alice
	token           | \x24326124313024683750704c533754797071766c3056694f534e33612e59517130346a574e574859726e537451764643416a724767347266537a4a71
	type            | user
	affiliation     | org2.department1
	attributes      | null
	state           | -1
	max_enrollments | -1

	user标中的state值为-1 revock之后

	user  certificates 数据表中的数据并没有删除
	*/
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

Register正常 EnrollUser用户使用其他密码：
Error enroling user: Enroll failed: Error response from server was: Authorization failure


 */
func main() {
	//ca_org1()
	//ca_org2()
	//ca_map()
	//ca_register_enroll()
	//ca_register_enroll2()

	//ca_register_enroll3()
	//ca_register_enroll4()

	ca_register_enroll_reenroll_revoke()
}

