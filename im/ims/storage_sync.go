package main

import (
	"github.com/qinxiaogit/go-by-example/im/lru"
	"net"
	"sync"
	log "github.com/golang/glog"
)

type SyncClient struct {
	conn *net.TCPConn
	ewt  chan *lru.Message
}

func NewSyncClient(conn *net.TCPConn)*SyncClient{
	return &SyncClient{
		conn:conn,
		ewt:make(chan *lru.Message,10),
	}
}

func (client *SyncClient)RunLoop(){
	seq := 0
	msg := lru.ReceiveMessage(client.conn)
	if msg == nil{
		return
	}
	if msg.Cmd != lru.MSG_STORAGE_SYNC_BEGIN{
		return
	}
	cursor := msg.Body.(*lru.SyncCursor)
	log.Info("cursor msgid:",cursor.Msgid)
	c := StorageFile{}
}


type Master struct {
	ewt chan *EMessage
	mutex sync.Mutex
	clients map[*SyncCl]
}
