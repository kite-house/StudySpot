package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"studyspot/internal/domain"

	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *domain.Event) error {
	query := `INSERT INTO events (id, title, description, category_id, date, location, created_by) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, event.ID, event.Title, event.Description, event.CategoryID, event.Date, event.Location, event.CreatedBy)
	return err
}

func (r *EventRepository) FindAll(limit, offset int) ([]domain.EventWithCategory, error) {
	query := `SELECT e.*, c.name as category_name 
              FROM events e 
              LEFT JOIN categories c ON e.category_id = c.id 
              ORDER BY e.date ASC 
              LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

func (r *EventRepository) Search(query string, categoryID uuid.UUID, limit, offset int) ([]domain.EventWithCategory, error) {
	var conditions []string
	var args []interface{}
	argCounter := 1

	if query != "" {
		conditions = append(conditions, fmt.Sprintf("to_tsvector('russian', e.title) @@ plainto_tsquery('russian', $%d)", argCounter))
		args = append(args, query)
		argCounter++
	}

	if categoryID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("e.category_id = $%d", argCounter))
		args = append(args, categoryID)
		argCounter++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery := fmt.Sprintf(`SELECT e.*, c.name as category_name 
                             FROM events e 
                             LEFT JOIN categories c ON e.category_id = c.id 
                             %s 
                             ORDER BY e.date ASC 
                             LIMIT $%d OFFSET $%d`,
		whereClause, argCounter, argCounter+1)

	args = append(args, limit, offset)
	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

func (r *EventRepository) FindByID(id uuid.UUID) (*domain.EventWithCategory, error) {
	query := `SELECT e.*, c.name as category_name 
              FROM events e 
              LEFT JOIN categories c ON e.category_id = c.id 
              WHERE e.id = $1`
	var event domain.EventWithCategory
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.CategoryID,
		&event.Date,
		&event.Location,
		&event.CreatedBy,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.CategoryName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(event *domain.Event) error {
	query := `UPDATE events SET title = $1, description = $2, category_id = $3, date = $4, location = $5 
              WHERE id = $6`
	_, err := r.db.Exec(query, event.Title, event.Description, event.CategoryID, event.Date, event.Location, event.ID)
	return err
}

func (r *EventRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EventRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	return count, err
}

func (r *EventRepository) scanEvents(rows *sql.Rows) ([]domain.EventWithCategory, error) {
	var events []domain.EventWithCategory
	for rows.Next() {
		var event domain.EventWithCategory
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.CategoryID,
			&event.Date,
			&event.Location,
			&event.CreatedBy,
			&event.CreatedAt,
			&event.UpdatedAt,
			&event.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepository) DeleteByCategory(categoryID uuid.UUID) error {
	query := `UPDATE events SET category_id = NULL WHERE category_id = $1`
	_, err := r.db.Exec(query, categoryID)
	return err
}
