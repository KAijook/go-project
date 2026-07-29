package main

import (
	"context"
	"log"
	"net/http"

	db "my-gin-app/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	store  *db.Queries
	router *gin.Engine
}

func main() {

	dbSource := "postgres://postgres:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"

	// 1. Tạo Connection Pool bằng pgx
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Không thể kết nối Database: %v\n", err)
	}
	defer connPool.Close()

	// Kiểm tra kết nối
	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database không phản hồi: %v\n", err)
	}

	// 2. Khởi tạo store từ sqlc
	store := db.New(connPool)

	// 3. Khởi tạo Gin Router
	router := gin.Default()

	server := &Server{
		store:  store,
		router: router,
	}

	// 4. Khai báo các Routes API
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)

	// 5. Chạy Server ở port 8080
	log.Println("Server đang chạy tại http://localhost:8080")
	router.Run(":8080")
}

// ---------------- HANDLERS ----------------

// Struct nhận JSON gửi lên từ Client
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance" binding:"required,min=0"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  req.Balance, // Sử dụng số dư từ request
	}

	account, err := server.store.CreateAccount(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := server.store.GetAccount(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài khoản"})
		return
	}

	c.JSON(http.StatusOK, account)
}

type deleteAccount struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(c *gin.Context) {}
