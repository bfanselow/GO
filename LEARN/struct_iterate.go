/***************************************************************************
 struct_iterate.go
 bfanselow 2017-10-31

 Playing around with struct types. Build an array of struct data and then
 iterate over that to generate some data totals

***************************************************************************/
package main

import "fmt"

// Define Person struct
type PersonData struct {
	Name          string
	height_inches int
	weight_lbs    int
}

// Create a slice of *Person* structs
var AllPersonData = []PersonData{
	{"Bill", 72, 170},
	{"Joe", 70, 160},
	{"Bob", 65, 155},
	{"Tom", 71, 175},
	{"Fred", 69, 171},
}

// Return total weight of all people in list
func GetTotalWeight(a_data []PersonData) int {
	total := 0
	for _, elem := range a_data {
		total += elem.weight_lbs
	}
	return total
}

// Return total height of all people in list
func GetTotalHeight(a_data []PersonData) int {
	total := 0
	for _, elem := range a_data {
		total += elem.height_inches
	}
	return total
}

// Return total height and weight of all people in list
func GetTotals(a_data []PersonData) (int, int) {
	total_h := 0
	total_w := 0
	for _, elem := range a_data {
		total_h += elem.height_inches
		total_w += elem.weight_lbs
	}
	return total_h, total_w
}

func main() {
	var tp = len(AllPersonData)
	//var tw = GetTotalWeight(AllPersonData)
	//var th = GetTotalHeight(AllPersonData)
	th, tw := GetTotals(AllPersonData)
	fmt.Printf("Total Weight (%d people): %d\n", tp, tw)
	fmt.Printf("Total Height (%d people): %d\n", tp, th)
}
