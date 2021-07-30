package sql

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4"
    "log"
    "os"
)

var db *pgx.Conn

func Start() {
    log.Println(os.Getenv("DATABASE_URL"))

    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
        os.Exit(1)
    }

    db = conn
}

func Stop() {
    defer func(db *pgx.Conn, ctx context.Context) {
        err := db.Close(ctx)
        if err != nil {
            log.Println(err)
        }
    }(db, context.Background())
}

func InsertMember(userId string, status string, payAmount int) {
    sql := "INSERT INTO patreon_status(patreon_id, status, pledge_amount) VALUES($1, $2, $3);"

    _, err := db.Query(context.Background(), sql, userId, status, payAmount)

    if err != nil {
        fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
        os.Exit(1)
    }
}

func UpdateMember(userId string, status string, payAmount int) {
    sql := "UPDATE patreon_status SET status = $2, pledge_amount = $3 WHERE patreon_id = $1;"

    _, err := db.Query(context.Background(), sql, userId, status, payAmount)

    if err != nil {
        fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
        os.Exit(1)
    }
}

func DeleteMember(userId string) {
    sql := "DELETE FROM patreon_status WHERE patreon_id = $1;"

    _, err := db.Query(context.Background(), sql, userId)

    if err != nil {
        fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
        os.Exit(1)
    }
}


