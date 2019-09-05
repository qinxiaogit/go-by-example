package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main(){
	dateCmd := exec.Command("ls","-la")
	dateOut,err:=dateCmd.Output()
	if err!=nil{
		panic(err)
	}
	fmt.Println(string(dateOut))

	binary,lookErr:= exec.LookPath("ls")
	if lookErr!=nil{
		panic(lookErr)
	}
	args:=[]string{"ls","-a","-l","-h"}
	env:=os.Environ()
	execErr:=syscall.Exec(binary,args,env)
	if execErr!=nil{
		panic(execErr)
	}
}
