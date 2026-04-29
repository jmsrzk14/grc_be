package service

import (
	"encoding/json"
	"net/http"
	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RiskService struct {
	uc  *biz.RiskUseCase
	log *log.Helper
}

func NewRiskService(uc *biz.RiskUseCase, logger log.Logger) *RiskService {
	return &RiskService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// --- Risk Category Handlers ---

// CreateCategory handles creating a new risk category.
// @Description Create a new risk category
// @Tags RiskService
// @Accept json
// @Produce json
// @Param category body biz.RiskCategory true "Risk Category"
// @Success 201 {object} biz.RiskCategory
// @Router /api/v1/risk-categories [post]
func (s *RiskService) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category biz.RiskCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		s.log.Errorf("failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := s.uc.CreateCategory(r.Context(), &category)
	if err != nil {
		s.log.Errorf("failed to create category: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GetCategory handles getting a risk category by ID.
// @Description Get a risk category by ID
// @Tags RiskService
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} biz.RiskCategory
// @Router /api/v1/risk-categories/{id} [get]
func (s *RiskService) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	res, err := s.uc.GetCategory(r.Context(), id)
	if err != nil {
		s.log.Errorf("failed to get category: %v", err)
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// ListCategories handles listing all risk categories.
// @Description List all risk categories
// @Tags RiskService
// @Produce json
// @Success 200 {array} biz.RiskCategory
// @Router /api/v1/risk-categories [get]
func (s *RiskService) ListCategories(w http.ResponseWriter, r *http.Request) {
	res, err := s.uc.ListCategories(r.Context())
	if err != nil {
		s.log.Errorf("failed to list categories: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateCategory handles updating a risk category.
// @Description Update a risk category
// @Tags RiskService
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body biz.RiskCategory true "Risk Category"
// @Success 200 {object} biz.RiskCategory
// @Router /api/v1/risk-categories/{id} [put]
func (s *RiskService) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var category biz.RiskCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	category.ID = id

	res, err := s.uc.UpdateCategory(r.Context(), &category)
	if err != nil {
		s.log.Errorf("failed to update category: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DeleteCategory handles deleting a risk category.
// @Description Delete a risk category
// @Tags RiskService
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Router /api/v1/risk-categories/{id} [delete]
func (s *RiskService) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	if err := s.uc.DeleteCategory(r.Context(), id); err != nil {
		s.log.Errorf("failed to delete category: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Risk Handlers ---

// CreateRisk handles creating a new risk item.
// @Description Create a new risk item
// @Tags RiskService
// @Accept json
// @Produce json
// @Param risk body biz.Risk true "Risk Item"
// @Success 201 {object} biz.Risk
// @Router /api/v1/risks [post]
func (s *RiskService) CreateRisk(w http.ResponseWriter, r *http.Request) {
	var risk biz.Risk
	if err := json.NewDecoder(r.Body).Decode(&risk); err != nil {
		s.log.Errorf("failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := s.uc.CreateRisk(r.Context(), &risk)
	if err != nil {
		s.log.Errorf("failed to create risk: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GetRisk handles getting a risk item by ID.
// @Description Get a risk item by ID
// @Tags RiskService
// @Produce json
// @Param id path string true "Risk ID"
// @Success 200 {object} biz.Risk
// @Router /api/v1/risks/{id} [get]
func (s *RiskService) GetRisk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid risk id", http.StatusBadRequest)
		return
	}

	res, err := s.uc.GetRisk(r.Context(), id)
	if err != nil {
		s.log.Errorf("failed to get risk: %v", err)
		http.Error(w, "risk not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// ListRisks handles listing all risks for a tenant.
// @Description List all risks for a tenant
// @Tags RiskService
// @Produce json
// @Param tenant_id query string true "Tenant ID"
// @Success 200 {array} biz.Risk
// @Router /api/v1/risks [get]
func (s *RiskService) ListRisks(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	if tenantIDStr == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Error(w, "invalid tenant_id", http.StatusBadRequest)
		return
	}

	res, err := s.uc.ListRisks(r.Context(), tenantID)
	if err != nil {
		s.log.Errorf("failed to list risks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateRisk handles updating a risk item.
// @Description Update a risk item
// @Tags RiskService
// @Accept json
// @Produce json
// @Param id path string true "Risk ID"
// @Param risk body biz.Risk true "Risk Item"
// @Success 200 {object} biz.Risk
// @Router /api/v1/risks/{id} [put]
func (s *RiskService) UpdateRisk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid risk id", http.StatusBadRequest)
		return
	}

	var risk biz.Risk
	if err := json.NewDecoder(r.Body).Decode(&risk); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	risk.ID = id

	res, err := s.uc.UpdateRisk(r.Context(), &risk)
	if err != nil {
		s.log.Errorf("failed to update risk: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DeleteRisk handles deleting a risk item.
// @Description Delete a risk item
// @Tags RiskService
// @Param id path string true "Risk ID"
// @Success 204 "No Content"
// @Router /api/v1/risks/{id} [delete]
func (s *RiskService) DeleteRisk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid risk id", http.StatusBadRequest)
		return
	}

	if err := s.uc.DeleteRisk(r.Context(), id); err != nil {
		s.log.Errorf("failed to delete risk: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
