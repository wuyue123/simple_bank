package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"pxsemic.com/simplebank/util"
)

func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

	return user
}

// test create User
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// test get User
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user1.Username, user.Username)
	require.Equal(t, user1.HashedPassword, user.HashedPassword)
	require.Equal(t, user1.FullName, user.FullName)
	require.Equal(t, user1.Email, user.Email)
	require.WithinDuration(t, user1.CreatedAt, user.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullname(t *testing.T) {
	user := createRandomUser(t)
	newFullName := util.RandomOwner()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, newFullName, user.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)

	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user := createRandomUser(t)
	newEmail := util.RandomEmail()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, newEmail, user.Email)
	require.Equal(t, newEmail, updatedUser.Email)

	require.Equal(t, user.FullName, updatedUser.FullName)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}
