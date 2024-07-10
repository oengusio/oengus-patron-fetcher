package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io"
	"log"
	"oenugs-patreon/cache"
	"oenugs-patreon/patreon"
	"oenugs-patreon/sql"
	"oenugs-patreon/structs"
	"os"
	"os/signal"
	"strings"
	"time"
)

func UpdatePatrons() {
	log.Println("Updating patrons")

	patrons, err := patreon.FetchPatrons(cache.Tokens)

	// 401 response, refresh the tokens
	if err != nil && err.Error() == "StatusUnauthorized" {
		newTokens, refreshErr := patreon.RefreshToken(cache.Tokens)

		if refreshErr != nil {
			log.Println("Error while refreshing tokens", refreshErr.Error())
			return
		}

		cache.Tokens = newTokens
		patrons, err = patreon.FetchPatrons(cache.Tokens)
	}

	if err != nil {
		log.Println(err)
		return
	}

	newCache := structs.PatronOutput{
		Patrons: make([]structs.PatronDisplay, 0),
	}

	// for {key}, {value} := range {list}
	for i, patron := range patrons.Data {
		attr := patron.Attributes

		// is an active patron that pays $1 or more
		if attr.PatronStatus == "active_patron" && attr.WillPayAmountCents >= 100 {
			userId := patron.Relationships.User.Data.Id
			imageUrl := patrons.Included[i].Attributes.ImageUrl

			newCache.Patrons = append(newCache.Patrons, structs.PatronDisplay{
				Id:       userId,
				Name:     attr.FullName,
				ImageUrl: imageUrl,
			})
		}
	}

	log.Println("Found", len(newCache.Patrons), "patrons")

	cache.PatronCache = newCache

	log.Println("Updating database")
	updatePatronsInDatabase(patrons.Data)
	log.Println("Database update complete")
}

func StartUpdatePatronTimer() {
	go UpdatePatrons()

	// tick every 24 hours to update the patrons
	ticker := time.NewTicker(24 * time.Hour)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		for {
			select {
			case <-ticker.C:
				UpdatePatrons()
			case <-sigint:
				ticker.Stop()
				return
			}
		}
	}()
}

func LoadPatronCredentials() {
	if _, err := os.Stat("/storage/oengus-patreon/patreon-credentials.json"); os.IsNotExist(err) {
		// fetch credentials
		log.Fatal("Failed to load credentials file at \"/storage/oengus-patreon/patreon-credentials.json\"")
	}

	// Open our jsonFile
	jsonFile, err := os.Open("/storage/oengus-patreon/patreon-credentials.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &cache.Tokens)
}

func updatePatronsInDatabase(data []structs.PatreonMembersData) {
	conn := sql.GetConnection()

	defer sql.CloseConnection(conn)

	batch := &pgx.Batch{}

	for _, patron := range data {
		attr := patron.Attributes

		query := "INSERT INTO  patreon_status(patreon_id, status, pledge_amount) VALUES($1, $2, $3) " +
			"ON CONFLICT (patreon_id) DO UPDATE SET status = $2, pledge_amount = $3;"

		userId := patron.Relationships.User.Data.Id
		status := strings.ToUpper(attr.PatronStatus)
		payAmount := attr.WillPayAmountCents

		log.Println(userId, "is an", status, "paying", payAmount, "cents")

		if status == "" {
			// Ignore blank statuses
			continue
		}

		batch.Queue(query, userId, status, payAmount)

		//_, err := conn.Query(context.Background(), query, userId, status, payAmount)
		//
		//if err != nil {
		//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		//}
	}

	br := conn.SendBatch(context.Background(), batch)

	for range batch.Len() {
		_, err := br.Exec()

		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		}
	}

	defer br.Close()
}
