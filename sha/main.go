package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func main() {
	s:= "hello golang sha"
	h:=sha1.New()
	h.Write([]byte(s))

	bs:= h.Sum(nil)
	fmt.Println(s)
	fmt.Printf("%x\n",bs)
	//fmt.Printf("%x\n",h)

	data := "abc123!?$*&()'-=@~"

	enBs:=base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(enBs)
	fmt.Println(base64.StdEncoding.DecodeString(enBs))

	uEnc := base64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)
	uDec, _ := base64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}
