package test

import (
	"net/rpc"
	"os"
	"testing"
	"time"
	"workshop/client"
	"workshop/util"
)

func TestTime(t *testing.T) {
	os.Stdout = nil

	// RPC dial to server - can only be tested locally due to IP address and time
	server := "127.0.0.1:8080"
	c, err := rpc.Dial("tcp", server)
	util.Check(err)
	defer c.Close()

	// TODO: look at the website below to make a table to create different combinations of children to test with client.Run
	// 	https://blog.alexellis.io/golang-writing-unit-tests/

	testSets := GetSets()
	for i := range testSets {
		children := testSets[i]
		results := []util.Child{}
		var th util.TimeHandler
		th.SetStartTime()

		client.Run(c, children, &results)

		// TODO: figure out a way to calculate the time it will take for elves to work, rather than hardcoding it
		//		or get Santa to order the tasks into the most efficient order and send that to the workshop
		duration := th.GetTime()
		if duration < (4 * time.Second) {
			t.Errorf("Duration too fast: %d", duration)
		} else if duration > (5 * time.Second) {
			t.Errorf("Duration too slow: %d", duration)
		}
	}



}