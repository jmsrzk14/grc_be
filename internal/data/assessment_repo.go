package data

import (
	"context"
	"errors"

	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- AssessmentSession Repository Implementation ---

type assessmentSessionRepo struct {
	data *Data
	log  *log.Helper
}

// NewAssessmentSessionRepo membuat instance repository AssessmentSession.
func NewAssessmentSessionRepo(data *Data, logger log.Logger) biz.AssessmentSessionRepo {
	return &assessmentSessionRepo{data: data, log: log.NewHelper(logger)}
}

func (r *assessmentSessionRepo) Create(ctx context.Context, session *biz.AssessmentSession) (*biz.AssessmentSession, error) {
	m := &AssessmentSessionModel{
		ID:         session.ID,
		TenantID:   session.TenantID,
		Title:      session.Title,
		PeriodYear: session.PeriodYear,
		Status:     session.Status,
		CreatedAt:  session.CreatedAt,
	}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	session.CreatedAt = m.CreatedAt
	return session, nil
}

func (r *assessmentSessionRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.AssessmentSession, error) {
	var m AssessmentSessionModel
	if result := r.data.db.WithContext(ctx).First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toSessionDomain(&m), nil
}

func (r *assessmentSessionRepo) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*biz.AssessmentSession, error) {
	var models []*AssessmentSessionModel
	if result := r.data.db.WithContext(ctx).Find(&models, "tenant_id = ?", tenantID); result.Error != nil {
		return nil, result.Error
	}
	sessions := make([]*biz.AssessmentSession, 0, len(models))
	for _, m := range models {
		sessions = append(sessions, toSessionDomain(m))
	}
	return sessions, nil
}

func (r *assessmentSessionRepo) FindAll(ctx context.Context) ([]*biz.AssessmentSession, error) {
	var models []*AssessmentSessionModel
	if result := r.data.db.WithContext(ctx).Find(&models); result.Error != nil {
		return nil, result.Error
	}
	sessions := make([]*biz.AssessmentSession, 0, len(models))
	for _, m := range models {
		sessions = append(sessions, toSessionDomain(m))
	}
	return sessions, nil
}

func (r *assessmentSessionRepo) Update(ctx context.Context, session *biz.AssessmentSession) (*biz.AssessmentSession, error) {
	m := &AssessmentSessionModel{
		ID:         session.ID,
		TenantID:   session.TenantID,
		Title:      session.Title,
		PeriodYear: session.PeriodYear,
		Status:     session.Status,
	}
	if result := r.data.db.WithContext(ctx).Save(m); result.Error != nil {
		return nil, result.Error
	}
	return toSessionDomain(m), nil
}

func (r *assessmentSessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&AssessmentSessionModel{}, "id = ?", id).Error
}

func toSessionDomain(m *AssessmentSessionModel) *biz.AssessmentSession {
	return &biz.AssessmentSession{
		ID:         m.ID,
		TenantID:   m.TenantID,
		Title:      m.Title,
		PeriodYear: m.PeriodYear,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt,
	}
}

// --- AssessmentResult Repository Implementation ---

type assessmentResultRepo struct {
	data *Data
	log  *log.Helper
}

// NewAssessmentResultRepo membuat instance repository AssessmentResult.
func NewAssessmentResultRepo(data *Data, logger log.Logger) biz.AssessmentResultRepo {
	return &assessmentResultRepo{data: data, log: log.NewHelper(logger)}
}

func (r *assessmentResultRepo) Create(ctx context.Context, result *biz.AssessmentResult) (*biz.AssessmentResult, error) {
	m := &AssessmentResultModel{
		ID:               result.ID,
		SessionID:        result.SessionID,
		RegulationItemID: result.RegulationItemID,
		ComplianceStatus: result.ComplianceStatus,
		EvidenceLink:     result.EvidenceLink,
		Remarks:          result.Remarks,
	}
	if res := r.data.db.WithContext(ctx).Create(m); res.Error != nil {
		return nil, res.Error
	}
	result.UpdatedAt = m.UpdatedAt
	return result, nil
}

func (r *assessmentResultRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.AssessmentResult, error) {
	var m AssessmentResultModel
	if result := r.data.db.WithContext(ctx).First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toResultDomain(&m), nil
}

func (r *assessmentResultRepo) FindBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*biz.AssessmentResult, error) {
	var models []*AssessmentResultModel
	if result := r.data.db.WithContext(ctx).Find(&models, "session_id = ?", sessionID); result.Error != nil {
		return nil, result.Error
	}
	results := make([]*biz.AssessmentResult, 0, len(models))
	for _, m := range models {
		results = append(results, toResultDomain(m))
	}
	return results, nil
}

// Upsert menyimpan atau memperbarui hasil berdasarkan (session_id, regulation_item_id).
func (r *assessmentResultRepo) Upsert(ctx context.Context, result *biz.AssessmentResult) (*biz.AssessmentResult, error) {
	m := &AssessmentResultModel{
		ID:               result.ID,
		SessionID:        result.SessionID,
		RegulationItemID: result.RegulationItemID,
		ComplianceStatus: result.ComplianceStatus,
		EvidenceLink:     result.EvidenceLink,
		Remarks:          result.Remarks,
	}
	res := r.data.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "session_id"}, {Name: "regulation_item_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"compliance_status", "evidence_link", "remarks", "updated_at"}),
		}).
		Create(m)
	if res.Error != nil {
		return nil, res.Error
	}
	result.UpdatedAt = m.UpdatedAt
	return result, nil
}

