package lru

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	log "github.com/golang/glog"
)

//平台账号
const PLATFORM_IOS = 1
const PLATFORM_ANDROID = 2
const PLATFORM_WEB = 3

const DEFAULT_VERSION = 2
const MSG_HEADER_SIZE = 12

var message_descriptions  map[int]string = make(map[int]string)

type MessageCreator func()IMessage

var message_creators map[int]MessageCreator = make(map[int]MessageCreator)

type VersionMessageCreator func()IVersionMessage
var vmessage_creators map[int]VersionMessageCreator = make(map[int]VersionMessageCreator)

//true client->server
var external_messages[256]bool

//写头部数据
func WriteHeader(len ,seq int32 ,cmd byte,version ,flag byte,buffer io.Writer){
	binary.Write(buffer,binary.BigEndian,len)
	binary.Write(buffer,binary.BigEndian,seq)
	t:= []byte{cmd,byte(version),flag,0}
	buffer.Write(t)
}
//读头文件
func ReadHeader(buff []byte)(int,int,int,int,int){
	var length int32
	var seq int32
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer,binary.BigEndian,&length)
	binary.Read(buffer,binary.BigEndian,&seq)

	cmd,_:= buffer.ReadByte()
	version,_:= buffer.ReadByte()
	flag,_:= buffer.ReadByte()
	return int(length),int(seq),int(cmd),int(version),int(flag)
}
//写消息数据
func WriteMessage(w *bytes.Buffer,msg *Message){
	body:= msg.ToData()
	WriteHeader(int32(len(body)), int32(msg.seq), byte(msg.cmd), byte(msg.version), byte(msg.flag), w)
	w.Write(body)
}
//发送消息
func SendMessage(conn io.Writer,msg *Message)error{
	buffer:= new(bytes.Buffer)
	WriteMessage(buffer,msg)
	buf:= buffer.Bytes()
	n,err:= conn.Write(buf)
	if err!= nil{
		log.Infof("write less:%d %d",n,len(buf))
		return errors.New("write less")
	}
	return nil
}
func ReceiveLimitMessage(conn io.Reader,limit_size int,external bool)*Message{
	buff := make([]byte,12)
	_,err := io.ReadFull(conn,buff)
	if err!=nil{
		log.Info("sock read error:",err)
		return nil
	}
	length,seq,cmd,version,flag := ReadHeader(buff)
	if length<0||length>= limit_size{
		log.Info("invalid len:",length)
		return nil
	}
	//0 <= cmd <= 255
	//收到客户端非法消息，断开链接
	if external&&!external_messages[cmd]{
		log.Warning("invalid external message cmd:",Command(cmd))
		return nil
	}
	buff = make([]byte,length)
	_,err = io.ReadFull(conn,buff)
	if err!=nil{
		log.Info("sock read error:",err)
		return nil
	}

	message := new(Message)
	message.Cmd = cmd
	message.seq = seq
	message.Version = version
	message.Flag = flag
	if !message.FromData(buff){
		log.Warning("parse error:%d %d %d %d %s",cmd,seq,version,flag,
			hex.EncodeToString(buff))
		return nil
	}
	return message
}

func ReceiveMessage(conn io.Reader)*Message{
	return ReceiveLimitMessage(conn,32*1024,false)
}
// 接受客户端消息
func ReceiveClientMessage(conn io.Reader)*Message{
	return ReceiveLimitMessage(conn,32*1024,true)
}
//消息大小限制在32M
func ReceiveStorageSyncMessage(conn io.Reader)*Message{
	return ReceiveLimitMessage(conn,1024*1024,false)
}
//消息带下限制在1M
func ReceiveStorageMessage(conn io.Reader)*Message{
	return ReceiveLimitMessage(conn,1024*1024,false)
}



