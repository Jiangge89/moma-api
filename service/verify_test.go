package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifyReceipt(t *testing.T) {
	var transactionID = "2000000476186502"
	expiredDate, err := VerifyReceipt(transactionID)
	require.Nil(t, err)

	fmt.Println(expiredDate)
}
