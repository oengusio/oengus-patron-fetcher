package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type People struct {
	Number int `json:"number"`
	People []Person `json:"people"`
}

type Person struct {
	Craft string `json:"craft"`
	Name string `json:"name"`
}

func main() {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	url := "http://api.open-notify.org/astros.json"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	fmt.Println("Response status:", res.Status)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}


	people := People{}

	err2 := json.Unmarshal(body, &people)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	fmt.Println(people.Number)

	for _, person := range people.People {
		fmt.Println(person.Name, "is in craft", person.Craft)
	}
}