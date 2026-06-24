package handler

import (
    "studyspot/internal/domain"
    "studyspot/internal/service"
    "studyspot/pkg/response"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type CategoryHandler struct {
    categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
    return &CategoryHandler{
        categoryService: categoryService,
    }
}

// GetCategories @Summary Получить список категорий
// @Tags Categories
// @Produce json
// @Success 200 {object} response.Response{data=[]domain.Category}
// @Failure 500 {object} response.Response
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
    categories, err := h.categoryService.FindAll()
    if err != nil {
        response.InternalError(c, "Ошибка получения категорий: "+err.Error())
        return
    }

    response.Success(c, categories)
}

// CreateCategory @Summary Создать категорию (только админ)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateCategoryRequest true "Данные категории"
// @Success 201 {object} response.Response{data=domain.Category}
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
    var req domain.CreateCategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Неверный формат запроса: "+err.Error())
        return
    }

    category, err := h.categoryService.Create(&req)
    if err != nil {
        if err.Error() == "category already exists" {
            response.BadRequest(c, "Категория с таким именем уже существует")
            return
        }
        response.InternalError(c, "Ошибка создания категории: "+err.Error())
        return
    }

    response.Created(c, category)
}

// UpdateCategory @Summary Обновить категорию (только админ)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID категории"
// @Param request body domain.UpdateCategoryRequest true "Данные для обновления"
// @Success 200 {object} response.Response{data=domain.Category}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.BadRequest(c, "Неверный формат ID")
        return
    }

    var req domain.UpdateCategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Неверный формат запроса: "+err.Error())
        return
    }

    category, err := h.categoryService.Update(id, &req)
    if err != nil {
        if err.Error() == "category not found" {
            response.NotFound(c, "Категория не найдена")
            return
        }
        if err.Error() == "category with this name already exists" {
            response.BadRequest(c, "Категория с таким именем уже существует")
            return
        }
        response.InternalError(c, "Ошибка обновления категории: "+err.Error())
        return
    }

    response.Success(c, category)
}

// DeleteCategory @Summary Удалить категорию (только админ)
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID категории"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.BadRequest(c, "Неверный формат ID")
        return
    }

    if err := h.categoryService.Delete(id); err != nil {
        if err.Error() == "category not found" {
            response.NotFound(c, "Категория не найдена")
            return
        }
        response.InternalError(c, "Ошибка удаления категории: "+err.Error())
        return
    }

    response.SuccessWithMessage(c, "Категория успешно удалена", nil)
}