package tools

import "strings"
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func failOnError(err error,msg string){
	if err != nil{
		log.Fatalf("%s:%s",err,msg)
	}
}
