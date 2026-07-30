// db/main_test.go
package db

import (
	"context"
	"log"
	"my-gin-app/util"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Không thể tải cấu hình:", err)
	}

	// Khởi tạo connection pool cho DB test
	testDB, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal("Không thể kết nối DB test:", err)
	}

	testQueries = New(testDB)

	// Chạy tất cả bài test và thoát
	os.Exit(m.Run())
}
