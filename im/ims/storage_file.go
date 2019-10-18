package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/qinxiaogit/go-by-example/im/lru"
	log "github.com/golang/glog"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const HEADER_SIZE  =  32
const MAGIC = 0x494d494d
const F_VERSION = 1<< 16

const BLOCK_SIZE = 128*1024*1024
const LRU_SIZE  = 128

type StorageFile struct {
	root string
	mutex sync.Mutex

	dirty bool 		//write file dirty
	block_NO int    //write file block no
	file  *os.File  //write
	files *lru.Cache//read ,block files

	last_id	int64 //peer &group_message_index记录的最大消息id
	last_saved_id int64//索引文件中最大的消息id
}

func onFileEvicted(key lru.Key,value interface{}){
	f:= value.(os.File)
	defer f.Close()
}

func NewStorageFile(root string)*StorageFile{
	storage := new(StorageFile)
	storage.files = lru.New(LRU_SIZE)
	storage.files.OnEvicted = onFileEvicted

	//find the last block file
	pattern:= fmt.Sprintf("%s/message_",storage.root)
	files,_:=filepath.Glob(pattern)
	block_NO := 0 //begin from 0
	for _,f := range files{
		base := filepath.Base(f)
		if strings.HasPrefix(base,"message_"){
			if checkFile(f){
				log.Fatal("check file failure")
			}else {
				log.Infof("check file pass:%s",f)
			}
			b,err:= strconv.ParseInt(base[8:],10,64)
			if err!= nil{
				log.Fatal("invalid message file:",f)
			}
			if int(b)>block_NO{
				block_NO = int(b)
			}
		}
	}
	storage.openWriteFile(block_NO)
	return storage
}

func checkFile(file_path string )bool{
	file,err := os.Open(file_path)
	if err!=nil{
		log.Fatal("open file:",err)
	}
	file_size,err := file.Seek(0,os.SEEK_END)
	if err!=nil{
		log.Fatal("seek file")
	}
	if file_size == HEADER_SIZE{
		return true
	}
	if file_size<HEADER_SIZE{
		return false
	}
	_,err = file.Seek(file_size-4,os.SEEK_SET)
	if err!=nil{
		log.Fatal("see file")
	}
	mf := make([]byte,4)
	n,err:= file.Read(mf)
	if err!=nil||n!=4{
		log.Fatal("read file err:",err)
	}
	buffer:= bytes.NewBuffer(mf)
	var m int32
	binary.Read(buffer,binary.BigEndian,&m)
	return int(m) == MAGIC
}
// open write file
func (storage *StorageFile)openWriteFile(block_NO int){
	path := fmt.Sprintf("%s/message_%d",storage.root,block_NO)
	log.Info("open/create message file path:",path)
	file,err := os.OpenFile(path,os.O_RDWR|os.O_APPEND|os.O_CREATE,0644)
	if err!=nil{
		log.Fatal("open file:",err)
	}
	file_size,err := file.Seek(0,os.SEEK_END)
	if err!=nil{
		log.Fatal("seek file")
	}
	if file_size < HEADER_SIZE&&file_size>0{
		log.Info("file header is't complete")
		err = file.Truncate(0)
		if err!=nil{
			log.Fatal("truncate file")
		}
		file_size = 0
	}
	if file_size == 0{
		storage.WriteHeader(file)
	}
	storage.file = file
	storage.block_NO = block_NO
	storage.dirty = false
}
// read file head
func (storage *Storage)openReadFile(block_NO int)*os.File{
	//open file readonly mode
	path := fmt.Sprintf("%s/message_%d", storage.root, block_NO)
	log.Info("open message block file path:",path)
	file,err:= os.Open(path)
	if err!=nil{
		if os.IsNotExist(err){
			log.Infof("message block file:%s nonexist", path)
			return nil
		} else {
			log.Fatal(err)
		}
	}
	file_size,err := file.Seek(0,os.SEEK_END)
	if err!=nil{
		log.Fatal("seek file")
	}
	if file_size<HEADER_SIZE&&file_size>0{
		if err != nil {
			log.Fatal("file header is't complete")
		}
	}
	return file
}
// 写文件头信息
func (storage *StorageFile)WriteHeader(file *os.File){
	var m int32 = MAGIC
	err := binary.Write(file,binary.BigEndian,m)
	if err!=nil{
		log.Fatalln(err)
	}
	var v int32 = F_VERSION
	err= binary.Write(file,binary.BigEndian,v)
	if err!=nil{
		log.Fatalln(err)
	}
	pad := make([]byte,HEADER_SIZE-8)
	n,err:= file.Write(pad)
	if err!=nil||n!=(HEADER_SIZE-8){
		log.Fatalln(err)
	}
}
//读文件头信息
func (Storage *StorageFile)ReadHeader(file *os.File)(magic int,version int){
	header := make([]byte,HEADER_SIZE)
	n,err := file.Read(header)
	if err!=nil||n!= HEADER_SIZE{
		return
	}
	buffer := bytes.NewBuffer(header)
	var m,v int32
	binary.Read(buffer,binary.BigEndian,&m)
	binary.Read(buffer,binary.BigEndian,&v)
	magic = int(m)
	version = int(v)
	return
}


