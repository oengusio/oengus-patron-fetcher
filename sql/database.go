package sql

import (
    "github.com/go-pg/pg/v10"
    "oenugs-patreon/structs"
)

var db *pg.DB

func Start() {
    db = pg.Connect(&pg.Options{})
}

func InsertMember(status string, userId string, payAmount int) {
    // sql := "INSERT INTO patreon_status(patreon_id, status, pledge_amount) VALUES(?, ?, ?);"

    model := &structs.PatreonStatus{
        PatreonId:    userId,
        Status:       status,
        PledgeAmount: payAmount,
    }

    _, err := db.Model(model).Insert()

    if err != nil {
        panic(err)
    }
}

func UpdateMember(status string, userId string, payAmount int) {
    // sql := "UPDATE patreon_status SET status = ?, pledge_amount = ? WHERE patreon_id = ?;"

    model := &structs.PatreonStatus{
        PatreonId:    userId,
        Status:       status,
        PledgeAmount: payAmount,
    }

    _, err := db.Model(model).WherePK().Update()

    if err != nil {
        panic(err)
    }
}

func DeleteMember(userId string) {
    // sql := "DELETE FROM patreon_status WHERE patreon_id = ?;"

    model := &structs.PatreonStatus{
        PatreonId:    userId,
    }

    _, err := db.Model(model).WherePK().Delete()

    if err != nil {
        panic(err)
    }
}


