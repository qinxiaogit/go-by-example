package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	dateCmd := exec.Command("date")

	dateOutput,err:= dateCmd.Output()

	if err!=nil{
		panic(err)
	}
	fmt.Println(string(dateOutput))

	grepCmd := exec.Command("grep", "hello")

	grepIn, _ := grepCmd.StdinPipe()
	grepOut ,_:= grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep\ngoodbye grep"))

	grepIn.Close()
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()

	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	//lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	//lsOut, err := lsCmd.Output()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("> ls -a -l -h")
	//fmt.Println(string(lsOut))

	//binary,lookError := exec.LookPath("ls")
	//if lookError!=nil{
	//	panic(lookError)
	//}
	//args := []string{"ls", "-a", "-l", "-h"}
	//
	//env := os.Environ()
	//execError:= syscall.Exec(binary,args,env)
	//if execError != nil {
	//	panic(execError)
	//}

	//sigs:= make(chan os.Signal,1)
	//done:= make(chan bool,1)
	//
	//signal.Notify(sigs,syscall.SIGINT,syscall.SIGTERM)
	//
	//go func() {
	//	sig:=<-sigs
	//	fmt.Println("test：",sig)
	//	done<-true
	//
	//}()
	//fmt.Println("awaiting signal")
	//<-done
	//fmt.Println("exiting")

	defer println("hello")
	os.Exit(3)
}
