package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/Yuruka00/go-user-subs/internal/domain"
	"github.com/Yuruka00/go-user-subs/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
	lg  *slog.Logger
	vld *validator.Validate
}

func NewSubscriptionHandler(s *service.SubscriptionService, l *slog.Logger) *SubscriptionHandler {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &SubscriptionHandler{
		svc: s,
		lg:  l,
		vld: v,
	}
}

// @Summary Создать подписку
// @Description Позволяет добавить новую подписку в базу данных для конкретного пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body SubscriptionCreateRequest true "Данные новой подписки"
// @Success 201 {object} domain.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [post]
func (ss *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request SubscriptionCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ss.lg.Debug("SubscriptionHandler.Create", "request decode error", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ss.lg.Debug("SubscriptionHandler.Create", "request data", request)

	err := ss.vld.Struct(&request)
	if err != nil {
		respondWithValidationError(w, err)
		return
	}

	item, err := request.ToDomain()

	if err != nil {
		ss.lg.Debug("SubscriptionHandler.Create", "request to domain error", err)

		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err = ss.svc.Create(r.Context(), item); err != nil {
		ss.lg.Error("SubscriptionHandler.Create", "error", err)

		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// @Summary Получить подписку
// @Description Позволяет получить информацию о подписке по её ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки в формате UUID"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (ss *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if len(idString) == 0 {
		respondWithError(w, http.StatusBadRequest, "id is missing")
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong id format")
		return
	}

	res, err := ss.svc.GetByID(r.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "subscription not found")
			return
		}

		ss.lg.Error("SubscriptionHandler.Get", "error", err)

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Обновить информацию о подписке
// @Description Позволяет выборочно обновить инфомрацию о подписке в базе данных
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки в формате UUID"
// @Param input body SubscriptionUpdateRequest true "Тело запроса с полями для обновления"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [patch]
func (ss *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if len(idString) == 0 {
		respondWithError(w, http.StatusBadRequest, "id is missing")
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong id format")
		return
	}

	var request SubscriptionUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ss.lg.Debug("SubscriptionHandler.Update", "request decode error", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ss.lg.Debug("SubscriptionHandler.Update", "request data", request)

	err = ss.vld.Struct(&request)
	if err != nil {
		respondWithValidationError(w, err)
		return
	}

	item, hasChanges, err := request.ToDomain()
	if err != nil {
		ss.lg.Debug("SubscriptionHandler.Update", "request to domain error", err)

		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	if !hasChanges {
		respondWithError(w, http.StatusBadRequest, "missing fields to update")
		return
	}

	item.ID = id
	if err := ss.svc.Update(r.Context(), item); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "subscription not found")
			return
		}

		ss.lg.Error("SubscriptionHandler.Update", "error", err)

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Удалить подписку
// @Description Позволяет удалить запись с подпиской из базы данных
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки в формате UUID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (ss *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if len(idString) == 0 {
		respondWithError(w, http.StatusBadRequest, "id is missing")
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong id format")
		return
	}

	err = ss.svc.Delete(r.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "subscription not found")
			return
		}

		ss.lg.Error("SubscriptionHandler.Delete", "error", err)

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Получить список подписок по фильтру
// @Description Позволяет получить список подписок по фильтру
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя в формате UUID"
// @Param service_name query string false "Название сервиса"
// @Param date_from query string true "Дата начала (YYYY-MM)"
// @Param date_to query string true "Дата конца (YYYY-MM)"
// @Success 200 {array} domain.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [get]
func (ss *SubscriptionHandler) GetList(w http.ResponseWriter, r *http.Request) {
	var filterRequest SubscriptionFilterRequest

	filterRequest.ServiceName = r.URL.Query().Get("service_name")

	userIDString := r.URL.Query().Get("user_id")
	if userIDString != "" {
		userID, err := uuid.Parse(userIDString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "wrong user_id format")
			return
		}
		filterRequest.UserID = userID
	}

	filterRequest.DateFrom = r.URL.Query().Get("date_from")
	filterRequest.DateTo = r.URL.Query().Get("date_to")

	ss.lg.Debug("SubscriptionHandler.GetList", "request data", filterRequest)

	filter, err := filterRequest.ToDomain()
	if err != nil {
		ss.lg.Debug("SubscriptionHandler.GetList", "request to domain error", err)

		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	res, err := ss.svc.GetList(r.Context(), filter)
	if err != nil {
		ss.lg.Error("SubscriptionHandler.GetList", "error", err)

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Получить итоговую стоимость подписок
// @Description Позволяет получить итоговую стоимость подписок по фильтру
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя в формате UUID"
// @Param service_name query string false "Название сервиса"
// @Param date_from query string true "Дата начала (YYYY-MM)"
// @Param date_to query string true "Дата конца (YYYY-MM)"
// @Success 200 {object} map[string]int "Пример: {"total": 1500}"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/total [get]
func (ss *SubscriptionHandler) GetTotalPrice(w http.ResponseWriter, r *http.Request) {
	var filterRequest SubscriptionFilterRequest

	filterRequest.ServiceName = r.URL.Query().Get("service_name")

	userIDString := r.URL.Query().Get("user_id")
	if userIDString != "" {
		userID, err := uuid.Parse(userIDString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "wrong user_id format")
			return
		}
		filterRequest.UserID = userID
	}

	filterRequest.DateFrom = r.URL.Query().Get("date_from")
	filterRequest.DateTo = r.URL.Query().Get("date_to")

	ss.lg.Debug("SubscriptionHandler.GetTotalPrice", "request data", filterRequest)

	filter, err := filterRequest.ToDomain()
	if err != nil {
		ss.lg.Debug("SubscriptionHandler.GetTotalPrice", "request to domain error", err)

		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	res, err := ss.svc.CalculateTotalPrice(r.Context(), filter)
	if err != nil {
		ss.lg.Error("SubscriptionHandler.GetTotalPrice", "error", err)

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Total int `json:"total"`
	}{Total: res})
}
