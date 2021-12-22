package test

import (
	"net/rpc"
	"os"
	"reflect"
	"testing"
	"workshop/client"
	"workshop/util"
)

func TestResults(t *testing.T) {
	os.Stdout = nil

	// RPC dial to server - can only be tested locally due to IP address and time
	server := "127.0.0.1:8080"
	c, err := rpc.Dial("tcp", server)
	util.Check(err)
	defer c.Close()

	testSets := GetSets([]rune{'A', 'B', 'C', 'D', 'E'})
	for i := range testSets {
		testName := "Set:"+string(rune(65+i))
		t.Run(testName, func(t *testing.T) {
			children := testSets[i]
			results := []util.Child{}
			route := []util.Address{}
			client.Run(0, c, children, &results, &route)

			if len(results) != len(children) {
				t.Errorf("Incorrect number of children returned: %d != %d", len(results), len(children))
			}

			for c := range(results) {
				child := results[c]
				switch child.Behaviour {
				case util.Good:
					if !reflect.DeepEqual(child.Presents, child.WishList) {
						t.Errorf("Incorrect presents for good child %s: %d != %d",
							child.Name, len(child.Presents), len(child.WishList))
					}
				case util.Bad:
					if !(len(child.Presents) == 1 && child.Presents[0].Type == util.Coal) {
						t.Errorf("Incorrect presents for bad child %s: %d != %d",
							child.Name, len(child.Presents), 1)
					}
				}

			}
		})

	}
}
