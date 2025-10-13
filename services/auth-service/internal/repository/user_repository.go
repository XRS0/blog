package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uint64    `bun:"id,pk,autoincrement"`
	Email     string    `bun:"email,notnull,unique"`
	Username  string    `bun:"username,notnull"`
	Password  string    `bun:"password,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, username, passwordHash string) (*User, error) {
	ctx := context.Background()
	user := &User{
		Email:     email,
		Username:  username,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	ctx := context.Background()
	user := new(User)

	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(id uint64) (*User, error) {
	ctx := context.Background()
	user := new(User)

	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	ctx := context.Background()
	exists, err := r.db.NewSelect().
		Model((*User)(nil)).
		Where("email = ?", email).
		Exists(ctx)

	return exists, err
}
