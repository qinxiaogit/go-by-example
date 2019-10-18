package lru

import (
	"bytes"
	"encoding/binary"
)

//主从同步消息
const MSG_STORAGE_SYNC_BEGIN = 220
const MSG_STORAGE_SYNC_MESSAGE = 221
const MSG_STORAGE_SYNC_MESSAGE_BATCH = 222

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
	device_id int64
	Msg *Message
}

func (emsg *EMessage) ToData() []byte {
	if emsg.Msg == nil {
		return nil
	}

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, emsg.Msgid)
	binary.Write(buffer, binary.BigEndian, emsg.device_id)
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
	binary.Read(buffer, binary.BigEndian, &emsg.device_id)
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
	appid    int64
	receiver int64 //用户id or 群组id
	msgid    int64 //消息本体的id
	device_id int64
	seq_id   int64      //v4 消息序号, 1,2,3...
	prev_msgid  int64 //个人消息队列(点对点消息，群组消息)
	prev_peer_msgid int64 //v2 点对点消息队列
	prev_batch_msgid int64 //v3 0<-1000<-2000<-3000...构成一个消息队列
}
func (off *OfflineMessage)body() *OfflineMessage{
	return off
}

type OfflineMessage1 struct {
	OfflineMessage
}
