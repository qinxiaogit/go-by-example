package lru

import (
	"bytes"
	"encoding/binary"
)

//主从同步消息
const MSG_STORAGE_SYNC_BEGIN = 220
const MSG_STORAGE_SYNC_MESSAGE = 221
const MSG_STORAGE_SYNC_MESSAGE_BATCH = 222

//内部文件存储使用
//超级群消息队列 代替MSG_GROUP_IM_LIST
const MSG_GROUP_OFFLINE = 247
//个人消息队列 代替MSG_OFFLINE_V3
const MSG_OFFLINE_V4  = 248
const MSG_OFFLINE_V3  = 249
//个人消息队列 代替MSG_OFFLINE_V2
const MSG_OFFLINE_V2 = 250
//im 实例使用
const MSG_PENDING_GROUP_MESSAGE  = 251
//超级群消息队列
// deprecated兼容性
const MSG_GROUP_IM_LIST = 252
//deprecated
const MSG_GROUP_ACK_IN = 253

//deprecated 兼容性
const MSG_OFFLINE = 254

//deprecated
const MSG_ACK_IN = 255

func init(){
	message_creators[MSG_GROUP_OFFLINE] = func() IMessage {return &OfflineMessage4{}}
	message_creators[MSG_OFFLINE_V4] = func() IMessage {
			return &OfflineMessage4{}
	}
	message_creators[MSG_OFFLINE_V3] = func()IMessage{return new (OfflineMessage3)}
	message_creators[MSG_OFFLINE_V2] = func()IMessage{return new (OfflineMessage2)}
	message_creators[MSG_PENDING_GROUP_MESSAGE] = func()IMessage{return new (PendingGroupMessage)}
	message_creators[MSG_GROUP_IM_LIST] = func()IMessage{return new(GroupOfflineMessage)}
	message_creators[MSG_GROUP_ACK_IN] = func()IMessage{return new(IgnoreMessage)}

	message_creators[MSG_OFFLINE] = func()IMessage{return new(OfflineMessage1)}
	message_creators[MSG_ACK_IN] = func()IMessage{return new(IgnoreMessage)}

	message_creators[MSG_STORAGE_SYNC_BEGIN] = func()IMessage{return new(SyncCursor)}
	message_creators[MSG_STORAGE_SYNC_MESSAGE] = func()IMessage{return new(EMessage)}
	message_creators[MSG_STORAGE_SYNC_MESSAGE_BATCH] = func()IMessage{return new(MessageBatch)}


	message_descriptions[MSG_STORAGE_SYNC_BEGIN] = "MSG_STORAGE_SYNC_BEGIN"
	message_descriptions[MSG_STORAGE_SYNC_MESSAGE] = "MSG_STORAGE_SYNC_MESSAGE"
	message_descriptions[MSG_STORAGE_SYNC_MESSAGE_BATCH] = "MSG_STORAGE_SYNC_MESSAGE_BATCH"

	message_descriptions[MSG_GROUP_OFFLINE] = "MSG_GROUP_OFFLINE"
	message_descriptions[MSG_OFFLINE_V4] = "MSG_OFFLINE_V4"
	message_descriptions[MSG_OFFLINE_V3] = "MSG_OFFLINE_V3"
	message_descriptions[MSG_OFFLINE_V2] = "MSG_OFFLINE_V2"
	message_descriptions[MSG_PENDING_GROUP_MESSAGE] = "MSG_PENDING_GROUP_MESSAGE"
	message_descriptions[MSG_GROUP_IM_LIST] = "MSG_GROUP_IM_LIST"
}

type SyncCursor struct {
	Msgid int64
}

func (cursor *SyncCursor) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, cursor.Msgid)
	return buffer.Bytes()
}

func (cursor *SyncCursor) FromData(buff []byte) bool {
	if len(buff) < 8 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &cursor.Msgid)
	return true
}

type EMessage struct {
	Msgid int64
	Device_id int64
	Msg *Message
}

func (emsg *EMessage) ToData() []byte {
	if emsg.Msg == nil {
		return nil
	}

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, emsg.Msgid)
	binary.Write(buffer, binary.BigEndian, emsg.Device_id)
	mbuffer := new(bytes.Buffer)
	WriteMessage(mbuffer, emsg.Msg)
	msg_buf := mbuffer.Bytes()
	var l int16 = int16(len(msg_buf))
	binary.Write(buffer, binary.BigEndian, l)
	buffer.Write(msg_buf)
	buf := buffer.Bytes()
	return buf
}

