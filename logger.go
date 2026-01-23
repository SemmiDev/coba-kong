package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger adalah custom middleware untuk logging setiap HTTP request
// Middleware ini akan mencatat informasi penting seperti method, path, status code, dan duration
// Menggunakan structured logging dengan logrus membuat log lebih mudah di-parse oleh tools monitoring
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Catat waktu mulai request
		startTime := time.Now()

		// Proses request (lanjutkan ke handler berikutnya)
		c.Next()

		// Hitung durasi pemrosesan request
		duration := time.Since(startTime)

		// Log dengan structured fields
		// Ini memudahkan filtering dan searching di log aggregation tools seperti ELK, Datadog, dll
		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP Request")
	}
}

// CORS middleware untuk menghandle Cross-Origin Resource Sharing
// Ini penting jika frontend Anda di domain berbeda dengan backend
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
