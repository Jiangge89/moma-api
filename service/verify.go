package service

import (
	"context"
	"fmt"
	"github.com/richzw/appstore"
	"os"
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
		Sandbox:    false,
	}
	a := appstore.NewStoreClient(c)
	response, err := a.GetTransactionInfo(context.TODO(), transactionId)
	if err != nil {
		return err
	}

	transactions, err := a.ParseSignedTransactions([]string{response.SignedTransactionInfo})
	if err != nil {
		return err
	}

	if transactions[0].TransactionID == transactionId {
		// the transaction is valid
		return nil
	}

	return fmt.Errorf("transaction not match, expect: %v but got %v", transactions[0].TransactionID, transactionId)
}
