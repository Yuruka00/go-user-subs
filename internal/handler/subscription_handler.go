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

func (ss *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request SubscriptionCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := ss.vld.Struct(&request)
	if err != nil {
		respondWithValidationError(w, err)
		return
	}

	item, err := request.ToDomain()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err = ss.svc.Create(r.Context(), item); err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

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
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(res)
}

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ss.vld.Struct(&request)
	if err != nil {
		respondWithValidationError(w, err)
		return
	}

	item, hasChanges, err := request.ToDomain()
	if err != nil {
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
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
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
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (ss *SubscriptionHandler) GetList(w http.ResponseWriter, r *http.Request) {
	var filterRequest SubscriptionFilterRequest

	filterRequest.ServiceName = r.URL.Query().Get("service_name")

	userIDString := r.URL.Query().Get("user_id")
	if userIDString != "" {
		userID, err := uuid.Parse(r.URL.Query().Get("user_id"))

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "wrong user_id format")
			return
		}
		filterRequest.UserID = userID
	}

	filterRequest.DateFrom = r.URL.Query().Get("date_from")
	filterRequest.DateTo = r.URL.Query().Get("date_to")

	filter, err := filterRequest.ToDomain()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	res, err := ss.svc.GetList(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
func (ss *SubscriptionHandler) GetTotalPrice(w http.ResponseWriter, r *http.Request) {
	var filterRequest SubscriptionFilterRequest

	filterRequest.ServiceName = r.URL.Query().Get("service_name")

	userIDString := r.URL.Query().Get("user_id")
	if userIDString != "" {
		userID, err := uuid.Parse(r.URL.Query().Get("user_id"))

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "wrong user_id format")
			return
		}
		filterRequest.UserID = userID
	}

	filterRequest.DateFrom = r.URL.Query().Get("date_from")
	filterRequest.DateTo = r.URL.Query().Get("date_to")

	filter, err := filterRequest.ToDomain()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	res, err := ss.svc.CalculateTotalPrice(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(struct {
		Total int `json:"total"`
	}{Total: res})
}
