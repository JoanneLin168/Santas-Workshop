package test

import (
	"net/rpc"
	"os"
	"reflect"
	"testing"
	"workshop/util"
)

func TestInput(t *testing.T) {
	os.Stdout = nil

	// RPC dial to server - can only be tested locally due to IP address and time
	server := "127.0.0.1:8080"
	c, err := rpc.Dial("tcp", server)
	util.Check(err)
	defer c.Close()

	// Test that valid csv files create the correct slice of children
	valid := []rune{'A', 'B', 'C', 'D', 'E'}
	testSets := GetSets([]rune{'A', 'B', 'C', 'D', 'E'})
	for i := range valid {
		testName := "Set:" + string(rune(65+i))
		t.Run(testName, func(t *testing.T) {
			filename := "check/valid/input" + string(valid[i]) + ".csv"
			results := util.ConvertCSV(filename)
			expected := testSets[i]
			if !reflect.DeepEqual(results, expected) {
				t.Errorf("Incorrect CSV conversion. Expected %v but got %v", expected, results)
			}
		})
	}

	// Test that invalid csv files will panic
	invalid := []rune{'X', 'Y', 'Z'}
	for i := range invalid {
		testName := "Set:" + string(rune(88+i))
		t.Run(testName, func(t *testing.T) {
			filename := "check/invalid/input" + string(invalid[i]) + ".csv"
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code incorrectly processed an invalid CSV without panicking")
				}
			}()
			util.ConvertCSV(filename)
		})
	}
}
