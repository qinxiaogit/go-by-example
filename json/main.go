package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Response1 struct {
	Fruits []string
	Page int
}
type Response2 struct {
	Fruits []string `json:fruits`
	Page int `json:page`
}

func main() {
	boolB,_:=json.Marshal(true)
	fmt.Println("json:",boolB,string(boolB))

	intB, _ := json.Marshal(1)
	fmt.Println("jsonInt:",intB,string(intB))

	floatB ,_:=json.Marshal(2.345)
	fmt.Println("jsonFloat:",floatB,string(floatB))

	strB,_:=json.Marshal("hello world")
	fmt.Println("str: ",strB,string(strB))

	slcB,_:=json.Marshal([]string{"hello","json","golang"})
	fmt.Println("slc",slcB,string(slcB))

	mapB,_:=json.Marshal(map[string]string{"hello":"world","lang":"golang"})
	fmt.Println("map",mapB,string(mapB))

	resq:=&Response1{Page:1,Fruits:[]string{"hello","world"}}

	resq1B,_:= json.Marshal(resq)

	fmt.Println("Response1:",resq1B,string(resq1B))

	res2D := Response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	byt :=[]byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat map[string]interface{}

	if err:=json.Unmarshal(byt,&dat);err!=nil {
		panic(err)
	}
	fmt.Println(dat)
	num:=dat["num"].(float64)

	fmt.Println(num)

	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := &Response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	enc:=json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)


	}
