package sql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

func GetConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func CloseConnection(db *pgx.Conn) {
	err := db.Close(context.Background())

	if err != nil {
		log.Println(err)
	}
}

func InsertMember(userId string, status string, payAmount int) {
	sql := "INSERT INTO patreon_status(patreon_id, status, pledge_amount) VALUES($1, $2, $3);"

	db := GetConnection()
	defer CloseConnection(db)

	_, err := db.Query(context.Background(), sql, userId, status, payAmount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}

func UpdateMember(userId string, status string, payAmount int) {
	sql := "UPDATE patreon_status SET status = $2, pledge_amount = $3 WHERE patreon_id = $1;"

	db := GetConnection()
	defer CloseConnection(db)

	_, err := db.Query(context.Background(), sql, userId, status, payAmount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}

func DeleteMember(userId string) {
	db := GetConnection()
	defer CloseConnection(db)

	DeleteMemberConn(userId, db)
}

func DeleteMemberConn(userId string, db *pgx.Conn) {
	sql := "DELETE FROM patreon_status WHERE patreon_id = $1;"

	_, err := db.Query(context.Background(), sql, userId)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}
