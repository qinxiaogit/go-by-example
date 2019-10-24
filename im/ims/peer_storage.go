package main

import (
	"github.com/qinxiaogit/go-by-example/im/lru"
	log "github.com/golang/glog"
	)
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

	user_index := storage.getPeerIndex(appid, uid)
	last_id := user_index.last_id
	last_peer_id := user_index.last_peer_id
	last_batch_id := user_index.last_batch_id
	last_seq_id := user_index.last_seq_id

	off := &lru.OfflineMessage4{}
	off.Appid = appid
	off.Receiver = uid
	off.Msgid = msgid
	off.Device_id = device_id
	off.Seq_id = last_seq_id
	off.Prev_msgid = last_id
	off.Prev_peer_msgid = last_peer_id
	off.Prev_batch_msgid = last_batch_id

	var flag int
	if storage.isGroupMessage(msg){
		flag = lru.MESSAGE_FLAG_GROUP
	}
	m := &lru.Message{Cmd:lru.MSG_OFFLINE_V4,Flag:flag,Body:off}
	last_id = storage.saveMessage(m)
	if !storage.isGroupMessage(msg){
		last_peer_id = last_id
	}
	last_seq_id+= 1
	if last_seq_id%BATCH_SIZE == 0{
		last_batch_id = last_id
	}
	ui := &UserIndex{msgid,last_id,last_peer_id,last_batch_id,last_seq_id}
	storage.setPeerIndex(appid,uid,ui)
	return msgid,user_index.last_msgid
}

func (storage *PeerStorage) setPeerIndex(appid int64,receiver int64,ui *UserIndex){
	id := UserID{appid,receiver}

	storage.message_index[id] = ui

	if ui.last_id > storage.last_id{
		storage.last_id = ui.last_id
	}
}

func (storage *PeerStorage)getPeerIndex(appid,receiver int64)*UserIndex{
	id := UserID{appid:receiver}
	if ui,ok := storage.message_index[id];ok{
		return ui
	}
	return &UserIndex{}
}

