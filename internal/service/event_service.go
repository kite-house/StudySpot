package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"studyspot/internal/cache"
	"studyspot/internal/domain"
	"studyspot/internal/repository"

	"github.com/google/uuid"
)

type EventService struct {
	eventRepo    *repository.EventRepository
	categoryRepo *repository.CategoryRepository
	cache        *cache.RedisCache
}

func NewEventService(eventRepo *repository.EventRepository, categoryRepo *repository.CategoryRepository, cache *cache.RedisCache) *EventService {
	return &EventService{
		eventRepo:    eventRepo,
		categoryRepo: categoryRepo,
		cache:        cache,
	}
}

func (s *EventService) Create(req *domain.CreateEventRequest, createdBy uuid.UUID) (*domain.Event, error) {
	if req.CategoryID != uuid.Nil {
		category, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("category not found")
		}
	}

	event := &domain.Event{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Date:        req.Date,
		Location:    req.Location,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	ctx := context.Background()
	s.cache.DeletePattern(ctx, "events:*")

	return event, nil
}

func (s *EventService) FindAll(page, limit int) ([]domain.EventWithCategory, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("events:all:%d:%d", page, limit)

	ctx := context.Background()
	var events []domain.EventWithCategory
	err := s.cache.Get(ctx, cacheKey, &events)
	if err == nil && events != nil {
		total, _ := s.eventRepo.Count()
		return events, total, nil
	}

	events, err = s.eventRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.eventRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	s.cache.Set(ctx, cacheKey, events, 5*time.Minute)

	return events, total, nil
}

func (s *EventService) Search(query string, categoryID uuid.UUID, page, limit int) ([]domain.EventWithCategory, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("events:search:%s:%s:%d:%d", query, categoryID.String(), page, limit)

	ctx := context.Background()
	var events []domain.EventWithCategory
	err := s.cache.Get(ctx, cacheKey, &events)
	if err == nil && events != nil {
		return events, len(events), nil
	}

	events, err = s.eventRepo.Search(query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	s.cache.Set(ctx, cacheKey, events, 2*time.Minute)

	return events, len(events), nil
}

func (s *EventService) FindByID(id uuid.UUID) (*domain.EventWithCategory, error) {
	cacheKey := fmt.Sprintf("event:%s", id.String())

	ctx := context.Background()
	var event domain.EventWithCategory
	err := s.cache.Get(ctx, cacheKey, &event)
	if err == nil && event.ID != uuid.Nil {
		return &event, nil
	}

	eventPtr, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if eventPtr == nil {
		return nil, errors.New("event not found")
	}

	s.cache.Set(ctx, cacheKey, eventPtr, 10*time.Minute)

	return eventPtr, nil
}

func (s *EventService) Update(id uuid.UUID, req *domain.UpdateEventRequest) (*domain.Event, error) {
	currentEvent, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if currentEvent == nil {
		return nil, errors.New("event not found")
	}

	event := &domain.Event{
		ID:          id,
		Title:       currentEvent.Title,
		Description: currentEvent.Description,
		CategoryID:  currentEvent.CategoryID,
		Date:        currentEvent.Date,
		Location:    currentEvent.Location,
		CreatedBy:   currentEvent.CreatedBy,
		CreatedAt:   currentEvent.CreatedAt,
	}

	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.CategoryID != nil {
		if *req.CategoryID != uuid.Nil {
			category, err := s.categoryRepo.FindByID(*req.CategoryID)
			if err != nil {
				return nil, err
			}
			if category == nil {
				return nil, errors.New("category not found")
			}
		}
		event.CategoryID = *req.CategoryID
	}
	if req.Date != nil {
		event.Date = *req.Date
	}
	if req.Location != nil {
		event.Location = *req.Location
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	ctx := context.Background()
	s.cache.Delete(ctx, "event:"+id.String())
	s.cache.DeletePattern(ctx, "events:*")

	return event, nil
}

func (s *EventService) Delete(id uuid.UUID) error {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	if err := s.eventRepo.Delete(id); err != nil {
		return err
	}

	ctx := context.Background()
	s.cache.Delete(ctx, "event:"+id.String())
	s.cache.DeletePattern(ctx, "events:*")

	return nil
}
