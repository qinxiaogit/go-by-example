package main

import "github.com/qinxiaogit/go-by-example/im/lru"

const BATCH_SIZE = 1000
const PEER_INDEX_FILE_NAM = "peer_index.v3"

type UserID struct {
	appid int64
	uid   int64
}

type UserIndex struct {
	last_msgid int64
	last_id		int64
	last_peer_id	int64
	last_batch_id 	int64
	last_seq_id		int64
}
//在取离线消息时，可以对群组消息和点对点消息分别获取，
//这样可以做到分别控制点对点消息和群组消息读取量，避免单次读取超量的离线消息
type PeerStorage struct {
	*StorageFile
	//消息索引全部放在内存中,在程序退出时,再全部保存到文件中，
	//如果索引文件不存在或上次保存失败，则在程序启动的时候，从消息DB中重建索引，这需要遍历每一条消息
	message_index map[UserID]*UserIndex
}

func NewPeerStorage(f *StorageFile)*PeerStorage{
	storage := &PeerStorage{StorageFile:f}
	storage.message_index = make(map[UserID]*UserIndex)
	return storage
}

func (storage *PeerStorage)SavePeerMessage(appid int64, uid int64, device_id int64, msg *lru.Message) (int64, int64) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	msgid := storage.saveMessage(msg)

	user_index := storage.getPeerIndex(appid,uid)
	last_id := user_index.last_id
	last_peer_id := user_index.last_peer_id
	last_batch_id := user_index.last_batch_id
	last_seq_id := user_index.last_seq_id

	off := &OfflineM
}

func (storage *PeerStorage)getPeerIndex(appid,receiver int64)*UserIndex{
	id := UserID{appid:receiver}
	if ui,ok := storage.message_index[id];ok{
		return ui
	}
	return &UserIndex{}
}