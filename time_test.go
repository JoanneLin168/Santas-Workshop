package main

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
	children := []util.Child{
		{"Alice",
			util.Good,
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Bob",
			util.Bad,
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Charlie",
			util.Good,
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
	}

	results := []util.Child{}
	var th util.TimeHandler
	th.SetStartTime()
	client.Run(c, children, &results)

	duration := th.GetTime()
	if duration < (6 * time.Second) {
		t.Errorf("Duration too fast: %d", duration)
	} else if duration > (7 * time.Second) {
		t.Errorf("Duration too slow: %d", duration)
	}

}