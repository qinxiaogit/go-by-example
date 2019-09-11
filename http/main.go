package main

import "net/http"
import "./controller"

func main() {
	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	writer.Write([]byte("hello wold"))
	//})

	// wxc96c0498b35b17c8
	//ba369c7f737dae89cc274e730c2738c5
	controller.StartUp()
	//http.ListenAndServe(":8080",nil)
	err := http.ListenAndServeTLS(":443", "/Users/owlet/Documents/code/go/http/cert/server.csr",
		"/Users/owlet/Documents/code/go/ht/cert/server.key", nil)
	panic(err)
}
