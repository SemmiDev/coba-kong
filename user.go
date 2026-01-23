package main

import "time"

// User merepresentasikan struktur data user dalam aplikasi
// Tags json digunakan untuk serialization/deserialization JSON
// Tags validate digunakan untuk validasi input (jika menggunakan validator library)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest adalah DTO (Data Transfer Object) untuk request membuat user baru
// Memisahkan request model dari domain model adalah best practice karena:
// 1. API contract terpisah dari internal data structure
// 2. Lebih mudah untuk validasi input
// 3. Mencegah user mengirim field yang tidak boleh di-set (seperti ID)
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=3,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest untuk request update user
// Semua field opsional karena user mungkin hanya ingin update sebagian data
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Email string `json:"email,omitempty" binding:"omitempty,email"`
}
