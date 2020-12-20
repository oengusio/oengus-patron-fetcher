package main

import "fmt"

func main() {
	patrons, _ := FetchPatrons()

	fmt.Println("There are", len(patrons.Data), "patrons")

	// for {key}, {value} := range {list}
	for _, patron := range patrons.Data {
		fmt.Println("Patron", patron.Attributes.FullName)
	}
}