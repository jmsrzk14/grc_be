package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type RegulationService struct {
	uc  *biz.RegulationUseCase
	log *log.Helper
}

func NewRegulationService(uc *biz.RegulationUseCase, logger log.Logger) *RegulationService {
	return &RegulationService{uc: uc, log: log.NewHelper(logger)}
}

// --- DTOs ---

type CreateRegulationRequest struct {
	Title          string `json:"title"`
	RegulationType string `json:"regulation_type"`
	IssuedDate     string `json:"issued_date"`
	Status         string `json:"status"`
	Category       string `json:"category"`
}

type RegulationResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	RegulationType string    `json:"regulation_type"`
	IssuedDate     time.Time `json:"issued_date"`
	Status         string    `json:"status"`
	Category       string    `json:"category"`
	AmountPass     int       `json:"amount_pass"`
	AmountFail     int       `json:"amount_fail"`
	AmountNA       int       `json:"amount_na"`
}

type UpdateRegulationRequest struct {
	Title          string `json:"title"`
	RegulationType string `json:"regulation_type"`
	IssuedDate     string `json:"issued_date"`
	Status         string `json:"status"`
	Category       string `json:"category"`
}

type CreateRegulationItemRequest struct {
	ItemCode        string   `json:"item_code"`
	ReferenceNumber string   `json:"reference_number"`
	Content         string   `json:"content"`
	PropertyIDs     []string `json:"property_ids"`
}

type RegulationItemResponse struct {
	ID              string   `json:"id"`
	RegulationID    string   `json:"regulation_id"`
	PropertyIDs     []string `json:"property_ids"`
	ItemCode        string   `json:"item_code"`
	ReferenceNumber string   `json:"reference_number"`
	Content         string   `json:"content"`
}

type AddMappingRequest struct {
	PropertyID string `json:"property_id"`
}

// --- Regulation Handlers ---

// CreateRegulation godoc
// @Tags RegulationsService
// @Accept json
// @Produce json
// @Param regulation body CreateRegulationRequest true "Regulation Data"
// @Success 201 {object} RegulationResponse
// @Router /api/v1/regulations [post]
func (s *RegulationService) CreateRegulation(w http.ResponseWriter, r *http.Request) {
	var req CreateRegulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	issuedDate, _ := time.Parse("2006-01-02", req.IssuedDate)
	reg := &biz.Regulation{
		Title:          req.Title,
		RegulationType: req.RegulationType,
		IssuedDate:     issuedDate,
		Status:         req.Status,
		Category:       req.Category,
	}
	if reg.Status == "" {
		reg.Status = "Active"
	}
	result, err := s.uc.CreateRegulation(r.Context(), reg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toRegulationResponse(result))
}

func (s *RegulationService) UpsertRegulation(w http.ResponseWriter, r *http.Request) {
	var req CreateRegulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	issuedDate, _ := time.Parse("2006-01-02", req.IssuedDate)
	reg := &biz.Regulation{
		Title:          req.Title,
		RegulationType: req.RegulationType,
		IssuedDate:     issuedDate,
		Status:         req.Status,
		Category:       req.Category,
	}
	if reg.Status == "" {
		reg.Status = "Active"
	}
	result, err := s.uc.UpsertRegulation(r.Context(), reg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toRegulationResponse(result))
}

