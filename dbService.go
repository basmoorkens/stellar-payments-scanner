package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func createIncomingTxRecord(pRecord paymentRecord) {
	db, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("Inserting incoming payment record for " + pRecord.paymentID)
	defer db.Close()
	stmnt, err := db.Prepare("INSERT INTO lsa_payment_records(type, trade_id, amount, asset, asset_issuer_id, status, created_date) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = stmnt.Exec("INCOMING_PAYMENT", pRecord.tradeID, pRecord.amount, pRecord.asset, pRecord.issuer, "NEW", time.Now)
	if err != nil {
		panic(err)
	}
}
