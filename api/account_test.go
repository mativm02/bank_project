package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/mativm02/bank_system/db/mock"
	db "github.com/mativm02/bank_system/db/sqlc"
	"github.com/mativm02/bank_system/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: -1,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore, functionBody createAccountRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		functionBody  createAccountRequest
	}{
		{
			name: "OK",
			buildStub: func(store *mockdb.MockStore, functionBody createAccountRequest) {
				arg := db.CreateAccountParams{
					Owner:    functionBody.Owner,
					Currency: functionBody.Currency,
				}

				store.EXPECT().CreateAccount(gomock.Any(), arg).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
			functionBody: createAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
		},
		{
			name: "InternalError",
			buildStub: func(store *mockdb.MockStore, functionBody createAccountRequest) {
				arg := db.CreateAccountParams{
					Owner:    functionBody.Owner,
					Currency: functionBody.Currency,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
			functionBody: createAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
		},
		{
			name: "InvalidCurrency",
			buildStub: func(store *mockdb.MockStore, functionBody createAccountRequest) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			functionBody: createAccountRequest{
				Owner:    account.Owner,
				Currency: "invalid",
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store, tc.functionBody)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			body, err := json.Marshal(tc.functionBody)
			if err != nil {
				t.Fatal(err)
			}
			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestListAccounts(t *testing.T) {
	accounts := []db.Account{
		randomAccount(),
		randomAccount(),
	}

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore, functionBody listAccountsRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		functionBody  listAccountsRequest
	}{
		{
			name: "OK",
			buildStub: func(store *mockdb.MockStore, functionBody listAccountsRequest) {
				arg := db.ListAccountsParams{
					Limit:  functionBody.PageSize,
					Offset: (functionBody.PageID - 1) * functionBody.PageSize,
				}
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotAccounts []db.Account
				err = json.Unmarshal(data, &gotAccounts)
				require.NoError(t, err)

				require.Equal(t, accounts, gotAccounts)
			},
			functionBody: listAccountsRequest{
				PageID:   1,
				PageSize: 5,
			},
		},
		{
			name: "InternalError",
			buildStub: func(store *mockdb.MockStore, functionBody listAccountsRequest) {
				arg := db.ListAccountsParams{
					Limit:  functionBody.PageSize,
					Offset: (functionBody.PageID - 1) * functionBody.PageSize,
				}
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
			functionBody: listAccountsRequest{
				PageID:   1,
				PageSize: 5,
			},
		},
		{
			name: "InvalidPageID",
			buildStub: func(store *mockdb.MockStore, functionBody listAccountsRequest) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			functionBody: listAccountsRequest{
				PageID:   -1,
				PageSize: 5,
			},
		},
		{
			name: "InvalidPageSize",
			buildStub: func(store *mockdb.MockStore, functionBody listAccountsRequest) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			functionBody: listAccountsRequest{
				PageID:   1,
				PageSize: 0,
			},
		},
		{
			name: "NotFound",
			buildStub: func(store *mockdb.MockStore, functionBody listAccountsRequest) {
				arg := db.ListAccountsParams{
					Limit:  functionBody.PageSize,
					Offset: (functionBody.PageID - 1) * functionBody.PageSize,
				}
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
			functionBody: listAccountsRequest{
				PageID:   1,
				PageSize: 5,
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store, tc.functionBody)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.functionBody.PageID, tc.functionBody.PageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
