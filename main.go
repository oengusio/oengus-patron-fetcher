package main

import (
	"log"
	"net/http"
)

func main() {
	/*patrons, _ := FetchPatrons()

	fmt.Println("There are", len(patrons.Data), "patrons")

	// for {key}, {value} := range {list}
	for _, patron := range patrons.Data {
		fmt.Println("Patron", patron.Attributes.FullName)
	}*/

	mux := http.NewServeMux()

	// Using /pagetier here becuase "/" would catch all routes sent to it
	mux.HandleFunc("/pagetier", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte("{\"patrons\": []}"))
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}