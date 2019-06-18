package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/operations"
)

const depositEvent = "TRADEDIRECT_DEPOSIT_INTO_ESCROW"

type paymentRecord struct {
	tradeID   string
	amount    string
	asset     string
	issuer    string
	paymentID string
}

func main() {
	var escrowAcc = "GAEMQ2KLGBAOEWN6VQCF6DDRUMBS327UYGZWOYZCMGR2DMBFOAZGKZUZ"
	pair, err := keypair.Parse(escrowAcc)
	if err != nil {
		log.Fatal(err)
	}

	client := horizonclient.DefaultPublicNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: pair.Address()}
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}
	for {
		scanHorizon(client, account, escrowAcc)
	}
}

func scanHorizon(client *horizonclient.Client, account hProtocol.Account, escrowAcc string) {
	cur := getLastScannedHorizonCursor()
	lastCur := ""
	txRequest := horizonclient.TransactionRequest{ForAccount: account.AccountID, Cursor: cur}
	transactions, err := client.Transactions(txRequest)
	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("Found " + strconv.Itoa(len(transactions.Embedded.Records)) + " records for cursor " + cur)
	for _, tx := range transactions.Embedded.Records {
		opRequest := horizonclient.OperationRequest{ForTransaction: tx.Hash}
		ops, err := client.Payments(opRequest)
		if err != nil {
			log.Fatal(err)
		}
		logrus.Info("Processing " + strconv.Itoa(len(ops.Embedded.Records)) + " operations for tx hash " + tx.Hash)
		for _, op := range ops.Embedded.Records {
			if op.GetType() == "payment" {
				payment := op.(operations.Payment)
				if payment.To == escrowAcc {
					processRecord(payment, tx)
				}
				lastCur = payment.PT
			}
		}
	}
	if len(transactions.Embedded.Records) != 0 {
		saveCursor(lastCur)
	} else {
		logrus.Info("Nothing to process, sleeping for a minute...")
		time.Sleep(60 * 1000 * time.Millisecond)
	}

}

func processRecord(payment operations.Payment, tx hProtocol.Transaction) {
	logrus.Info("Processing payment " + payment.ID)
	if strings.Contains(tx.Memo, depositEvent) {
		pRecord := paymentRecord{
			strings.TrimPrefix(tx.Memo, depositEvent),
			payment.Amount,
			payment.Asset.Code,
			"",
			payment.ID,
		}
		if payment.Asset.Type != "native" {
			pRecord.issuer = payment.Asset.Issuer
		}
		createIncomingTxRecord(pRecord)
	}

}
