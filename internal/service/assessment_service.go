package service

import (
	"encoding/json"
	// "errors"
	"net/http"
	"time"

	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type AssessmentService struct {
	uc  *biz.AssessmentUseCase
	log *log.Helper
}

func NewAssessmentService(uc *biz.AssessmentUseCase, logger log.Logger) *AssessmentService {
	return &AssessmentService{uc: uc, log: log.NewHelper(logger)}
}

// --- DTOs ---

type CreateSessionRequest struct {
	TenantID   string `json:"tenant_id"`
	Title      string `json:"title"`
	PeriodYear int    `json:"period_year"`
}

type SessionResponse struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	Title      string    `json:"title"`
	PeriodYear int       `json:"period_year"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type SubmitResultRequest struct {
	RegulationItemID string `json:"regulation_item_id"`
	ComplianceStatus string `json:"compliance_status"`
	EvidenceLink     string `json:"evidence_link"`
	Remarks          string `json:"remarks"`
}

type ResultResponse struct {
	ID               string    `json:"id"`
	SessionID        string    `json:"session_id"`
	RegulationItemID string    `json:"regulation_item_id"`
	ComplianceStatus string    `json:"compliance_status"`
	EvidenceLink     string    `json:"evidence_link"`
	Remarks          string    `json:"remarks"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// --- Handlers ---

func (s *AssessmentService) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tenant id")
		return
	}
	session := &biz.AssessmentSession{
		TenantID:   tenantID,
		Title:      req.Title,
		PeriodYear: req.PeriodYear,
	}
	result, err := s.uc.CreateSession(r.Context(), session)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toSessionResponse(result))
}

func (s *AssessmentService) GetSession(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	result, err := s.uc.GetSession(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toSessionResponse(result))
}

func (s *AssessmentService) UpdateSession(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tenantID, _ := uuid.Parse(req.TenantID)
	session := &biz.AssessmentSession{
		ID:         id,
		TenantID:   tenantID,
		Title:      req.Title,
		PeriodYear: req.PeriodYear,
	}
	result, err := s.uc.UpdateSession(r.Context(), session)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toSessionResponse(result))
}

func (s *AssessmentService) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.uc.DeleteSession(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "session deleted"})
}

func (s *AssessmentService) ListSessions(w http.ResponseWriter, r *http.Request) {
	var tenantID *uuid.UUID
	tIDStr := r.URL.Query().Get("tenant_id")
	if tIDStr != "" {
		id, err := uuid.Parse(tIDStr)
		if err == nil {
			tenantID = &id
		}
	}
	results, err := s.uc.ListSessions(r.Context(), tenantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]*SessionResponse, 0, len(results))
	for _, s := range results {
		resp = append(resp, toSessionResponse(s))
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *AssessmentService) SubmitResult(w http.ResponseWriter, r *http.Request) {
	sessionID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	var req SubmitResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	itemID, err := uuid.Parse(req.RegulationItemID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	result := &biz.AssessmentResult{
		SessionID:        sessionID,
		RegulationItemID: itemID,
		ComplianceStatus: req.ComplianceStatus,
		EvidenceLink:     req.EvidenceLink,
		Remarks:          req.Remarks,
	}
	saved, err := s.uc.SubmitResult(r.Context(), result)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toResultResponse(saved))
}

func (s *AssessmentService) DeleteResult(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDFromRequest(r, "result_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid result id")
		return
	}
	if err := s.uc.DeleteResult(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "result deleted"})
}

func (s *AssessmentService) ListResults(w http.ResponseWriter, r *http.Request) {
	sessionID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	results, err := s.uc.ListResultsBySession(r.Context(), sessionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]*ResultResponse, 0, len(results))
	for _, res := range results {
		resp = append(resp, toResultResponse(res))
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *AssessmentService) GetSummaries(w http.ResponseWriter, r *http.Request) {
	sessionID, err := parseUUIDFromRequest(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	results, err := s.uc.GetRegulationAssessments(r.Context(), sessionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, results)
}

func toSessionResponse(s *biz.AssessmentSession) *SessionResponse {
	return &SessionResponse{
		ID:         s.ID.String(),
		TenantID:   s.TenantID.String(),
		Title:      s.Title,
		PeriodYear: s.PeriodYear,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt,
	}
}

func toResultResponse(r *biz.AssessmentResult) *ResultResponse {
	return &ResultResponse{
		ID:               r.ID.String(),
		SessionID:        r.SessionID.String(),
		RegulationItemID: r.RegulationItemID.String(),
		ComplianceStatus: r.ComplianceStatus,
		EvidenceLink:     r.EvidenceLink,
		Remarks:          r.Remarks,
		UpdatedAt:        r.UpdatedAt,
	}
}
