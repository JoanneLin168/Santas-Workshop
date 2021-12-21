package util

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// ChildInSlice - checks if a child is in the slice. Note: compares names - might need to make a comparison function for Child
func ChildInSlice(child Child, children []Child) bool {
	for c := range(children) {
		if child.Name == children[c].Name {
			return true
		}
	}
	return false
}

func RemoveChild(i int, xs []Child) []Child {
	return append(xs[:i], xs[i+1:]...)
}

func MaxInt(xs []int) int {
	max := xs[0]
	for i := range xs {
		x := xs[i]
		if x > max {
			max = x
		}
	}
	return max
}

func MinInt(xs []int) int {
	min := xs[0]
	for i := range xs {
		x := xs[i]
		if x < min {
			min = x
		}
	}
	return min
}

func MaxFloat64(xs []float64) float64 {
	max := xs[0]
	for i := range xs {
		x := xs[i]
		if x > max {
			max = x
		}
	}
	return max
}

func MinFloat64(xs []float64) float64 {
	min := xs[0]
	for i := range xs {
		x := xs[i]
		if x < min {
			min = x
		}
	}
	return min
}

// convertWishList - converts a string-formatted list into a slice of presents
func convertWishlist(str string) []Present {
	// remove [] on either side of the string input
	presents := []Present{}
	str = str[1:len(str)-1]
	strSlice := strings.Split(str, ";")
	for s := range strSlice {
		strPresent := strSlice[s]
		x, err := strconv.Atoi(strPresent)
		Check(err)

		if x > int(Lego) { // Note: may change depending on if you change the presents enum
			panic("Invalid present type!")
		}
		presents = append(presents, Present{PresentType(x)})
	}

	return presents
}

// REFERENCE: https://ankurraina.medium.com/reading-a-simple-csv-in-go-36d7a269cecd
// ConvertCSV - converts the csv file to a slice
func ConvertCSV(filename string) []Child {
	results := []Child{}
	csvfile, err := os.Open(filename)
	Check(err)
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// TODO: consider changing this depending on whether you want to continue adding children or just returning an empty set
		if len(record) < 5 {
			continue
		}

		// convert each value from csv
		name := record[0]
		behaviour, err := strconv.Atoi(record[1])
		Check(err)
		if behaviour > int(Bad) { // Note: may change depending on if you change the behaviour enum
			panic("Invalid behaviour!")
		}
		x, err := strconv.Atoi(record[2])
		Check(err)
		y, err := strconv.Atoi(record[3])
		Check(err)
		address := Address{name, x, y}
		wishlist := convertWishlist(record[4])

		child := Child{
			name, BehaviourType(behaviour), address, wishlist, []Present{},
		}
		results = append(results, child)
	}


	// Check to see if there are any invalid children in the list
	checkValidity(results)

	return results
}

func checkValidity(children []Child) {
	// check if there are duplicate names/addresses
	for i := 0; i < len(children); i++ {
		childA := children[i]
		for j := i+1; j < len(children); j++ {
			childB := children[j]
			if childA.Name == childB.Name {
				panic("Cannot have duplicate names!")
			}
			if childA.Address.X == childB.Address.X && childA.Address.Y == childB.Address.Y {
				panic("Cannot have duplicate addresses!")
			}
		}
	}
}
