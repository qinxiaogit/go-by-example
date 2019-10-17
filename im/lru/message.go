package lru

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
)

const MSG_AUTH_STATUS  =  3
const MSG_IM = 4
const MSG_ACK = 5
const MSG_RST = 6

const MSG_GROUP_NOTIFICATION  = 7
const MSG_GROUP_IM = 8

const MSG_PING = 13
const MSG_PONG = 14
const MSG_AUTH_TOKEN = 15

const MSG_RT = 17
const MSG_ENTER_ROOM = 18
const MSG_LEAVE_ROOM = 19
const MSG_ROOM_IM =20

const MSG_SYSTEM = 21

const MSG_UNREAD_COUNT = 22
const MSG_CUSTOMER_SERVICE_ = 23

//persistent
const MSG_CUSTOMER = 24 //顾客=》顾客
const MSG_CUSTOMER_SUPPORT = 25 //客户到客服

const MSG_SYNC = 26 //同步消息-客户端-》服务端
const MSG_SYNC_BEGIN = 27
const MSG_SYNC_END = 28
//通知客户端有新消息
const MSG_SYNC_NOTIFY = 29

//客户端->服务端
const MSG_SYNC_GROUP = 30 //同步超级群消息
//服务端到客服端
const MSG_SYNC_GROUP_BEGIN = 31
const MSG_SYNC_GROUP_END = 32
//通知客户端有新消息
const MSG_SYNC_GROUP_NOTIFY = 33

//客服端->服务端 ，更新服务器的syncKey
const MSG_SYNC_KEY = 34
const MSG_GROUP_SYNC_KEY  = 35

//系统通知消息
const MSG_NOTIFICATION  =  36
//消息的meta信息
const MSG_METADATA = 37
const MSG_VOIP_CONTROL  = 64

//消息标志
//文本消息 c<->s
const MESSAGE_FLAG_TEXT = 0x01
//消息不持久化
const MESSAGE_FLAG_UNPERSISTENT = 0x02
//群组消息
const MESSAGE_FLAG_GROUP = 0x04
//离线消息由当前登录的用户在当前设备发出
const MESSAGE_F  = 0x08
//消息由服务器主动推送到客户端
const MESSAGE_FLAG_PUSH = 0x10
//超级群消息 c<- s
const MESSAGE_FLAG_SUPER_GROUP  = 0x20

const ACK_SUCCESS = 0
const ACK_NOT_MY_FRIEND = 1
const ACK_NOT_YOUR_FRUEND = 2
const ACK_IN_YOUR_BLACKLIST = 3

const ACK_NOT_GROUP_MEMBER = 64

func init(){
	message_creators[]
}

type Command int

func (cmd Command)String()string{
	c:= int(cmd)
	if desc,ok := message_descriptions[c];ok{
		return desc
	}
	return fmt.Sprintf("%d",c)
}

type IMessage interface {
	ToData()[]byte
	FromData(buff []byte)bool
}

type IVersionMessage interface {
	ToData(version int)[]byte
	FromData(version int,buff[]byte)bool
}

type Message struct {
	cmd int
	seq int
	version int
	flag int

	body interface{}
	meta *MetaData // non searialize
}

func (message *Message)ToData()[]byte{
	if message.body !=nil{
		if m,ok := message.body.(IMessage);ok{
			return m.ToData()
		}
		if m,ok := message.body.(IVersionMessage);ok{
			return m.ToData(message.version)
		}
	}
	return nil
}
func (message *Message)FromData(buff []byte)bool{
	cmd := message.cmd
	if creator,ok := message_creators[cmd];ok{
		c:= creator()
		r:= c.FromData(buff)
		message.body = c
		return r
	}
	if creator,ok := vmessage_creators[cmd];ok{
		c := creator()
		r := c.FromData(message.version,buff)
		message.body = c
		return r
	}
	return len(buff) == 0
}
//保存在磁盘中但不再需要处理消息
type IgnoreMessage struct {
}

func (ignore *IgnoreMessage)ToData()[]byte{
	return nil
}
func (ignore *IgnoreMessage)FromData(buff []byte)bool{
	return true
}

type AuthenticationToken struct {
	token string
	platform_id int8
	device_id string
}

func (auth *AuthenticationToken)ToData()[]byte{
	var l int8

	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.BigEndian,auth.platform_id)

	l = int8(len(auth.token))
	buffer.Write([]byte(auth.token))

	l = int8(len(auth.device_id))
	binary.Write(buffer,binary.BigEndian,l)
	buffer.Write([]byte(auth.device_id))
	return buffer.Bytes()
}

func (auth *AuthenticationToken)FromData(buff []byte)bool{
	var l int8
	if len(buff)<3{
		return false
	}
	auth.platform_id = int8(buff[0])
	buffer := bytes.NewBuffer(buff[1:])

	binary.Read(buffer,binary.BigEndian,&l)
	if int(l)>buffer.Len()||int(l)<0{
		return false
	}
	token := make([]byte,1)
	buffer.Read(token)

	binary.Read(buffer,binary.BigEndian,&l)
	if int(l)>buffer.Len()||int(l)<0{
		return false
	}
	device_id :=  make([]byte,l)
	buffer.Read(device_id)

	auth.token = string(token)
	auth.device_id = string(device_id)
	return true
}

type AuthenticationStatus struct {
	status int32
}
func (auth *AuthenticationStatus)ToData()[]byte{
	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.BigEndian,auth.status)
	buf := buffer.Bytes()
	return buf
}

func (auth *AuthenticationStatus)FromData(buff []byte)bool{
	if len(buff)<4{
		return false
	}
	buffer:= bytes.NewBuffer(buff)
	binary.Read(buffer,binary.BigEndian,&auth.status)
	return true
}



type MetaData struct {
	sync_key int64
	prev_sync_key int64
}

