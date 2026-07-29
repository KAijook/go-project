// db/main_test.go
package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

const (
	dbSource = "postgres://postgres:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()

	// Khởi tạo connection pool cho DB test
	testDB, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("Không thể kết nối DB test:", err)
	}

	testQueries = New(testDB)

	// Chạy tất cả bài test và thoát
	os.Exit(m.Run())
}
