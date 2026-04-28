package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Role, sistemdeki kullanıcı rollerini tanımlar.
type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleSeller Role = "SELLER"
	RoleBuyer  Role = "BUYER"
)

// User, veritabanındaki users tablosunun Go karşılığıdır.
// PasswordHash alanı hiçbir zaman JSON olarak serialize edilmez.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // API yanıtına asla dahil edilmez
	Role         Role      `json:"role"`
	TenantID     string    `json:"tenant_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// DemoClaims, JWT token içinde taşınan payload yapısıdır.
type DemoClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}
