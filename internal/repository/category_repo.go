package repository

import (
    "database/sql"
    "errors"

    "studyspot/internal/domain"

    "github.com/google/uuid"
)

type CategoryRepository struct {
    db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *domain.Category) error {
    query := `INSERT INTO categories (id, name) VALUES ($1, $2)`
    _, err := r.db.Exec(query, category.ID, category.Name)
    return err
}

func (r *CategoryRepository) FindAll() ([]domain.Category, error) {
    query := `SELECT id, name, created_at FROM categories ORDER BY name`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []domain.Category
    for rows.Next() {
        var category domain.Category
        err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func (r *CategoryRepository) FindByID(id uuid.UUID) (*domain.Category, error) {
    query := `SELECT id, name, created_at FROM categories WHERE id = $1`
    var category domain.Category
    err := r.db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &category, nil
}

func (r *CategoryRepository) FindByName(name string) (*domain.Category, error) {
    query := `SELECT id, name, created_at FROM categories WHERE name = $1`
    var category domain.Category
    err := r.db.QueryRow(query, name).Scan(&category.ID, &category.Name, &category.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &category, nil
}

func (r *CategoryRepository) Update(category *domain.Category) error {
    query := `UPDATE categories SET name = $1 WHERE id = $2`
    _, err := r.db.Exec(query, category.Name, category.ID)
    return err
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
    query := `DELETE FROM categories WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}