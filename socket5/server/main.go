package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"io"
	"net"
	"net/http"
	"socket5/common"
	"sync"
)

var gVersion = "taosocks/20190722"
var gForWard = "https://example.com"
var gListen string
var gKey string
var gPath string

var tslog = &common.TSLog{}

func doRelay(conn net.Conn,bio *bufio.ReadWriter) error{
	var err error
	defer conn.Close()

	enc :=  gob.NewEncoder(bio)
	dec :=  gob.NewDecoder(bio)

	var openMsg common.OpenMessage

	if err = dec.Decode(&openMsg);err!=nil{
		return err
	}
	tslog.Log("> %s",openMsg.Addr)

	outConnm,err := net.Dial("tcp",openMsg.Addr)
	if err!=nil{
		err2:=enc.Encode(common.OpenAckMessage{
			Status:false,
		})
		if err2!=nil{
			tslog.Log("> %s","加密失败")
			return err2
		}
		err3 := bio.Flush()
		if err3!=nil{
			tslog.Log("> %s","刷新缓存失败")
			return err3
		}
		return err
	}
	defer outConnm.Close()
	if err := enc.Encode(common.OpenAckMessage{
		Status:true,
	}); err != nil {
		tslog.Red("%s", err.Error())
		return err
	}
	err = bio.Flush()
	if err !=nil{
		tslog.Log("> %s","-刷新缓存失败-")
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		buf := make([]byte,common.ReadBuffSize)
		for{
			n,err:=outConnm.Read(buf)
			if err !=nil{
				if err!= io.EOF{
					tslog.Red("%s",err.Error())
				}
				return
			}
			var msg common.RelayMessage
			msg.Data = buf[:n]
			if err :=enc.Encode(&msg);err!=nil{
				tslog.Red("%s",err.Error())
				return
			}
			bio.Flush()
		}
	}()

	go func() {
		defer wg.Done()

		for {
			var msg common.RelayMessage
			if err:=dec.Decode(&msg);err!=nil{
				if err!=io.EOF{
					tslog.Red("%s",err.Error())
				}
				return
			}
			_,err :=outConnm.Write(msg.Data)
			if err!=nil{
				tslog.Red("%s",err.Error())
				return
			}
		}
	}()
	wg.Wait()
	tslog.Gray("<%s",openMsg.Addr)
	return nil
}

func doForward(w http.ResponseWriter,req *http.Request){
	resp,err := http.Get(gForWard+req.RequestURI)
	if err!=nil{
		w.WriteHeader(resp.StatusCode)
		io.Copy(w,resp.Body)
	}
}

func handleRequest(w http.ResponseWriter,req *http.Request){
	ver := req.Header.Get("Upgrade")
	auth := req.Header.Get("Authorization")
	path := req.URL.Path
	tslog.Green("%s",req.URL.Path)
	tslog.Gray("-----------------------------")
	if path == gPath && ver == gVersion && auth == "taosocks" + gKey{
		w.WriteHeader(http.StatusSwitchingProtocols)
		conn,bio,_ := w.(http.Hijacker).Hijack()
		doRelay(conn,bio)
	}else{
		doForward(w,req)
	}
}

func main(){
	flag.StringVar(&gListen,"listen",":1081","listen address(host:port)")
	flag.StringVar(&gForWard,"forward","https://example.com", "delegate website, format must be https://example.com")
	flag.StringVar(&gKey,"key","","the key")
	flag.StringVar(&gPath,"path","/","/your/path")

	flag.Parse()

	http.HandleFunc("/",handleRequest)
	if err := http.ListenAndServeTLS(gListen,
		"config/server.crt",
		"config/server.key",
		nil,
		); err!=nil{
			panic(err)
	}
}