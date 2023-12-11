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

	SubscriptionGroupId = "21355441"
)

func VerifyReceipt(originTransactionId string) (expiredDate int64, err error) {
	return GetALLSubscriptionStatuses(originTransactionId)
}

func GetALLSubscriptionStatuses(originTransactionId string) (expiredDate int64, err error) {
	authKey, err := os.ReadFile(AccountPrivateKeyFile)
	if err != nil {
		return 0, err
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
	response, err := a.GetALLSubscriptionStatuses(ctx, originTransactionId)
	if err != nil {
		c.Sandbox = true // retry with sandbox
		a = appstore.NewStoreClient(c)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer func() {
			cancel()
		}()

		response, err = a.GetALLSubscriptionStatuses(ctx, originTransactionId)
		if err != nil {
			fmt.Printf("fail get transactions due to: %v\n", err)
			return 0, err
		}
	}

	if len(response.Data) < 1 {
		return 0, fmt.Errorf("fail since data in all subscription status response is empty")
	}

	for _, data := range response.Data {
		if data.SubscriptionGroupIdentifier != SubscriptionGroupId {
			continue
		}

		if len(data.LastTransactions) < 1 {
			return 0, fmt.Errorf("fail since data.LastTransactions in all subscription status response is empty")
		}

		if data.LastTransactions[0].Status != 1 {
			return 0, fmt.Errorf("fail since last transaction status is %d in all subscription status is empty", data.LastTransactions[0].Status)
		}

		signedTransaction := data.LastTransactions[0].SignedTransactionInfo
		transaction, err := a.ParseSignedTransactions([]string{signedTransaction})
		if err != nil {
			fmt.Printf("fail parse signed transaction: %v due to: %v\n", signedTransaction, err)
			return 0, err
		}
		fmt.Printf("GetALLSubscriptionStatuses returns transaction: %+v \n", transaction)

		return transaction[0].ExpiresDate, nil
	}

	err = fmt.Errorf("the Data in response is empty")
	return 0, err
}

func GetTransactionInfo(transactionId string) (expiredDate int64, err error) {
	authKey, err := os.ReadFile(AccountPrivateKeyFile)
	if err != nil {
		return 0, err
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
			return 0, err
		}
	}

	transactions, err := a.ParseSignedTransactions([]string{response.SignedTransactionInfo})
	if err != nil {
		fmt.Printf("fail parse transactions due to: %v\n", err)
		return 0, err
	}
	fmt.Printf("GetTransactionInfo returns transactions: %+v \n", transactions)

	if transactions[0].TransactionID == transactionId && transactions[0].ExpiresDate > time.Now().UnixNano()/1e6 {
		// the transaction is valid
		return transactions[0].ExpiresDate, nil
	}

	err = fmt.Errorf("transaction not match or expired, expect: %v but got %v, expired date is: %v",
		transactions[0].TransactionID, transactionId, transactions[0].ExpiresDate)

	return 0, err
}
