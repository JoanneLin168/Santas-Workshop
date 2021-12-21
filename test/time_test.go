package test

import (
	"fmt"
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

	testSets := GetSets([]rune{'A', 'B', 'C', 'D', 'E'})
	for i := range testSets {
		testName := "Set:"+string(rune(65+i))
		t.Run(testName, func(t *testing.T) {
			children := testSets[i]
			results := []util.Child{}
			route := []util.Address{}
			var th util.TimeHandler
			th.SetStartTime()
			client.Run(c, children, &results, &route)

			expectedTime := expectedDuration(children)
			if th.GetTime() < expectedTime * time.Second || th.GetTime() > (expectedTime+1) * time.Second {
				t.Errorf("Unexpected time taken: %d != %d", th.GetTime(), expectedTime*time.Second)
			}
		})
	}
}

func expectedDuration(set []util.Child) time.Duration {
	semList := []int{0, 0, 0, 0}

	for c := range set {
		child := set[c]
		m := 0
		for i := 0; i < 4; i++ {
			if semList[i] < semList[m] {
				m = i
			}
		}

		switch child.Behaviour {
		case util.Good:
			semList[m] += len(child.WishList) * 3
		case util.Bad:
			semList[m] += 1
		}
		fmt.Println(semList)
	}

	m := 0
	for i := 0; i < 4; i++ {
		if semList[i] > semList[m] {
			m = i
		}
	}

	return time.Duration(semList[m])
}