func (client *PeerStorage)isGroupMessage(msg *lru.Message)bool{
	return msg.Cmd == lru.MSG_GROUP_IM||msg.Flag &lru.MESSAGE_FLAG_GROUP != 0
}
//获取最近离线消息id
func (storage *PeerStorage)getLastMessageId(appid int64,receiver int64)(int64,int64){
	id := UserID{appid,receiver}
	if ui,ok := storage.message_index[id];ok{
		return ui.last_id,ui.last_peer_id
	}
	return 0,0
}
//lock
func (storage *PeerStorage)GetLastMessageId(appid int64,receiver int64)(int64,int64){
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	return storage.getLastMessageId(appid,receiver)
}
//获取所有消息id大于sync_msgid的消息,
//group_limit&limit:0 表示无限制
//消息超过group_limit后，只获取点对点消息
//总消息数限制在limit
func (storage *PeerStorage)LoadHistoryMessages(appid,receiver,sync_msgid int64,group_limit,limit int)([]*lru.EMessage ,int64,bool){
	var last_msgid int64
	last_id,_:= storage.GetLastMessageId(appid,receiver)
	messages := make([]*lru.EMessage,0,10)
	for{
		if last_id == 0{
			break
		}
		msg := storage.LoadMessage(last_id)
		if msg == nil{
			break
		}
		var off *lru.OfflineMessage
		if ioff,ok := msg.Body.(lru.IOfflineMessage);ok{
			off = ioff.Body()
		}else{
			log.Warning("invalid message cmd:",msg.Cmd)
			break
		}
		if last_msgid == 0{
			last_msgid = off.Msgid
		}
		if off.Msgid <= sync_msgid{
			break
		}
		msg = storage.LoadMessage(off.Msgid)
		if msg == nil{
			break
		}
		if msg.Cmd != lru.MSG_GROUP_IM&&
			msg.Cmd!=lru.MSG_GROUP_NOTIFICATION&&
			msg.Cmd!=lru.MSG_IM&&
			msg.Cmd!=lru.MSG_CUSTOMER&&
			msg.Cmd!=lru.MSG_CUSTOMER_SUPPORT&&
			msg.Cmd!= lru.MSG_SYSTEM{
			if group_limit>0&&len(messages)>=group_limit{
				last_id = off.Prev_peer_msgid
			}else{
				last_id = off.Prev_msgid
			}
			continue
		}
		emsg := &lru.EMessage{Msgid:off.Msgid, Device_id:off.Device_id, Msg:msg}
		messages = append(messages, emsg)
		if limit>0&&len(messages)>=limit{
			break
		}
		if group_limit>0&&len(messages)>=group_limit{
			last_id = off.Prev_peer_msgid
		}else{
			last_id = off.Prev_msgid
		}
	}
	if len(messages) > 1000 {
		log.Warningf("appid:%d uid:%d sync msgid:%d history message overflow:%d",
			appid, receiver, sync_msgid, len(messages))
	}
	log.Infof("appid:%d uid:%d sync msgid:%d history message loaded:%d %d",
		appid, receiver, sync_msgid, len(messages), last_msgid)

	return messages, last_msgid, false
}
//加载离线消息
//limit 单次加载限制
//hard_limit 总量限制
//hard_limit == 0 or ( hard_limit >= BATCH_SIZE and hard_limit > limit)
func (storage *PeerStorage) LoadHistoryMessagesV3(appid int64, receiver int64, sync_msgid int64, limit int, hard_limit int) ([]*EMessage, int64, bool) {
	var last_msgid int64
	var last_offline_msgid int64
	msg_index := storage.getPeerIndex(appid,receiver)

	var last_batch_id int64
	if msg_index != nil{
		last_batch_id = msg_index.last_batch_id
	}
	hard_batch_count := hard_limit/BATCH_SIZE
	batch_count := limit/BATCH_SIZE

	batch_ids := make([]int64,0,10)
	//搜索和sync_msgid最近的batch_id
	for{
		if last_batch_id<= sync_msgid{
			break
		}
		msg:= storage.LoadMessage(last_batch_id)
		if msg == nil{
			break
		}
		var off *lru.OfflineMessage
		if ioff,ok := msg.Body.(lru.IOfflineMessage);ok{
			off = ioff.Body()
		}else{
			log.Warning("invalid message cmd:",msg.Cmd)
			break
		}
		if off.Msgid <= sync_msgid{
			break
		}
		batch_ids = append(batch_ids,last_batch_id)
		last_batch_id = off.Prev_batch_msgid
		if hard_batch_count>0 && len(batch_ids)>= hard_batch_count{
			break
		}
	}
	//只获取点对点消息
	var is_peer bool
	var last_id int64
	if hard_batch_count > 0 && len(batch_ids) >= hard_batch_count && hard_batch_count >= batch_count {
		index := len(batch_ids) - batch_count
		last_id = batch_ids[index]
		//离线消息总数超过hard_limit时，首先加载数量不超过limit的点对点消息
		//尽量保证点对点消息的投递
		is_peer = true
	} else if len(batch_ids) >= batch_count {
		index := len(batch_ids) - batch_count
		last_id = batch_ids[index]
	} else if msg_index != nil {
		last_id = msg_index.last_id
	}
	messages := make([]*lru.EMessage,0,10)
	for{
		msg := storage.LoadMessage(last_id)
		if msg == nil{
			break
		}
		var off *lru.OfflineMessage
		if ioff,ok:=msg.Body.(lru.IOfflineMessage);ok{
			off = ioff.Body()
		}else{
			log.Warning("invalid message cmd:",msg.Cmd)
			break
		}
		if last_msgid == 0{
			last_msgid = off.Msgid
			last_offline_msgid = last_id
		}
		if off.Msgid <= sync_msgid{
			break
		}
		msg = storage.LoadMessage(off.Msgid)
		if msg == nil{
			break
		}
		if msg.Cmd != lru.MSG_GROUP_IM &&
			msg.Cmd != lru.MSG_GROUP_NOTIFICATION &&
			msg.Cmd != lru.MSG_IM &&
			msg.Cmd != lru.MSG_CUSTOMER &&
			msg.Cmd != lru.MSG_CUSTOMER_SUPPORT &&
			msg.Cmd != lru.MSG_SYSTEM {
			if is_peer {
				last_id = off.Prev_peer_msgid
			} else {
				last_id = off.Prev_msgid
			}
			continue
		}
		emsg := &lru.EMessage{Msgid:off.Msgid,Device_id:off.Device_id}
		messages = append(messages,emsg)
		if limit>0&&len(messages)>=limit{
			break
		}
		if is_peer{
			last_id = off.Prev_peer_msgid
		}else{
			last_id = off.Prev_msgid
		}
	}
	if len(messages)>1000{
		log.Warningf("appid:%d uid:%d sync msgid:%d history message overflow:%d",
			appid, receiver, sync_msgid, len(messages))
	}
	var hasMore bool
	if msg_index != nil && last_offline_msgid > 0 &&
		last_offline_msgid < msg_index.last_id {
		hasMore = true
	}

	log.Infof("appid:%d uid:%d sync msgid:%d history message loaded:%d %d, last_id:%d, %d last batch id:%d last seq id:%d has more:%t only peer:%t",
		appid, receiver, sync_msgid, len(messages), last_msgid,
		last_offline_msgid, msg_index.last_id, msg_index.last_batch_id,
		msg_index.last_seq_id, hasMore, is_peer)

	return messages, last_msgid, hasMore
}
//读取最新消息
func (storage *PeerStorage)LoadLatestMessages(appid int64, receiver int64, limit int) []*lru.EMessage {
	last_id,_ := storage.getLastMessageId(appid,receiver)
	messages := make([]*lru.EMessage,0,10)
	for{
		if last_id == 0{
			break
		}
		msg := storage.LoadMessage(last_id)
		if msg == nil{
			break
		}
		var off *lru.OfflineMessage
		if ioff,ok := msg.Body.(lru.IOfflineMessage);ok{
			off = ioff.Body()
		}else{
			log.Warning("invalid message cmd:",msg.Cmd)
			break
		}
		msg = storage.LoadMessage(off.Msgid)
		if msg == nil{
			break
		}
		if msg.Cmd != lru.MSG_GROUP_IM &&
			msg.Cmd != lru.MSG_GROUP_NOTIFICATION &&
			msg.Cmd != lru.MSG_IM &&
			msg.Cmd != lru.MSG_CUSTOMER &&
			msg.Cmd != lru.MSG_CUSTOMER_SUPPORT {
			last_id = off.Prev_msgid
			continue
		}
		emsg := &lru.EMessage{Msgid:off.Msgid,Device_id:off.Device_id,Msg:msg}
		messages = append(messages,emsg)
		if len(messages)>= limit{
			break
		}
		last_id = off.Prev_msgid
	}
	return messages
}