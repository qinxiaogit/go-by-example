package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"
	log "github.com/golang/glog"
)
var (
	VERSION string
	BUILD_TIME string
	GO_VERSION string
	GIT_COMMIT_ID string
	GIT_BRANCH	string
)

var server_summary *ServerSummary
var config *StorageConfig
var master *Master

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
	log.Infof("rpc listen:%s storage root:%s sync listen:%s master address:%s is push system:%t group limit:%d offline message limit:%d hard limit:%d\n",
		config.rpc_listen, config.storage_root, config.sync_listen,
		config.master_address, config.is_push_system, config.group_limit,
		config.limit, config.hard_limit)
	log.Infof("http listen address:%s", config.http_listen_address)

	if config.limit == 0{
		log.Error("config limit is 0")
		return
	}
	if config.hard_limit >0 && config.hard_limit/config.limit<2{
		log.Errorf("config limit:%d, hard limit:%d invalid, hard limit/limit must gte 2", config.limit, config.hard_limit)
		return
	}

	storage

}

