package main

import (
	"github.com/imroc/req"
	"fmt"
)

func main()  {

	param := req.Param{
		"user": "manu", "password": "123",
	}
	url := "http://127.0.0.1:8080/loginJSON"
	r,err := req.Post(url,param)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(r.String())
}
