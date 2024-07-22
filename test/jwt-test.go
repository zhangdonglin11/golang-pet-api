package main

import (
	"fmt"
	"golang-pet-api/utils/jwt"
	"time"
)

func main() {
	//	生成token
	data := map[string]interface{}{
		"AA": 10,
		"BB": uint(20),
		"CC": "30",
	}
	token, err := jwt.GetJwtToken("abc", time.Now().Unix(), 60, data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

	//解析
	claims, err := jwt.ParseJwtToken("abc", token)
	if err != nil {
		return
	}
	fmt.Println(claims)
	for k, v := range claims {
		fmt.Println(k, v)
	}
	aa := claims["AA"]
	bb := claims["BB"]
	cc := claims["CC"]
	i := aa.(float64)
	j := bb.(float64)
	k := cc.(string)
	fmt.Println(i)
	fmt.Println(j)
	fmt.Println(k)

	fmt.Println(uint(i))

}
