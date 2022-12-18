package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	/*
		客户端调用:

	*/

	// 1.传输协议:http
	c := http.DefaultClient
	// 2.callID:URL
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/add?a=1&b=2", nil)
	res, _ := c.Do(req)
	body, _ := io.ReadAll(res.Body)
	var r RespData
	// 2. 数据格式:json
	_ = json.Unmarshal(body, &r)
	defer res.Body.Close()
	fmt.Print(string(body))
}

type RespData struct {
	Data string `json:"data"`
}
