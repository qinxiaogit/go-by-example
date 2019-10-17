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
			if chec
		}
	}

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
		storage.W
	}
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

//写消息数据
func (storage *StorageFile)WriteMessage(file io.Writer , msg *Message){

}