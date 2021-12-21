package main

import (
	"flag"
	"fmt"
	"net/rpc"
	client2 "workshop/client"
	"workshop/test"
	"workshop/util"
)

func main() {
	// RPC dial to server
	server := flag.String("server", "127.0.0.1:8080", "IP:port string to connect to")
	client, err := rpc.Dial("tcp", *server)
	flag.Parse()
	util.Check(err)
	defer client.Close()

	// temporary list of children, in the future use text file/csv
	children := test.SetC()

	results := []util.Child{}
	client2.Run(client, children, &results)
	fmt.Println(fmt.Sprintf("All %d children have received presents from Santa!", len(results)))
}
