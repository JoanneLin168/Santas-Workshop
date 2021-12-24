package test

import (
	"net/rpc"
	"os"
	"testing"
	"workshop/client"
	"workshop/util"
)

func TestDistances(t *testing.T) {
	os.Stdout = nil

	// RPC dial to server - can only be tested locally due to IP address and time
	server := "127.0.0.1:8080"
	c, err := rpc.Dial("tcp", server)
	util.Check(err)
	defer c.Close()

	testSets := GetSets([]rune{'A', 'B', 'C'})
	for i := range testSets {
		testName := "Set:"+string(rune(65+i))
		t.Run(testName, func(t *testing.T) {
			children := testSets[i]
			results := []util.Child{}
			route := []util.Address{}
			client.Run(0, c, children, &results, &route)

			// Note: +2 is for Santa start and Santa end
			switch i {
			case 0:
				if len(route) > 0 {
					t.Errorf("Unexpected length of route for 0 children: %d > 0", len(route))
				}
			case 1:
				if len(route) != 1+2 {
					t.Errorf("Unexpected length of route for 1 child: %d != 1", len(route))
				} else {
					expected := []util.Address{
						{"Santa", 0, 0}, {"Alice", 1, 1}, {"Santa", 0, 0},
					}
					for i := range route {
						if route[i] != expected[i] {
							t.Errorf("Incorrect route for Set B. %v != %v", route, expected)
						}
					}
				}
			case 2:
				if len(route) != 3+2 {
					t.Errorf("Unexpected length of route for 3 children: %d != 5", len(route))
				} else {
					expected := []util.Address{
						{"Santa", 0, 0},
						{"Alice", 1, 1},
						{"Charlie", -1, -1},
						{"Bob", 0, 5},
						{"Santa", 0, 0},
					}
					for i := range route {
						if route[i] != expected[i] {
							t.Errorf("Incorrect route for Set B. %v != %v", route, expected)
							break
						}
					}
				}
			}
		})
	}
}
