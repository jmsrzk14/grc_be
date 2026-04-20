package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// TenantService handles HTTP requests untuk Tenant.
type TenantService struct {
	uc  *biz.TenantUseCase
	log *log.Helper
}

// NewTenantService membuat instance TenantService.
func NewTenantService(uc *biz.TenantUseCase, logger log.Logger) *TenantService {
	return &TenantService{uc: uc, log: log.NewHelper(logger)}
}

// --- Request/Response DTOs ---

type CreateTenantRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type TenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateTenantRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// --- Handlers ---

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant (BPR/BPRS)
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body CreateTenantRequest true "Tenant Data"
// @Success 201 {object} TenantResponse
// @Router /api/v1/tenants [post]
func (s *TenantService) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenant := &biz.Tenant{
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
	}
	if tenant.Status == "" {
		tenant.Status = "Active"
	}

	result, err := s.uc.CreateTenant(r.Context(), tenant)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toTenantResponse(result))
}

// GetTenant godoc
// @Summary Get tenant by ID
// @Description Get a specific tenant details
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} TenantResponse
// @Router /api/v1/tenants/{id} [get]
func (s *TenantService) GetTenant(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	result, err := s.uc.GetTenant(r.Context(), id)
	if err != nil {
		if errors.Is(err, biz.ErrNotFound) {
			respondError(w, http.StatusNotFound, "tenant not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTenantResponse(result))
}

// ListTenants godoc
// @Summary List all tenants
// @Description Get a list of all tenants
// @Tags tenants
// @Produce json
// @Success 200 {array} TenantResponse
// @Router /api/v1/tenants [get]
func (s *TenantService) ListTenants(w http.ResponseWriter, r *http.Request) {
	results, err := s.uc.ListTenants(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]*TenantResponse, 0, len(results))
	for _, t := range results {
		resp = append(resp, toTenantResponse(t))
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateTenant godoc
// @Summary Update tenant
// @Description Update an existing tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body UpdateTenantRequest true "Updated Data"
// @Success 200 {object} TenantResponse
// @Router /api/v1/tenants/{id} [put]
func (s *TenantService) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tenant := &biz.Tenant{ID: id, Name: req.Name, Type: req.Type, Status: req.Status}
	result, err := s.uc.UpdateTenant(r.Context(), tenant)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTenantResponse(result))
}

// DeleteTenant godoc
// @Summary Delete tenant
// @Description Delete a tenant by ID
// @Tags tenants
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{id} [delete]
func (s *TenantService) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.uc.DeleteTenant(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "tenant deleted"})
}

func toTenantResponse(t *biz.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Type:      t.Type,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
	}
}

// --- PropertyService ---

// PropertyService handles HTTP requests untuk Property.
type PropertyService struct {
	uc  *biz.PropertyUseCase
	log *log.Helper
}

// NewPropertyService membuat instance PropertyService.
func NewPropertyService(uc *biz.PropertyUseCase, logger log.Logger) *PropertyService {
	return &PropertyService{uc: uc, log: log.NewHelper(logger)}
}

type CreatePropertyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PropertyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateProperty godoc
// POST /api/v1/properties
func (s *PropertyService) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := s.uc.CreateProperty(r.Context(), &biz.Property{Name: req.Name, Description: req.Description})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toPropertyResponse(p))
}

// GetProperty godoc
// GET /api/v1/properties/{id}
func (s *PropertyService) GetProperty(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := s.uc.GetProperty(r.Context(), id)
	if err != nil {
		if errors.Is(err, biz.ErrNotFound) {
			respondError(w, http.StatusNotFound, "property not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toPropertyResponse(p))
}

// ListProperties godoc
// GET /api/v1/properties
func (s *PropertyService) ListProperties(w http.ResponseWriter, r *http.Request) {
	props, err := s.uc.ListProperties(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]*PropertyResponse, 0, len(props))
	for _, p := range props {
		resp = append(resp, toPropertyResponse(p))
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateProperty godoc
// PUT /api/v1/properties/{id}
func (s *PropertyService) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := s.uc.UpdateProperty(r.Context(), &biz.Property{ID: id, Name: req.Name, Description: req.Description})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toPropertyResponse(p))
}

// DeleteProperty godoc
// DELETE /api/v1/properties/{id}
func (s *PropertyService) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.uc.DeleteProperty(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "property deleted"})
}

func toPropertyResponse(p *biz.Property) *PropertyResponse {
	return &PropertyResponse{ID: p.ID.String(), Name: p.Name, Description: p.Description}
}

