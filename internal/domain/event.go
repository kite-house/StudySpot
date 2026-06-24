package domain

import (
    "time"

    "github.com/google/uuid"
)

type Event struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    CategoryID  uuid.UUID `json:"category_id"`
    Date        time.Time `json:"date" binding:"required"`
    Location    string    `json:"location"`
    CreatedBy   uuid.UUID `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type EventWithCategory struct {
    Event
    CategoryName string `json:"category_name"`
}

type CreateEventRequest struct {
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    CategoryID  uuid.UUID `json:"category_id"`
    Date        time.Time `json:"date" binding:"required"`
    Location    string    `json:"location"`
}

type UpdateEventRequest struct {
    Title       string     `json:"title"`
    Description *string    `json:"description"`
    CategoryID  *uuid.UUID `json:"category_id"`
    Date        *time.Time `json:"date"`
    Location    *string    `json:"location"`
}