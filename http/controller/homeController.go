package controller

import (
	"net/http"
)

type home struct {}

func (h home)RegisterRoutes(){
	http.HandleFunc("/",indexHandle)
}

func indexHandle(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json;charset=utf-8")
	w.WriteHeader(301)

	_,err:=w.Write([]byte("{'json':'xxx'}"))
	if err !=nil{
		panic(err)
	}

}