func (emsg *EMessage) FromData(buff []byte) bool {
	if len(buff) < 18 {
		return false
	}

	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &emsg.Msgid)
	binary.Read(buffer, binary.BigEndian, &emsg.Device_id)
	var l int16
	binary.Read(buffer, binary.BigEndian, &l)
	if int(l) > buffer.Len() {
		return false
	}

	msg_buf := make([]byte, l)
	buffer.Read(msg_buf)
	mbuffer := bytes.NewBuffer(msg_buf)
	//recusive
	msg := ReceiveMessage(mbuffer)
	if msg == nil {
		return false
	}
	emsg.Msg = msg

	return true
}

type OfflineMessage struct {
	Appid    int64
	Receiver int64 //用户id or 群组id
	Msgid    int64 //消息本体的id
	Device_id int64
	Seq_id   int64      //v4 消息序号, 1,2,3...
	Prev_msgid  int64 //个人消息队列(点对点消息，群组消息)
	Prev_peer_msgid int64 //v2 点对点消息队列
	Prev_batch_msgid int64 //v3 0<-1000<-2000<-3000...构成一个消息队列
}
func (off *OfflineMessage)Body() *OfflineMessage{
	return off
}

type OfflineMessage1 struct {
	OfflineMessage
}

func (off *OfflineMessage1)ToData() []byte {
	buffer := &bytes.Buffer{}
	binary.Write(buffer,binary.BigEndian,off.Appid)
	binary.Write(buffer, binary.BigEndian, off.Receiver)
	binary.Write(buffer, binary.BigEndian, off.Msgid)
	binary.Write(buffer, binary.BigEndian, off.Device_id)
	binary.Write(buffer, binary.BigEndian, off.Prev_msgid)
	buf := buffer.Bytes()
	return buf
}

func (off *OfflineMessage1) FromData(buff []byte) bool {
	if len(buff) < 32 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &off.Appid)
	binary.Read(buffer, binary.BigEndian, &off.Receiver)
	binary.Read(buffer, binary.BigEndian, &off.Msgid)
	if len(buff) == 40 {
		binary.Read(buffer, binary.BigEndian, &off.Device_id)
	}
	binary.Read(buffer, binary.BigEndian, &off.Prev_msgid)

	off.Prev_peer_msgid = off.Prev_msgid
	return true
}

type OfflineMessage2 struct {
	OfflineMessage
}
func (off *OfflineMessage2) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, off.Appid)
	binary.Write(buffer, binary.BigEndian, off.Receiver)
	binary.Write(buffer, binary.BigEndian, off.Msgid)
	binary.Write(buffer, binary.BigEndian, off.Device_id)
	binary.Write(buffer, binary.BigEndian, off.Prev_msgid)
	binary.Write(buffer, binary.BigEndian, off.Prev_peer_msgid)
	buf := buffer.Bytes()
	return buf
}

func (off *OfflineMessage2) FromData(buff []byte) bool {
	if len(buff) < 48 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &off.Appid)
	binary.Read(buffer, binary.BigEndian, &off.Receiver)
	binary.Read(buffer, binary.BigEndian, &off.Msgid)
	binary.Read(buffer, binary.BigEndian, &off.Device_id)
	binary.Read(buffer, binary.BigEndian, &off.Prev_msgid)
	binary.Read(buffer, binary.BigEndian, &off.Prev_peer_msgid)
	return true
}


type OfflineMessage3 struct {
	OfflineMessage
}

func (off *OfflineMessage3) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, off.Appid)
	binary.Write(buffer, binary.BigEndian, off.Receiver)
	binary.Write(buffer, binary.BigEndian, off.Msgid)
	binary.Write(buffer, binary.BigEndian, off.Device_id)
	binary.Write(buffer, binary.BigEndian, off.Prev_msgid)
	binary.Write(buffer, binary.BigEndian, off.Prev_peer_msgid)
	binary.Write(buffer, binary.BigEndian, off.Prev_batch_msgid)
	buf := buffer.Bytes()
	return buf
}

func (off *OfflineMessage3) FromData(buff []byte) bool {
	if len(buff) < 56 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &off.Appid)
	binary.Read(buffer, binary.BigEndian, &off.Receiver)
	binary.Read(buffer, binary.BigEndian, &off.Msgid)
	binary.Read(buffer, binary.BigEndian, &off.Device_id)
	binary.Read(buffer, binary.BigEndian, &off.Prev_msgid)
	binary.Read(buffer, binary.BigEndian, &off.Prev_peer_msgid)
	binary.Read(buffer, binary.BigEndian, &off.Prev_batch_msgid)
	return true
}

type OfflineMessage4 struct {
	OfflineMessage
}

func (off *OfflineMessage4) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, off.Appid)
	binary.Write(buffer, binary.BigEndian, off.Receiver)
	binary.Write(buffer, binary.BigEndian, off.Msgid)
	binary.Write(buffer, binary.BigEndian, off.Device_id)
	binary.Write(buffer, binary.BigEndian, off.Seq_id)
	binary.Write(buffer, binary.BigEndian, off.Prev_msgid)
	binary.Write(buffer, binary.BigEndian, off.Prev_peer_msgid)
	binary.Write(buffer, binary.BigEndian, off.Prev_batch_msgid)
	buf := buffer.Bytes()
	return buf
}

