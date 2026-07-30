package main

import (
	"context"
	"log"

	"my-gin-app/api"
	db "my-gin-app/db/sqlc"
	"my-gin-app/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Tải cấu hình ứng dụng từ file .env hoặc app.env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v\n", err)
	}

	// 2. Kết nối Database PostgreSQL với pgxpool
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Không thể kết nối Database: %v\n", err)
	}
	defer connPool.Close()

	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database không phản hồi: %v\n", err)
	}

	// 3. Khởi tạo Store và Server
	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("Không thể khởi tạo server: %v\n", err)
	}

	// 4. Lắng nghe và phục vụ Request
	log.Printf("Server đang chạy tại http://%s\n", config.ServerAddress)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("Không thể khởi chạy server: %v\n", err)
	}
}
