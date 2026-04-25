package test

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/slupx/smartest-backend/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type testResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(testResponse{
		Message: "test handler ok",
	})
}

func (h *Handler) CreateTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	var createdBy *uuid.UUID
	if user, ok := r.Context().Value("user").(*auth.User); ok {
		createdBy = &user.ID
	}

	test, err := h.service.CreateTest(r.Context(), req.Title, req.Questions, createdBy)
	if err != nil {
		http.Error(w, "failed to create test", http.StatusInternalServerError)
		return
	}

	writeJSON(w, test, http.StatusCreated)
}

func (h *Handler) UpdateTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing test id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid test id", http.StatusBadRequest)
		return
	}

	var req updateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	test, err := h.service.UpdateTest(r.Context(), id, req.Title, req.Questions)
	if err != nil {
		http.Error(w, "failed to update test", http.StatusInternalServerError)
		return
	}

	writeJSON(w, test, http.StatusOK)
}

func (h *Handler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing test id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid test id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTest(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to delete test", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing test id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid test id", http.StatusBadRequest)
		return
	}

	test, err := h.service.GetTest(r.Context(), id)
	if err != nil {
		http.Error(w, "test not found", http.StatusNotFound)
		return
	}

	writeJSON(w, test, http.StatusOK)
}

func (h *Handler) ListTests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("ListTests handler called")
	tests, err := h.service.ListTests(r.Context())
	if err != nil {
		http.Error(w, "failed to list tests", http.StatusInternalServerError)
		return
	}

	writeJSON(w, tests, http.StatusOK)
}

func (h *Handler) JoinTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	test, err := h.service.JoinTest(r.Context(), code)
	if err != nil {
		http.Error(w, "test not found", http.StatusNotFound)
		return
	}

	writeJSON(w, test, http.StatusOK)
}

func (h *Handler) SubmitTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req submitTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	testID, err := uuid.Parse(req.TestID)
	if err != nil {
		http.Error(w, "invalid test id", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("user").(*auth.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	submission, err := h.service.SubmitTest(r.Context(), testID, user.ID, req.Answers)
	if err != nil {
		http.Error(w, "failed to submit test", http.StatusInternalServerError)
		return
	}

	writeJSON(w, submission, http.StatusCreated)
}

func (h *Handler) GetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testIDStr := r.URL.Query().Get("test_id")
	if testIDStr == "" {
		http.Error(w, "test_id is required", http.StatusBadRequest)
		return
	}

	testID, err := uuid.Parse(testIDStr)
	if err != nil {
		http.Error(w, "invalid test id", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("user").(*auth.User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role == auth.RoleTeacher {
		results, err := h.service.GetResults(r.Context(), testID)
		if err != nil {
			http.Error(w, "failed to get results", http.StatusInternalServerError)
			return
		}
		writeJSON(w, results, http.StatusOK)
	} else {
		result, err := h.service.GetMyResult(r.Context(), testID, user.ID)
		if err != nil {
			http.Error(w, "result not found", http.StatusNotFound)
			return
		}
		writeJSON(w, result, http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
