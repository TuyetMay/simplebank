package db

import (
	"time"
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"	
)

func createRandomUser(t *testing.T) User { // hàm nhận tham số đầu vào t kiểu dữ liệu là *tetsting.T và trả về kiểu dữ liệu Account
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t,err)

	arg := CreateUserParams{
		Username:   util.RandomOwner(),
		HashedPassword:  hashedPassword,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user

}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
	
}

func TestGetUser (t * testing.T) {
	user1 := createRandomUser(t)// dữ liệu inserted vào database
	user2, err := testQueries.GetUser(context.Background(), user1.Username) // gán lại các giá trị trả về từ account 1 vào account2 và err
// kiểm tra dữ liệu khi lấy từ DB là account2 có lỗi hay khác với account1 được thêm vào database không
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t,user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t,user1.FullName, user2.FullName)
	require.Equal(t,user1.Email, user2.Email)
	require.WithinDuration(t,user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second) // so sánh lại 2 giá trị thời gian, xem có chênh lệch nhau không
	require.WithinDuration(t,user1.CreatedAt, user2.CreatedAt, time.Second) // so sánh lại 2 giá trị thời gian, xem có chênh lệch nhau không
}