// GetRegulation godoc
// @Tags RegulationsService
// @Produce json
// @Param id path string true "Regulation ID"
// @Param tenant_id query string false "Tenant ID for chart calculation"
// @Success 200 {object} RegulationResponse
// @Router /api/v1/regulations/{id} [get]
func (s *RegulationService) GetRegulation(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tenantIDStr := r.URL.Query().Get("tenant_id")
	var tenantID uuid.UUID
	if tenantIDStr != "" {
		tenantID, _ = uuid.Parse(tenantIDStr)
	}

	result, err := s.uc.GetRegulation(r.Context(), id, tenantID)
	if err != nil {
		if errors.Is(err, biz.ErrNotFound) {
			respondError(w, http.StatusNotFound, "regulation not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toRegulationResponse(result))
}

// ListRegulations godoc
// @Tags RegulationsService
// @Produce json
// @Param tenant_id query string false "Tenant ID for chart calculation"
// @Success 200 {array} RegulationResponse
// @Router /api/v1/regulations [get]
func (s *RegulationService) ListRegulations(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.URL.Query().Get("tenant_id")
	var tenantID uuid.UUID
	if tenantIDStr != "" {
		tenantID, _ = uuid.Parse(tenantIDStr)
	}

	results, err := s.uc.ListRegulations(r.Context(), tenantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]*RegulationResponse, 0, len(results))
	for _, reg := range results {
		resp = append(resp, toRegulationResponse(reg))
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateRegulation godoc
// @Tags RegulationsService
// @Accept json
// @Produce json
// @Param id path string true "Regulation ID"
// @Param regulation body UpdateRegulationRequest true "Updated Data"
// @Success 200 {object} RegulationResponse
// @Router /api/v1/regulations/{id} [put]
func (s *RegulationService) UpdateRegulation(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateRegulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	issuedDate, _ := time.Parse("2006-01-02", req.IssuedDate)
	reg := &biz.Regulation{
		ID:             id,
		Title:          req.Title,
		RegulationType: req.RegulationType,
		IssuedDate:     issuedDate,
		Status:         req.Status,
		Category:       req.Category,
	}
	result, err := s.uc.UpdateRegulation(r.Context(), reg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toRegulationResponse(result))
}

// DeleteRegulation godoc
// @Tags RegulationsService
// @Param id path string true "Regulation ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/regulations/{id} [delete]
func (s *RegulationService) DeleteRegulation(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.uc.DeleteRegulation(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "regulation deleted"})
}

// --- Item Handlers ---

// CreateItem godoc
// @Tags RegulationsService
// @Accept json
// @Produce json
// @Param id path string true "Regulation ID"
// @Param item body CreateRegulationItemRequest true "Item Data"
// @Success 201 {object} RegulationItemResponse
// @Router /api/v1/regulations/{id}/items [post]
func (s *RegulationService) CreateItem(w http.ResponseWriter, r *http.Request) {
	regID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid regulation id")
		return
	}
	var req CreateRegulationItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	propertyIDs := make([]uuid.UUID, 0, len(req.PropertyIDs))
	for _, idStr := range req.PropertyIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			propertyIDs = append(propertyIDs, id)
		}
	}
	item := &biz.RegulationItem{
		RegulationID:    regID,
		PropertyIDs:     propertyIDs,
		ItemCode:        req.ItemCode,
		ReferenceNumber: req.ReferenceNumber,
		Content:         req.Content,
	}
	result, err := s.uc.CreateRegulationItem(r.Context(), item)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toItemResponse(result))
}

func (s *RegulationService) UpsertItem(w http.ResponseWriter, r *http.Request) {
	regID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid regulation id")
		return
	}
	var req CreateRegulationItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	propertyIDs := make([]uuid.UUID, 0, len(req.PropertyIDs))
	for _, idStr := range req.PropertyIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			propertyIDs = append(propertyIDs, id)
		}
	}
	item := &biz.RegulationItem{
		RegulationID:    regID,
		PropertyIDs:     propertyIDs,
		ItemCode:        req.ItemCode,
		ReferenceNumber: req.ReferenceNumber,
		Content:         req.Content,
	}
	result, err := s.uc.UpsertRegulationItem(r.Context(), item)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toItemResponse(result))
}

// ListItems godoc
// @Tags RegulationsService
// @Produce json
// @Param id path string true "Regulation ID"
// @Success 200 {array} RegulationItemResponse
// @Router /api/v1/regulations/{id}/items [get]
func (s *RegulationService) ListItems(w http.ResponseWriter, r *http.Request) {
	regID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid regulation id")
		return
	}
	tenantIDStr := r.URL.Query().Get("tenant_id")
	var tenantID uuid.UUID
	if tenantIDStr != "" {
		tenantID, _ = uuid.Parse(tenantIDStr)
	}

	results, err := s.uc.ListRegulationItems(r.Context(), regID, tenantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]*RegulationItemResponse, 0, len(results))
	for _, item := range results {
		resp = append(resp, toItemResponse(item))
	}
	respondJSON(w, http.StatusOK, resp)
}

