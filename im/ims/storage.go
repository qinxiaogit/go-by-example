package main

import (
	log "github.com/golang/glog"
	"github.com/qinxiaogit/go-by-example/im/lru"
	"os"
)
type Storage struct {
	*StorageFile
	*PeerStorage
	*GroupStorage
}

func NewStorage(root string)*Storage{
	f := NewStorageFile(root)
	ps:= NewPeerStorage(f)
	gs:= NewGroupStorage(f)

	storage := &Storage{f,ps,gs}
	r1 := storage.ReadPeerIndex()
	r2 := storage.readGroupIndex()

	storage.last_saved_id = storage.last_id
	if r1{
		storage.RepairPeerIndex()
	}
	if r2 {
		storage.repairGroupIndex()
	}
	if !r1{
		storage.createPeerIndex()
	}
	log.Infof("last id:%d last saved id:%d", storage.last_id, storage.last_saved_id)

	storage.FlushIndex()
	return storage
}
func (storage *Storage)LoadSyncMessagesInBackground(cursor int64)chan *lru.MessageBatch{
	c := make(chan *lru.MessageBatch,10)
	go func() {
		defer close(c)
		block_NO := storage.getBlockNO(cursor)
		offset := storage.getBlockOffset(cursor)

		n := block_NO
		for {
			file := storage.openReadFile(n)
			if file == nil{
				break
			}
			if n == block_NO{
				file_size,err := file.Seek(0,os.SEEK_END)
				if err != nil{
					log.Fatal("seek file err:",err)
					return
				}
				if file_size < int64(offset){
					break
				}
				_,err = file.Seek(int64(offset),os.SEEK_SET)
				if err!
			}
		}
	}()
}

func (storage *Storage)flushIndex(){
	storage.mutex.Lock()
	last_id := storage.last_id
	peer_index := storage.clonePeerIndex()
	group_index:= storage.cloneGroupIndex()
	storage.mutex.Unlock()
	storage.savePeerIndex(peer_index)
	storage.saveGroupIndex(group_index)

	storage.last_saved_id = last_id
}

func (storage *Storage)FlushIndex(){
	do_flush := false
	if storage.last_id - storage.last_saved_id >2*BLOCK_SIZE{
		do_flush = true
	}
	if do_flush{
		storage.flushIndex()
	}
}



