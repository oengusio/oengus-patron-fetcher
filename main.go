package main

import (
    "encoding/json"
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

    // Using /patrons here becuase "/" would catch all routes sent to it
    mux.HandleFunc("/patrons", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        output := &PatronOutput{
            Patrons: make([]PatronDisplay, 0),
        }

        theJson, _ := json.Marshal(output)

        w.Write(theJson)
    })

    log.Println("Listening to port 9000")

    log.Fatal(http.ListenAndServe(":9000", mux))
}
