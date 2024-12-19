package gapi

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/daniel-adam-ce/go-bank/db/mock"
	db "github.com/daniel-adam-ce/go-bank/db/sqlc"
	"github.com/daniel-adam-ce/go-bank/pb"
	"github.com/daniel-adam-ce/go-bank/util"
	mockwk "github.com/daniel-adam-ce/go-bank/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actualArg.HashedPassword)

	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArg.HashedPassword
	return reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams)
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (db.User, string) {
	password, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, password)
	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: password,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}, password
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// important (see function comments)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			taskDistributor := mockwk.NewMockTaskDistributor(ctrl)
			tc.buildStubs(store)

			// start test http server and send request
			server := newTestServer(t, store, taskDistributor)
			res, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
