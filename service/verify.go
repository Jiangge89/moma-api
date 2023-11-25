package service

import (
	"context"
	"fmt"
	"github.com/richzw/appstore"
	"os"
	"time"
)

const (
	AccountPrivateKeyFile = "./service/SubscriptionKey_M3F27D733S.p8"

	AccountPrivateKeyId = "M3F27D733S"
	KeyIssuer           = "c61e84ec-ec57-4fb7-8078-0b27beb6ec9f"

	BundleId = "duftee-moma-free"
)

func VerifyReceipt(transactionId string) error {
	authKey, err := os.ReadFile(AccountPrivateKeyFile)
	if err != nil {
		return err
	}

	c := &appstore.StoreConfig{
		KeyContent: authKey,
		KeyID:      AccountPrivateKeyId,
		BundleID:   BundleId,
		Issuer:     KeyIssuer,
		Sandbox:    true,
	}
	a := appstore.NewStoreClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		cancel()
	}()
	response, err := a.GetTransactionInfo(ctx, transactionId)
	if err != nil {
		c.Sandbox = true // retry with sandbox
		a = appstore.NewStoreClient(c)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer func() {
			cancel()
		}()

		response, err = a.GetTransactionInfo(ctx, transactionId)
		if err != nil {
			fmt.Printf("fail get transactions due to: %v\n", err)
			return err
		}
	}

	transactions, err := a.ParseSignedTransactions([]string{response.SignedTransactionInfo})
	if err != nil {
		fmt.Printf("fail parse transactions due to: %v\n", err)
		return err
	}
	fmt.Printf("GetTransactionInfo returns the first of transactions: %+v \n", *transactions[0])

	if transactions[0].TransactionID == transactionId && transactions[0].ExpiresDate > time.Now().UnixNano()/1e6 {
		// the transaction is valid
		return nil
	}

	return fmt.Errorf("transaction not match or expired, expect: %v but got %v, expired date is: %v",
		transactions[0].TransactionID, transactionId, transactions[0].ExpiresDate)
}
