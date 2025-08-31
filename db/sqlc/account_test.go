package db

import (
	"time"
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
	"database/sql"
)

func createRandomAccount(t *testing.T) Account { // hàm nhận tham số đầu vào t kiểu dữ liệu là *tetsting.T và trả về kiểu dữ liệu Account
	arg := CreateAccountParams{
		Owner:   util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
	
}

func TestGetAccount(t * testing.T) {
	account1 := createRandomAccount(t)// dữ liệu inserted vào database
	account2, err := testQueries.GetAccount(context.Background(), account1.ID) // gán lại các giá trị trả về từ account 1 vào account2 và err
// kiểm tra dữ liệu khi lấy từ DB là account2 có lỗi hay khác với account1 được thêm vào database không
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t,account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t,account1.Balance, account2.Balance)
	require.Equal(t,account1.Currency, account2.Currency)
	require.WithinDuration(t,account1.CreatedAt, account2.CreatedAt, time.Second) // so sánh lại 2 giá trị thời gian, xem có chênh lệch nhau không
}

func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t,err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t,account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)

	err :=testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err,sql.ErrNoRows.Error()) // so sánh lỗi err với lỗi sql.ErrNoRows.Error() và lỗi đúng phải làm no error
	require.Empty(t, account2) // kiểm tra account2 rỗng
}

func TestListAccounts(t *testing.T){
	for i := 0; i< 10; i++{
		createRandomAccount(t)
	}

	arg := ListAccountsParams{	
		Limit: 5, // số luợng account trả về
		Offset: 5, // bỏ qua 5 account đầu tiên
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t,err)
	require.Len(t, accounts,5) // kiểm tra độ dài của accounts trả về có đúng 5 không

	for _, account := range accounts{
		require.NotEmpty(t,account) // kiểm tra account trả về có rỗng không
	}

}
