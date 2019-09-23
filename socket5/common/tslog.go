package common
//https://blog.csdn.net/u014470361/article/details/81512330
import (
	"fmt"
	"time"
)

type TSLog struct {
}

func (z *TSLog) now()string{
	return time.Now().Format("2016-01-02 15:04:05")
}

func (z *TSLog)log(c string,f string,v ...interface{}){
	out :=fmt.Sprintf(f,v...)
	if c!=""{
		out = fmt.Sprintf("\033[%sm%s %s\033[0m",c,z.now(),out)
	}else{
		out = fmt.Sprintf("%s %s",z.now(),out)
	}
	fmt.Println(out)
}

func (z *TSLog)Log(f string ,v ...interface{}){
	z.log("0",f,v...)
}
func (z *TSLog)Green(f string, v...interface{}){
	z.log("O;32",f,v...)
}
func (z *TSLog)Red(f string,v...interface{}){
	z.log("0;31",f,v...)
}
func (z *TSLog)Gray(f string,v...interface{}){
	z.log("1:30",f,v...)
}

