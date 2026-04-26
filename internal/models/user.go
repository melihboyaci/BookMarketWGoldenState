package models

import "github.com/golang-jwt/jwt/v5"

// Role, sistemdeki kullanıcı rollerini tanımlar.
type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleSeller Role = "SELLER"
	RoleBuyer  Role = "BUYER"
)

// DemoClaims, JWT token içinde taşınan payload yapısıdır.
type DemoClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}