func (storage *StorageFile) getMsgId(block_NO int, offset int) int64 {
	return int64(block_NO)*BLOCK_SIZE + int64(offset)
}

func (storage *StorageFile) getBlockNO(msg_id int64) int {
	return int(msg_id/BLOCK_SIZE)
}

func (storage *StorageFile) getBlockOffset(msg_id int64) int {
	return int(msg_id%BLOCK_SIZE)
}

func (storage *StorageFile) getFile(block_NO int) *os.File {
	v, ok := storage.files.Get(block_NO)
	if ok {
		return v.(*os.File)
	}
	file := storage.openReadFile(block_NO)
	if file == nil {
		return nil
	}

	storage.files.Add(block_NO, file)
	return file
}

//写消息数据
func (storage *StorageFile)WriteMessage(file io.Writer , msg *lru.Message){
	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.BigEndian,int32(MAGIC))
	lru.WriteMessage(buffer,msg)
	binary.Write(buffer,binary.BigEndian,int32(MAGIC))
	buf := buffer.Bytes()
	n,err := file.Write(buf)
	if err!=nil{
		log.Fatal("file write err:",err)
	}
	if n!=len(buf){
		log.Fatal("file write size",len(buf)," nwrite ",n)
	}
}
//save without lock
func (storage *StorageFile)saveMessage(msg *lru.Message) int64{
	msgid,err := storage.file.Seek(0,os.SEEK_END)
	if err!=nil{
		log.Fatalln(err)
	}
	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.BigEndian,int32(MAGIC))
	lru.WriteMessage(buffer,msg)
	binary.Write(buffer,binary.BigEndian,int32(MAGIC))
	buf := buffer.Bytes()

	if msgid+int64(len(buf)) > BLOCK_SIZE{
		err = storage.file.Sync()
		if err!=nil{
			log.Fatalln("sync storage file")
		}
		storage.file.Close()
		storage.openWriteFile(storage.block_NO+1)
		msgid,err = storage.file.Seek(0,os.SEEK_END)
		if err!=nil{
			log.Fatalln(err)
		}
	}
	if msgid+int64(len(buf))>BLOCK_SIZE{
		log.Fatalln("message size:",len(buf))
	}
	n ,err := storage.file.Write(buf)
	if err!=nil{
		log.Fatal("file write err:" ,err)
	}
	if n!=len(buf){
		log.Fatal("file write size:",len(buf)," nwrite:",n)
	}
	storage.dirty =true

	msgid = int64(storage.block_NO)*BLOCK_SIZE+msgid
	master.ewt <- &lru.EMessage{Msgid:msgid,Msg:msg}
	//log.Info("save message:", Command(msg.cmd), " ", msgid)
	return msgid
}

func (storage *StorageFile) SaveMessage(msg *lru.Message)int64{
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	return storage.saveMessage(msg)
}

func (storage *StorageFile)Flush(){
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	if storage.file != nil && storage.dirty{
		err := storage.file.Sync()
		if err!=nil{
			log.Fatal("sync err:",err)
		}
		storage.dirty = false
		log.Info("sync storage file success")
	}
}