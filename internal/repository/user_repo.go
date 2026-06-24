package repository

import (
    "database/sql"
    "errors"

    "studyspot/internal/domain"

    "github.com/google/uuid"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
    query := `INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, user.ID, user.Email, user.Password, user.Role)
    return err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
    query := `SELECT id, email, password, role, created_at FROM users WHERE email = $1`
    var user domain.User
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
        &user.Role,
        &user.CreatedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
    query := `SELECT id, email, password, role, created_at FROM users WHERE id = $1`
    var user domain.User
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
        &user.Role,
        &user.CreatedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
    query := `UPDATE users SET email = $1, password = $2, role = $3 WHERE id = $4`
    _, err := r.db.Exec(query, user.Email, user.Password, user.Role, user.ID)
    return err
}