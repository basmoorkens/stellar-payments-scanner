package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon/operations"
)

func main() {
	var escrowAcc = "GAEMQ2KLGBAOEWN6VQCF6DDRUMBS327UYGZWOYZCMGR2DMBFOAZGKZUZ"
	pair, err := keypair.Parse(escrowAcc)
	if err != nil {
		log.Fatal(err)
	}
	// Create and fund the address on TestNet, using friendbot
	client := horizonclient.DefaultPublicNetClient

	// Get information about the account we just created
	accountRequest := horizonclient.AccountRequest{AccountID: pair.Address()}
	hAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hAccount.AccountID)
	fmt.Printf("%v", hAccount.Balances)

	txRequest := horizonclient.TransactionRequest{ForAccount: hAccount.AccountID, Cursor: "76068307513434112"}
	transactions, err := client.Transactions(txRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range transactions.Embedded.Records {
		fmt.Println("Account: " + tx.Account)
		fmt.Println("Hash: " + tx.Hash)
		fmt.Println("Memo: " + tx.Memo)
		opRequest := horizonclient.OperationRequest{ForTransaction: tx.Hash}
		ops, err := client.Payments(opRequest)
		if err != nil {
			log.Fatal(err)
		}
		for _, op := range ops.Embedded.Records {
			if op.GetType() == "payment" {
				payment := op.(operations.Payment)
				if payment.To == escrowAcc {
					fmt.Println("Payment to our escrow acc" + payment.To)
					fmt.Print(payment.Amount + " ")
					if payment.Asset.Type == "native" {
						fmt.Println("XLM")
					} else {
						fmt.Println(payment.Asset.Code)
					}
				}
			}
		}
	}

}
