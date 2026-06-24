package domain

import (
    "time"

    "github.com/google/uuid"
)

type Category struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name" binding:"required"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateCategoryRequest struct {
    Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
    Name string `json:"name" binding:"required"`
}