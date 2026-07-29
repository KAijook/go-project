package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// Khởi tạo seed cho hàm rand để mỗi lần chạy test ra dữ liệu khác nhau
	rand.Seed(time.Now().UnixNano())
}

// RandomInt sinh ra một số nguyên ngẫu nhiên trong khoảng [min, max]
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString sinh ra một chuỗi ký tự ngẫu nhiên độ dài n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner sinh ra tên chủ tài khoản ngẫu nhiên (6 ký tự)
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney sinh ra số tiền ngẫu nhiên từ 100 đến 1000
func RandomMoney() int64 {
	return RandomInt(100, 1000)
}

// RandomCurrency chọn ngẫu nhiên loại tiền tệ
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "VND"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
