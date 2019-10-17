package main

import (
	"log"
	"strconv"
	"github.com/qinxiaogit/go-by-example/im/cfg"
)

//超级群离线消息数量限制，超过的部分被丢弃
const GROUP_OFFLINE_LIMIT = 100

//离线消息返回数量限制
const OFFLINE_DEFAULT_LIMIT = 3000
const GROUP_OFFLINE_DEFAULT_LIMIT = 0

//unlimit
const OFFLINE_DEFAULT_HARD_LIMIT = 0

type StorageConfig struct {
	rpc_listen            string
	storage_root        string
	kefu_appid          int64
	http_listen_address string

	sync_listen    string
	master_address string
	is_push_system bool
	group_limit    int //普通群离线消息的数量限制
	limit          int //单次离线消息数量限制
	hard_limit     int //离线消息总的数量限制
}

func get_int(app_cfg map[string]string, key string) int64 {
	concurrency, present := app_cfg[key]
	if !present {
		log.Fatalf("key:%s non exist", key)
	}
	n, err := strconv.ParseInt(concurrency, 10, 64)
	if err != nil {
		log.Fatalf("key:%s is't interger", key)
	}
	return n
}

func get_opt_int(app_cfg map[string]string, key string, default_value int64) int64 {
	concurrency, present := app_cfg[key]
	if !present {
		return default_value
	}
	n, err := strconv.ParseInt(concurrency, 10, 64)
	if err != nil {
		log.Fatalf("key:%s is't integer", key)
	}
	return n
}

//get_string 获取key
func get_string(app_cfg map[string]string, key string) string {
	concurrency, present := app_cfg[key]
	if !present {
		log.Fatalf("key:%s non exist", key)
	}
	return concurrency
}

//get_opt_string 获取key
func get_opt_string(app_cfg map[string]string, key string) string {
	concurrency, present := app_cfg[key]
	if !present {
		return ""
	}
	return concurrency
}

func read_storage_cfg(cfg_path string) *StorageConfig {
	config := StorageConfig{}
	app_cfg := make(map[string]string)
	err := cfg.Load(cfg_path, app_cfg)
	if err != nil {
		log.Fatal(err)
	}
	config.rpc_listen = get_string(app_cfg, "rpc_listen")
	config.http_listen_address = get_opt_string(app_cfg, "http_listen_address")
	config.storage_root = get_string(app_cfg,"storage_root")
	config.kefu_appid = get_int(app_cfg,"kefu_appid")
	config.sync_listen = get_string(app_cfg,"sync_listen")
	config.is_push_system = get_opt_int(app_cfg,"is_push_system",0) ==1
	config.limit = int(get_opt_int(app_cfg,"limit",OFFLINE_DEFAULT_LIMIT))
	config.group_limit= int(get_opt_int(app_cfg,"group_limit",GROUP_OFFLINE_DEFAULT_LIMIT))
	config.hard_limit = int(get_opt_int(app_cfg,"hard_limit",OFFLINE_DEFAULT_HARD_LIMIT))
	return &config
}
