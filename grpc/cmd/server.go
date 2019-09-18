package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var serverCmd =&cobra.Command{
	Use:"server",
	Short:"Run the gRPC hello-world server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err:=recover();err!=nil{
				log.Println("recover err: %v",err)
			}
		}()
		server.run()
	},
}

func init