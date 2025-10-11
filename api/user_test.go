package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Doris-Mwito5/simple-bank/internal/db/mock"
	db "github.com/Doris-Mwito5/simple-bank/internal/db/sqlc"
	"github.com/Doris-Mwito5/simple-bank/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.ValidatePassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v,", e.arg, e.password)
}

func eqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//marshall data
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{ // Remove colon to assign to the return variable
		ID:             utils.RandomInt(1, 1000),
		Username:       utils.RandomOwner(),
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	return user, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, expectedUser db.User) {
	data, err := io.ReadAll(body) // Use io.ReadAll instead of ioutil.ReadAll
	require.NoError(t, err)

	// Create a response struct that matches your actual API response
	var response struct {
		ID                int64  `json:"id"`
		Username          string `json:"username"`
		FullName          string `json:"full_name"`
		Email             string `json:"email"`
		PasswordChangedAt string `json:"password_changed_at"`
		CreatedAt         string `json:"created_at"`
	}

	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	require.Equal(t, expectedUser.Username, response.Username)
	require.Equal(t, expectedUser.FullName, response.FullName)
	require.Equal(t, expectedUser.Email, response.Email)
	// Don't compare password fields
}
