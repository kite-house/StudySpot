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

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	eventRepo    *repository.EventRepository
	cache        *cache.RedisCache
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, eventRepo *repository.EventRepository, cache *cache.RedisCache) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		eventRepo:    eventRepo,
		cache:        cache,
	}
}

func (s *CategoryService) Create(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	// Проверка существования категории с таким именем
	existing, err := s.categoryRepo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category already exists")
	}

	category := &domain.Category{
		ID:   uuid.New(),
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	// Очистка кэша
	ctx := context.Background()
	s.cache.Delete(ctx, "categories:all")

	return category, nil
}

func (s *CategoryService) FindAll() ([]domain.Category, error) {
	ctx := context.Background()
	cacheKey := "categories:all"

	// Попытка получить из кэша
	var categories []domain.Category
	err := s.cache.Get(ctx, cacheKey, &categories)
	if err == nil && categories != nil {
		return categories, nil
	}

	// Получение из БД
	categories, err = s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Сохранение в кэш
	s.cache.Set(ctx, cacheKey, categories, 30*time.Minute)

	return categories, nil
}

func (s *CategoryService) FindByID(id uuid.UUID) (*domain.Category, error) {
	cacheKey := fmt.Sprintf("category:%s", id.String())
	ctx := context.Background()

	// Попытка получить из кэша
	var category domain.Category
	err := s.cache.Get(ctx, cacheKey, &category)
	if err == nil && category.ID != uuid.Nil {
		return &category, nil
	}

	// Получение из БД
	categoryPtr, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if categoryPtr == nil {
		return nil, errors.New("category not found")
	}

	// Сохранение в кэш
	s.cache.Set(ctx, cacheKey, categoryPtr, 30*time.Minute)

	return categoryPtr, nil
}

func (s *CategoryService) Update(id uuid.UUID, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	// Проверка существования категории
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Проверка уникальности имени
	existing, err := s.categoryRepo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, errors.New("category with this name already exists")
	}

	category.Name = req.Name
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	// Очистка кэша
	ctx := context.Background()
	s.cache.Delete(ctx, "categories:all")
	s.cache.Delete(ctx, fmt.Sprintf("category:%s", id.String()))
	s.cache.DeletePattern(ctx, "events:*")

	return category, nil
}

func (s *CategoryService) Delete(id uuid.UUID) error {
	// Проверка существования категории
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	// Удаление категории
	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	// Обновление событий - убираем ссылку на категорию
	if err := s.eventRepo.DeleteByCategory(id); err != nil {
		return err
	}

	// Очистка кэша
	ctx := context.Background()
	s.cache.Delete(ctx, "categories:all")
	s.cache.Delete(ctx, fmt.Sprintf("category:%s", id.String()))
	s.cache.DeletePattern(ctx, "events:*")

	return nil
}
