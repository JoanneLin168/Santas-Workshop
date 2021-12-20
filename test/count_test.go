package test

import (
	"flag"
	"net/rpc"
	"os"
	"reflect"
	"testing"
	"workshop/client"
	"workshop/util"
)

func TestCount(t *testing.T) {
	os.Stdout = nil

	// RPC dial to server
	server := flag.String("server", "127.0.0.1:8080", "IP:port string to connect to")
	c, err := rpc.Dial("tcp", *server)
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
	client.Run(c, children, results)

	for c := range(results) {
		child := results[c]
		switch child.Behaviour {
		case util.Good:
			if !reflect.DeepEqual(child.Presents, child.WishList) {
				t.Errorf("Incorrect number of presents for %s: %d != %d",
					child.Name, len(child.Presents), len(child.WishList))
			}
		case util.Bad:
			if !(len(child.Presents) == 1 && child.Presents[0].Type == util.Coal) {
				t.Errorf("Incorrect number of presents for %s: %d != %d",
					child.Name, len(child.Presents), 1)
			}
		}

	}
}
