package main

import (
	"flag"
	"fmt"
	"net/rpc"
	c "workshop/client"
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
	children := util.ConvertCSV("input.csv")

	results := []util.Child{}
	route := []util.Address{}
	c.Run(client, children, &results, &route)
	fmt.Println(fmt.Sprintf("All %d children have received presents from Santa!", len(results)))
}
