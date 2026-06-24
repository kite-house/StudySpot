package handler

import (
	"strconv"

	"studyspot/internal/domain"
	"studyspot/internal/service"
	"studyspot/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	events, total, err := h.eventService.FindAll(page, limit)
	if err != nil {
		response.InternalError(c, "Ошибка получения списка мероприятий: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"events": events,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

func (h *EventHandler) SearchEvents(c *gin.Context) {
	query := c.Query("q")
	categoryIDStr := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var categoryID uuid.UUID
	if categoryIDStr != "" {
		parsed, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = parsed
		}
	}

	events, total, err := h.eventService.Search(query, categoryID, page, limit)
	if err != nil {
		response.InternalError(c, "Ошибка поиска мероприятий: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"events": events,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Неверный формат ID")
		return
	}

	event, err := h.eventService.FindByID(id)
	if err != nil {
		if err.Error() == "event not found" {
			response.NotFound(c, "Мероприятие не найдено")
			return
		}
		response.InternalError(c, "Ошибка получения мероприятия: "+err.Error())
		return
	}

	response.Success(c, event)
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	createdBy := c.GetString("userID")
	userID, err := uuid.Parse(createdBy)
	if err != nil {
		response.InternalError(c, "Ошибка получения ID пользователя")
		return
	}

	event, err := h.eventService.Create(&req, userID)
	if err != nil {
		if err.Error() == "category not found" {
			response.BadRequest(c, "Категория не найдена")
			return
		}
		response.InternalError(c, "Ошибка создания мероприятия: "+err.Error())
		return
	}

	response.Created(c, event)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Неверный формат ID")
		return
	}

	var req domain.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	event, err := h.eventService.Update(id, &req)
	if err != nil {
		if err.Error() == "event not found" {
			response.NotFound(c, "Мероприятие не найдено")
			return
		}
		if err.Error() == "category not found" {
			response.BadRequest(c, "Категория не найдена")
			return
		}
		response.InternalError(c, "Ошибка обновления мероприятия: "+err.Error())
		return
	}

	response.Success(c, event)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Неверный формат ID")
		return
	}

	if err := h.eventService.Delete(id); err != nil {
		if err.Error() == "event not found" {
			response.NotFound(c, "Мероприятие не найдено")
			return
		}
		response.InternalError(c, "Ошибка удаления мероприятия: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Мероприятие успешно удалено", nil)
}