func (off *OfflineMessage4) FromData(buff []byte) bool {
	if len(buff) < 64 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &off.Appid)
	binary.Read(buffer, binary.BigEndian, &off.Receiver)
	binary.Read(buffer, binary.BigEndian, &off.Msgid)
	binary.Read(buffer, binary.BigEndian, &off.Device_id)
	binary.Read(buffer, binary.BigEndian, &off.Seq_id)
	binary.Read(buffer, binary.BigEndian, &off.Prev_msgid)
	binary.Read(buffer, binary.BigEndian, &off.Prev_peer_msgid)
	binary.Read(buffer, binary.BigEndian, &off.Prev_batch_msgid)
	return true
}

type GroupOfflineMessage struct {
	OfflineMessage
}
func (off *GroupOfflineMessage) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, off.Appid)
	binary.Write(buffer, binary.BigEndian, off.Receiver)
	binary.Write(buffer, binary.BigEndian, off.Msgid)
	binary.Write(buffer, binary.BigEndian, off.Receiver)
	binary.Write(buffer, binary.BigEndian, off.Device_id)
	binary.Write(buffer, binary.BigEndian, off.Prev_msgid)
	buf := buffer.Bytes()
	return buf
}

func (off *GroupOfflineMessage) FromData(buff []byte) bool {
	if len(buff) < 40 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &off.Appid)
	binary.Read(buffer, binary.BigEndian, &off.Receiver)
	binary.Read(buffer, binary.BigEndian, &off.Msgid)
	binary.Read(buffer, binary.BigEndian, &off.Receiver)
	if len(buff) == 48 {
		binary.Read(buffer, binary.BigEndian, &off.Device_id)
	}
	binary.Read(buffer, binary.BigEndian, &off.Prev_msgid)
	return true
}



//待发送的群组消息临时存储结构
type PendingGroupMessage struct {
	appid     int64
	sender    int64
	device_ID int64 //发送者的设备id
	gid       int64
	timestamp int32

	members   []int64 //需要接受此消息的成员列表
	content   string
}

func (gm *PendingGroupMessage) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, gm.appid)
	binary.Write(buffer, binary.BigEndian, gm.sender)
	binary.Write(buffer, binary.BigEndian, gm.device_ID)
	binary.Write(buffer, binary.BigEndian, gm.gid)
	binary.Write(buffer, binary.BigEndian, gm.timestamp)

	count := int16(len(gm.members))
	binary.Write(buffer, binary.BigEndian, count)
	for _, uid := range gm.members {
		binary.Write(buffer, binary.BigEndian, uid)
	}

	buffer.Write([]byte(gm.content))
	buf := buffer.Bytes()
	return buf
}

func (gm *PendingGroupMessage) FromData(buff []byte) bool {
	if len(buff) < 38 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &gm.appid)
	binary.Read(buffer, binary.BigEndian, &gm.sender)
	binary.Read(buffer, binary.BigEndian, &gm.device_ID)
	binary.Read(buffer, binary.BigEndian, &gm.gid)
	binary.Read(buffer, binary.BigEndian, &gm.timestamp)

	var count int16
	binary.Read(buffer, binary.BigEndian, &count)

	if len(buff) < int(38 + count*8) {
		return false
	}

	gm.members = make([]int64, count)
	for i := 0; i < int(count); i++ {
		var uid int64
		binary.Read(buffer, binary.BigEndian, &uid)
		gm.members[i] = uid
	}
	offset := 38 + count*8
	gm.content = string(buff[offset:])

	return true
}

type IOfflineMessage interface {
	Body()*OfflineMessage
}
type MessageBatch struct {
	first_id int64
	last_id  int64
	msgs     []*Message
}

func (batch *MessageBatch) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, batch.first_id)
	binary.Write(buffer, binary.BigEndian, batch.last_id)
	count := int32(len(batch.msgs))
	binary.Write(buffer, binary.BigEndian, count)

	for _, m := range batch.msgs {
		WriteMessage(buffer, m)
	}

	buf := buffer.Bytes()
	return buf
}

func (batch *MessageBatch) FromData(buff []byte) bool {
	if len(buff) < 18 {
		return false
	}

	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &batch.first_id)
	binary.Read(buffer, binary.BigEndian, &batch.last_id)

	var count int32
	binary.Read(buffer, binary.BigEndian, &count)

	batch.msgs = make([]*Message, 0, count)
	for i := 0; i < int(count); i++ {
		msg := ReceiveMessage(buffer)
		if msg == nil {
			return false
		}
		batch.msgs = append(batch.msgs, msg)
	}

	return true
}




