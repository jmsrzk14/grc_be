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
	Title          string    `json:"title"`
	RegulationType string    `json:"regulation_type"`
	IssuedDate     time.Time `json:"issued_date"`
	Status         string    `json:"status"`
}

type RegulationResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	RegulationType string    `json:"regulation_type"`
	IssuedDate     time.Time `json:"issued_date"`
	Status         string    `json:"status"`
}

type UpdateRegulationRequest struct {
	Title          string    `json:"title"`
	RegulationType string    `json:"regulation_type"`
	IssuedDate     time.Time `json:"issued_date"`
	Status         string    `json:"status"`
}

type CreateRegulationItemRequest struct {
	ReferenceNumber string `json:"reference_number"`
	Content         string `json:"content"`
}

type RegulationItemResponse struct {
	ID              string `json:"id"`
	RegulationID    string `json:"regulation_id"`
	ReferenceNumber string `json:"reference_number"`
	Content         string `json:"content"`
}

type AddMappingRequest struct {
	PropertyID string `json:"property_id"`
}

// --- Regulation Handlers ---

func (s *RegulationService) CreateRegulation(w http.ResponseWriter, r *http.Request) {
	var req CreateRegulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	reg := &biz.Regulation{
		Title:          req.Title,
		RegulationType: req.RegulationType,
		IssuedDate:     req.IssuedDate,
		Status:         req.Status,
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

func (s *RegulationService) GetRegulation(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	result, err := s.uc.GetRegulation(r.Context(), id)
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

func (s *RegulationService) ListRegulations(w http.ResponseWriter, r *http.Request) {
	results, err := s.uc.ListRegulations(r.Context())
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
	reg := &biz.Regulation{
		ID:             id,
		Title:          req.Title,
		RegulationType: req.RegulationType,
		IssuedDate:     req.IssuedDate,
		Status:         req.Status,
	}
	result, err := s.uc.UpdateRegulation(r.Context(), reg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toRegulationResponse(result))
}

// DeleteRegulation godoc
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
	item := &biz.RegulationItem{
		RegulationID:    regID,
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

func (s *RegulationService) ListItems(w http.ResponseWriter, r *http.Request) {
	regID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid regulation id")
		return
	}
	results, err := s.uc.ListRegulationItems(r.Context(), regID)
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
	item := &biz.RegulationItem{
		ID:              id,
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
	}
}

func toItemResponse(r *biz.RegulationItem) *RegulationItemResponse {
	return &RegulationItemResponse{
		ID:              r.ID.String(),
		RegulationID:    r.RegulationID.String(),
		ReferenceNumber: r.ReferenceNumber,
		Content:         r.Content,
	}
}
