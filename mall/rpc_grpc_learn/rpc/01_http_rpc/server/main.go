package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// 服务端提供服务的函数:
	// callID: URL解决 http://127.0.0.1:8080/add?a=1&b=2
	// 网络传输协议: 使用http
	// 数据传输协议: 使用json
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm() // 解析参数
		fmt.Println("path:", r.URL.Path)
		a, _ := strconv.Atoi(r.Form["a"][0])
		b, _ := strconv.Atoi(r.Form["a"][0])
		w.Header().Set("Content-Type", "application/json")
		jData, _ := json.Marshal(map[string]int{
			"data": a + b,
		})
		_, _ = w.Write(jData)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
