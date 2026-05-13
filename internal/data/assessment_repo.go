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
	// Order by created_at DESC agar session terbaru selalu di index 0 (dipakai frontend)
	if result := r.data.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&models); result.Error != nil {
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

func (r *assessmentSessionRepo) FindByTenantAndYear(ctx context.Context, tenantID uuid.UUID, year int) (*biz.AssessmentSession, error) {
	var m AssessmentSessionModel
	if result := r.data.db.WithContext(ctx).
		Where("tenant_id = ? AND period_year = ?", tenantID, year).
		First(&m); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toSessionDomain(&m), nil
}

func (r *assessmentSessionRepo) CheckAndComplete(ctx context.Context, sessionID uuid.UUID) error {
	// Logika ini memastikan status 'Completed' hanya diberikan jika SEMUA item
	// dari SEMUA regulasi yang aktif dalam sesi ini telah dijawab.

	// 1. Ambil data session untuk mendapatkan tenant_id dan period_year (untuk audit/logika jika diperlukan)
	var session AssessmentSessionModel
	if err := r.data.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error; err != nil {
		return err
	}

	// 2. Dapatkan semua regulation_id yang seharusnya diperiksa:
	// - Semua regulasi Eksternal (karena berlaku umum)
	// - Regulasi Internal yang sudah aktif di sesi ini (tercatat di regulation_assesments)
	var activeRegIDs []uuid.UUID
	err := r.data.db.WithContext(ctx).
		Table("regulations").
		Where("category != 'Internal'").
		Or("id IN (SELECT regulation_id FROM regulation_assesments WHERE session_id = ? AND is_active = ?)", sessionID, true).
		Pluck("id", &activeRegIDs).Error
	if err != nil {
		return err
	}

	if len(activeRegIDs) == 0 {
		return nil
	}

	// 3. Cari semua id item dari tabel regulation_items berdasarkan regulation_id tadi
	var totalItemIDs []uuid.UUID
	err = r.data.db.WithContext(ctx).
		Table("regulation_items").
		Where("regulation_id IN ?", activeRegIDs).
		Pluck("id", &totalItemIDs).Error
	if err != nil {
		return err
	}

	if len(totalItemIDs) == 0 {
		return nil
	}

	// 4. Dapatkan semua regulation_item_id yang sudah ada hasilnya di assessment_results untuk session ini
	var answeredItemIDs []uuid.UUID
	err = r.data.db.WithContext(ctx).
		Table("assessment_results").
		Where("session_id = ?", sessionID).
		Pluck("regulation_item_id", &answeredItemIDs).Error
	if err != nil {
		return err
	}

	// 5. Hitung sisa item yang BELUM dikerjakan
	answeredMap := make(map[uuid.UUID]bool)
	for _, id := range answeredItemIDs {
		answeredMap[id] = true
	}

	var remainingUnanswered int
	for _, id := range totalItemIDs {
		if !answeredMap[id] {
			remainingUnanswered++
		}
	}

	// 6. Tentukan status baru
	var newStatus string
	if remainingUnanswered == 0 {
		newStatus = "Completed"
	} else {
		newStatus = "In_Progress"
	}

	// 7. Jika hasil pencarian (Query) memberikan hasil kosong (empty set), update menjadi completed
	// (dan handle transisi In_Progress)
	if session.Status != newStatus {
		return r.data.db.WithContext(ctx).
			Model(&AssessmentSessionModel{}).
			Where("id = ?", sessionID).
			Update("status", newStatus).Error
	}

	return nil
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

	// Ambil data terbaru dari DB untuk memastikan ID dan field lainnya sinkron
	var latest AssessmentResultModel
	r.data.db.WithContext(ctx).First(&latest, "session_id = ? AND regulation_item_id = ?", result.SessionID, result.RegulationItemID)

	return toResultDomain(&latest), nil
}

func (r *assessmentResultRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&AssessmentResultModel{}, "id = ?", id).Error
}

func (r *assessmentResultRepo) DeleteBySessionAndRegulation(ctx context.Context, sessionID, regulationID uuid.UUID) error {
	// Delete results where item belongs to the regulation
	return r.data.db.WithContext(ctx).
		Where("session_id = ? AND regulation_item_id IN (SELECT id FROM regulation_items WHERE regulation_id = ?)", sessionID, regulationID).
		Delete(&AssessmentResultModel{}).Error
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
		IsActive:     ra.IsActive,
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
		IsActive:     ra.IsActive,
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
		Table("assessment_results").
		Select("assessment_results.compliance_status as status, COUNT(*) as count").
		Joins("JOIN regulation_items ON regulation_items.id = assessment_results.regulation_item_id").
		Where("assessment_results.session_id = ? AND regulation_items.regulation_id = ?", sessionID, regulationID).
		Group("assessment_results.compliance_status").
		Scan(&rows).Error
	if err != nil {
		r.log.Errorf("recalculation failed for session %s: %v", sessionID, err)
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
		ID:           uuid.New(),
		RegulationID: regulationID,
		SessionID:    sessionID,
		AmountPass:   ra.AmountPass,
		AmountFail:   ra.AmountFail,
		AmountNA:     ra.AmountNA,
		IsActive:     true,
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
		IsActive:     m.IsActive,
	}
}

func (r *regulationAssessmentRepo) DeleteBySessionAndRegulation(ctx context.Context, sessionID, regulationID uuid.UUID) error {
	return r.data.db.WithContext(ctx).
		Where("session_id = ? AND regulation_id = ?", sessionID, regulationID).
		Delete(&RegulationAssessmentModel{}).Error
}

func (r *regulationAssessmentRepo) Deactivate(ctx context.Context, sessionID, regulationID uuid.UUID) error {
	// Gunakan Map untuk memastikan GORM tidak mengabaikan nilai false (zero value)
	values := map[string]interface{}{
		"is_active": false,
	}

	result := r.data.db.WithContext(ctx).
		Model(&RegulationAssessmentModel{}).
		Where("session_id = ? AND regulation_id = ?", sessionID, regulationID).
		Updates(values)

	if result.Error != nil {
		return result.Error
	}

	// Jika tidak ada baris yang terupdate, berarti record belum ada, maka buat baru
	if result.RowsAffected == 0 {
		m := &RegulationAssessmentModel{
			ID:           uuid.New(),
			RegulationID: regulationID,
			SessionID:    sessionID,
			IsActive:     false,
		}
		return r.data.db.WithContext(ctx).Create(m).Error
	}

	return nil
}

func (r *regulationAssessmentRepo) Activate(ctx context.Context, sessionID, regulationID uuid.UUID) error {
	values := map[string]interface{}{
		"is_active": true,
	}

	result := r.data.db.WithContext(ctx).
		Model(&RegulationAssessmentModel{}).
		Where("session_id = ? AND regulation_id = ?", sessionID, regulationID).
		Updates(values)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		m := &RegulationAssessmentModel{
			ID:           uuid.New(),
			RegulationID: regulationID,
			SessionID:    sessionID,
			IsActive:     true,
		}
		return r.data.db.WithContext(ctx).Create(m).Error
	}

	return nil
}

func (r *regulationAssessmentRepo) IsActive(ctx context.Context, sessionID, regulationID uuid.UUID) (bool, error) {
	var m RegulationAssessmentModel
	if result := r.data.db.WithContext(ctx).
		Where("session_id = ? AND regulation_id = ?", sessionID, regulationID).
		First(&m); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Jika belum ada record summary, defaultnya Aktif (true)
			return true, nil
		}
		return false, result.Error
	}
