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
		Sandbox:    false,
	}
	a := appstore.NewStoreClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		cancel()
	}()
	response, err := a.GetALLSubscriptionStatuses(ctx, originTransactionId)
	if err != nil {
		fmt.Printf("[Err] fail get transactions for id: %v in production mode due to: %v, will try to request in sandbox mode then.\n", originTransactionId, err)

		c.Sandbox = true // retry with sandbox
		a = appstore.NewStoreClient(c)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer func() {
			cancel()
		}()

		var sandboxErr error
		response, sandboxErr = a.GetALLSubscriptionStatuses(ctx, originTransactionId)
		if sandboxErr != nil {
			fmt.Printf("[Err] fail get transactions for id: %v in sandbox mode due to: %v, will return the production request err \n", originTransactionId, sandboxErr)
			return 0, err
		}
	}

	if len(response.Data) < 1 {
		return 0, fmt.Errorf("fail since data in GetALLSubscriptionStatuses response for transaction id>: %v is empty", originTransactionId)
	}

	for _, data := range response.Data {
		if data.SubscriptionGroupIdentifier != SubscriptionGroupId {
			continue
		}

		if len(data.LastTransactions) < 1 {
			return 0, fmt.Errorf("[Err] fail since data.LastTransactions in GetALLSubscriptionStatuses response is empty")
		}

		fmt.Printf("[INFO] last transaction in GetALLSubscriptionStatuses response is %v \n", data.LastTransactions[0])

		if data.LastTransactions[0].Status != 1 {
			return 0, fmt.Errorf("[Err] fail since last transaction status is %d in GetALLSubscriptionStatuses response", data.LastTransactions[0].Status)
		}

		signedTransaction := data.LastTransactions[0].SignedTransactionInfo
		transaction, err := a.ParseSignedTransactions([]string{signedTransaction})
		if err != nil {
			fmt.Printf("[Err] fail parse signed transaction: %v due to: %v\n", signedTransaction, err)
			return 0, err
		}

		fmt.Printf("[INFO] GetALLSubscriptionStatuses for id: %v returns %v transactions: \n", originTransactionId, len(transaction))
		for idx, t := range transaction {
			fmt.Printf("%v transaction for id: %v is: %+v \n", idx, originTransactionId, t)
		}

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
			fmt.Printf("fail get transactions for id: %v due to: %v\n", transactionId, err)
			return 0, err
		}
	}

	transactions, err := a.ParseSignedTransactions([]string{response.SignedTransactionInfo})
	if err != nil {
		fmt.Printf("fail parse transactions: %v due to: %v\n", response.SignedTransactionInfo, err)
		return 0, err
	}
	fmt.Printf("GetTransactionInfo for id: %v returns transactions: %+v \n", transactionId, transactions)

	if transactions[0].TransactionID == transactionId && transactions[0].ExpiresDate > time.Now().UnixNano()/1e6 {
		// the transaction is valid
		return transactions[0].ExpiresDate, nil
	}

	err = fmt.Errorf("transaction not match or expired, expect: %v but got %v, expired date is: %v",
		transactions[0].TransactionID, transactionId, transactions[0].ExpiresDate)

	return 0, err
}
