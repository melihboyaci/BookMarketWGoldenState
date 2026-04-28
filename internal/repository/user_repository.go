package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bookmarket/golden-state/internal/models"
)

// ErrUserNotFound, kullanıcı bulunamadığında döner.
// Handler katmanında 401 vs 404 ayrımı için kullanılır.
var ErrUserNotFound = errors.New("kullanıcı bulunamadı")

// FindUserByEmail, belirtilen e-posta ve tenant_id kombinasyonuna sahip
// kullanıcıyı veritabanından çeker.
// Kullanıcı yoksa ErrUserNotFound döner.
func FindUserByEmail(db *sql.DB, email, tenantID string) (*models.User, error) {
	const query = `
		SELECT id, username, email, password_hash, role, tenant_id, created_at
		FROM users
		WHERE email = $1 AND tenant_id = $2
		LIMIT 1
	`

	user := &models.User{}
	err := db.QueryRow(query, email, tenantID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.TenantID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("kullanıcı sorgulanırken hata: %w", err)
	}

	return user, nil
}

// CreateUser, yeni bir kullanıcıyı veritabanına ekler.
// user.TenantID çağıran tarafından ayarlanmış olmalıdır (her zaman 'demo_active').
// Oluşturulan kullanıcının UUID'si user.ID alanına yazılır.
func CreateUser(db *sql.DB, user *models.User) error {
	const query = `
		INSERT INTO users (username, email, password_hash, role, tenant_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.TenantID,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("kullanıcı oluşturulurken hata: %w", err)
	}

	return nil
}
