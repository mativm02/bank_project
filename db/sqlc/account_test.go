package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAcco(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Matias",
		Balance:  100,
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Owner, account.Owner)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
