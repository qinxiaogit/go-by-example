package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)
var (
	VERSION string
	BUILD_TIME string
	GO_VERSION string
	GIT_COMMIT_ID string
	GIT_BRANCH	string
)

var server_summary *ServerSummary

func init(){
	server_summary = NewServerSummary()
}


func main(){
	fmt.Printf("Version:     %s\nBuilt:       %s\nGo version:  %s\nGit branch:  %s\nGit commit:  %s\n", VERSION, BUILD_TIME, GO_VERSION, GIT_BRANCH, GIT_COMMIT_ID)

	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if len(flag.Args()) == 0{
		fmt.Println("usage: ims config")
		return
	}

	config = read_storage_cfg(flag.Args()[0])
}

