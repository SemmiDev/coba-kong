package main

import (
	"github.com/gin-gonic/gin"
)

// APIResponse adalah struktur standar untuk semua response API
// Menggunakan response format yang konsisten membuat API lebih predictable dan mudah dikonsumsi
// Format ini mengikuti best practice dengan memisahkan success/error state dan menyertakan metadata
type APIResponse struct {
	Success bool        `json:"success"`         // Indikator apakah request berhasil atau tidak
	Message string      `json:"message"`         // Pesan untuk user atau developer
	Data    interface{} `json:"data,omitempty"`  // Actual payload, omitempty jika data kosong
	Error   interface{} `json:"error,omitempty"` // Error details jika ada
}

// Success mengirim response untuk request yang berhasil
// Function ini memastikan semua success response punya format yang sama
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error mengirim response untuk request yang gagal
// Memisahkan error response memudahkan client untuk handle error dengan konsisten
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