func (r *assessmentResultRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&AssessmentResultModel{}, "id = ?", id).Error
}

func toResultDomain(m *AssessmentResultModel) *biz.AssessmentResult {
	return &biz.AssessmentResult{
		ID:               m.ID,
		SessionID:        m.SessionID,
		RegulationItemID: m.RegulationItemID,
		ComplianceStatus: m.ComplianceStatus,
		EvidenceLink:     m.EvidenceLink,
		Remarks:          m.Remarks,
		UpdatedAt:        m.UpdatedAt,
	}
}

// --- RegulationAssessment Repository Implementation ---

type regulationAssessmentRepo struct {
	data *Data
	log  *log.Helper
}

// NewRegulationAssessmentRepo membuat instance repository RegulationAssessment.
func NewRegulationAssessmentRepo(data *Data, logger log.Logger) biz.RegulationAssessmentRepo {
	return &regulationAssessmentRepo{data: data, log: log.NewHelper(logger)}
}

func (r *regulationAssessmentRepo) Create(ctx context.Context, ra *biz.RegulationAssessment) (*biz.RegulationAssessment, error) {
	m := &RegulationAssessmentModel{
		ID:           ra.ID,
		RegulationID: ra.RegulationID,
		SessionID:    ra.SessionID,
		AmountPass:   ra.AmountPass,
		AmountFail:   ra.AmountFail,
		AmountNA:     ra.AmountNA,
	}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	return ra, nil
}

func (r *regulationAssessmentRepo) FindBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*biz.RegulationAssessment, error) {
	var models []*RegulationAssessmentModel
	if result := r.data.db.WithContext(ctx).Find(&models, "session_id = ?", sessionID); result.Error != nil {
		return nil, result.Error
	}
	ras := make([]*biz.RegulationAssessment, 0, len(models))
	for _, m := range models {
		ras = append(ras, toRegulationAssessmentDomain(m))
	}
	return ras, nil
}

func (r *regulationAssessmentRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.RegulationAssessment, error) {
	var m RegulationAssessmentModel
	if result := r.data.db.WithContext(ctx).First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toRegulationAssessmentDomain(&m), nil
}

func (r *regulationAssessmentRepo) Update(ctx context.Context, ra *biz.RegulationAssessment) (*biz.RegulationAssessment, error) {
	m := &RegulationAssessmentModel{
		ID:           ra.ID,
		RegulationID: ra.RegulationID,
		SessionID:    ra.SessionID,
		AmountPass:   ra.AmountPass,
		AmountFail:   ra.AmountFail,
		AmountNA:     ra.AmountNA,
	}
	if result := r.data.db.WithContext(ctx).Save(m); result.Error != nil {
		return nil, result.Error
	}
	return toRegulationAssessmentDomain(m), nil
}

// RecalculateForSession menghitung ulang summary pass/fail/na dari assessment_results.
func (r *regulationAssessmentRepo) RecalculateForSession(ctx context.Context, sessionID uuid.UUID, regulationID uuid.UUID) (*biz.RegulationAssessment, error) {
	type CountResult struct {
		Status string
		Count  int
	}
	var rows []CountResult

	err := r.data.db.WithContext(ctx).
		Table("assessment_results ar").
		Select("ar.compliance_status as status, COUNT(*) as count").
		Joins("JOIN regulation_items ri ON ri.id = ar.regulation_item_id").
		Where("ar.session_id = ? AND ri.regulation_id = ?", sessionID, regulationID).
		Group("ar.compliance_status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	ra := &biz.RegulationAssessment{
		RegulationID: regulationID,
		SessionID:    sessionID,
	}
	for _, row := range rows {
		switch row.Status {
		case "YES":
			ra.AmountPass = row.Count
		case "NO":
			ra.AmountFail = row.Count
		case "N/A":
			ra.AmountNA = row.Count
		}
	}

	// Upsert regulation assessment summary:
	m := &RegulationAssessmentModel{
		RegulationID: regulationID,
		SessionID:    sessionID,
		AmountPass:   ra.AmountPass,
		AmountFail:   ra.AmountFail,
		AmountNA:     ra.AmountNA,
	}
	res := r.data.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "regulation_id"}, {Name: "session_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"amount_pass", "amount_fail", "amount_na"}),
		}).
		Create(m)
	if res.Error != nil {
		return nil, res.Error
	}
	if m.ID != uuid.Nil {
		ra.ID = m.ID
	}
	return ra, nil
}

func toRegulationAssessmentDomain(m *RegulationAssessmentModel) *biz.RegulationAssessment {
	return &biz.RegulationAssessment{
		ID:           m.ID,
		RegulationID: m.RegulationID,
		SessionID:    m.SessionID,
		AmountPass:   m.AmountPass,
		AmountFail:   m.AmountFail,
		AmountNA:     m.AmountNA,
	}
}
