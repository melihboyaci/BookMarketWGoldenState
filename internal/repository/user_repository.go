package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bookmarket/golden-state/internal/config"
	"github.com/bookmarket/golden-state/internal/models"
)

// ErrUserNotFound, kullanıcı bulunamadığında döner.
var ErrUserNotFound = errors.New("kullanıcı bulunamadı")

// UserRepository, kullanıcı veri işlemlerini soyutlar.
type UserRepository interface {
	FindByEmail(email, tenantID string) (*models.User, error)
	Create(user *models.User) error
	GetAll() ([]models.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// FindByEmail, belirtilen e-posta ve tenant_id kombinasyonuna sahip
// kullanıcıyı veritabanından çeker.
func (r *PostgresUserRepository) FindByEmail(email, tenantID string) (*models.User, error) {
	const query = `
		SELECT id, username, email, password_hash, role, tenant_id, created_at
		FROM users
		WHERE email = $1 AND tenant_id = $2
		LIMIT 1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email, tenantID).Scan(
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

// Create, yeni bir kullanıcıyı veritabanına ekler.
func (r *PostgresUserRepository) Create(user *models.User) error {
	const query = `
		INSERT INTO users (username, email, password_hash, role, tenant_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query,
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

// GetAll, demo_active tenant'ındaki tüm kullanıcıları döner.
func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	const query = `
		SELECT id, username, email, role, tenant_id, created_at
		FROM users
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, config.TenantActive)
	if err != nil {
		return nil, fmt.Errorf("kullanıcılar sorgulanırken hata: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.TenantID, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("kullanıcı satırı okunurken hata: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