// GetItem godoc
// @Tags RegulationsService
// @Produce json
// @Param id path string true "Regulation ID"
// @Param item_id path string true "Item ID"
// @Success 200 {object} RegulationItemResponse
// @Router /api/v1/regulations/{id}/items/{item_id} [get]
func (s *RegulationService) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "item_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	result, err := s.uc.GetRegulationItem(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toItemResponse(result))
}

// UpdateItem godoc
// @Tags RegulationsService
// @Accept json
// @Produce json
// @Param id path string true "Regulation ID"
// @Param item_id path string true "Item ID"
// @Param item body CreateRegulationItemRequest true "Updated Data"
// @Success 200 {object} RegulationItemResponse
// @Router /api/v1/regulations/{id}/items/{item_id} [put]
func (s *RegulationService) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "item_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	var req CreateRegulationItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	propertyIDs := make([]uuid.UUID, 0, len(req.PropertyIDs))
	for _, idStr := range req.PropertyIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			propertyIDs = append(propertyIDs, id)
		}
	}
	item := &biz.RegulationItem{
		ID:              id,
		PropertyIDs:     propertyIDs,
		ItemCode:        req.ItemCode,
		ReferenceNumber: req.ReferenceNumber,
		Content:         req.Content,
	}
	result, err := s.uc.UpdateRegulationItem(r.Context(), item)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toItemResponse(result))
}

// DeleteItem godoc
// @Tags RegulationsService
// @Param id path string true "Regulation ID"
// @Param item_id path string true "Item ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/regulations/{id}/items/{item_id} [delete]
func (s *RegulationService) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "item_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	if err := s.uc.DeleteRegulationItem(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "item deleted"})
}

// --- Mapping Handlers ---

// AddMapping godoc
// @Tags RegulationsService
// @Accept json
// @Produce json
// @Param id path string true "Regulation ID"
// @Param mapping body AddMappingRequest true "Mapping Data"
// @Success 201 {object} biz.RegulationPropertyMapping
// @Router /api/v1/regulations/{id}/mappings [post]
func (s *RegulationService) AddMapping(w http.ResponseWriter, r *http.Request) {
	regID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid regulation id")
		return
	}
	var req AddMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	propID, err := uuid.Parse(req.PropertyID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid property id")
		return
	}
	mapping := &biz.RegulationPropertyMapping{
		RegulationID: regID,
		PropertyID:   propID,
	}
	result, err := s.uc.AddPropertyToRegulation(r.Context(), mapping)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

// ListMappings godoc
// @Tags RegulationsService
// @Produce json
// @Param id path string true "Regulation ID"
// @Success 200 {array} biz.RegulationPropertyMapping
// @Router /api/v1/regulations/{id}/mappings [get]
func (s *RegulationService) ListMappings(w http.ResponseWriter, r *http.Request) {
	regID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid regulation id")
		return
	}
	results, err := s.uc.ListRegulationMappings(r.Context(), regID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, results)
}

// DeleteMapping godoc
// @Tags RegulationsService
// @Param id path string true "Regulation ID"
// @Param mapping_id path string true "Mapping ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/regulations/{id}/mappings/{mapping_id} [delete]
func (s *RegulationService) DeleteMapping(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "mapping_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid mapping id")
		return
	}
	if err := s.uc.DeleteRegulationMapping(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "mapping deleted"})
}

func toRegulationResponse(r *biz.Regulation) *RegulationResponse {
	return &RegulationResponse{
		ID:             r.ID.String(),
		Title:          r.Title,
		RegulationType: r.RegulationType,
		IssuedDate:     r.IssuedDate,
		Status:         r.Status,
		Category:       r.Category,
		AmountPass:     r.AmountPass,
		AmountFail:     r.AmountFail,
		AmountNA:       r.AmountNA,
	}
}

func toItemResponse(r *biz.RegulationItem) *RegulationItemResponse {
	propertyIDs := make([]string, 0, len(r.PropertyIDs))
	for _, id := range r.PropertyIDs {
		propertyIDs = append(propertyIDs, id.String())
	}
	return &RegulationItemResponse{
		ID:              r.ID.String(),
		RegulationID:    r.RegulationID.String(),
		PropertyIDs:     propertyIDs,
		ItemCode:        r.ItemCode,
		ReferenceNumber: r.ReferenceNumber,
		Content:         r.Content,
	}
}
