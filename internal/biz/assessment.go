package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// AssessmentUseCase menangani logika bisnis untuk modul Assessment.
type AssessmentUseCase struct {
	sessionRepo      AssessmentSessionRepo
	resultRepo       AssessmentResultRepo
	regulationAssRepo RegulationAssessmentRepo
	itemRepo         RegulationItemRepo
	log              *log.Helper
}

// NewAssessmentUseCase membuat instance baru.
func NewAssessmentUseCase(
	sessionRepo AssessmentSessionRepo,
	resultRepo AssessmentResultRepo,
	regulationAssRepo RegulationAssessmentRepo,
	itemRepo RegulationItemRepo,
	logger log.Logger,
) *AssessmentUseCase {
	return &AssessmentUseCase{
		sessionRepo:      sessionRepo,
		resultRepo:       resultRepo,
		regulationAssRepo: regulationAssRepo,
		itemRepo:         itemRepo,
		log:              log.NewHelper(logger),
	}
}

// --- AssessmentSession Use Cases ---

// CreateSession membuat sesi assessment baru.
// Setelah session dibuat, otomatis seed N/A untuk semua items yang propertinya
// tidak cocok dengan properti tenant, agar amount_na terhitung sejak awal.
func (uc *AssessmentUseCase) CreateSession(ctx context.Context, session *AssessmentSession) (*AssessmentSession, error) {
	if session.TenantID == uuid.Nil {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if session.Title == "" {
		return nil, fmt.Errorf("session title is required")
	}
	session.ID = uuid.New()
	session.Status = "Draft"

	created, err := uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	// Auto-seed N/A untuk items yang tidak relevan bagi tenant ini
	excludedItems, err := uc.itemRepo.FindExcludedByTenantID(ctx, session.TenantID)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("failed to fetch excluded items for tenant %s: %v", session.TenantID, err)
		return created, nil
	}

	// Kumpulkan regulation_id unik untuk recalculate setelah seed
	regulationIDs := make(map[uuid.UUID]struct{})
	for _, item := range excludedItems {
		naResult := &AssessmentResult{
			ID:               uuid.New(),
			SessionID:        created.ID,
			RegulationItemID: item.ID,
			ComplianceStatus: "N/A",
			Remarks:          "Item tidak relevan untuk properti tenant ini",
		}
		if _, err := uc.resultRepo.Upsert(ctx, naResult); err != nil {
			uc.log.WithContext(ctx).Warnf("failed to seed N/A for item %s: %v", item.ID, err)
		}
		regulationIDs[item.RegulationID] = struct{}{}
	}

	// Recalculate summary per regulasi
	for regID := range regulationIDs {
		if _, err := uc.regulationAssRepo.RecalculateForSession(ctx, created.ID, regID); err != nil {
			uc.log.WithContext(ctx).Warnf("recalculate failed for regulation %s: %v", regID, err)
		}
	}

	return created, nil
}

// GetSession mengambil sesi berdasarkan ID.
func (uc *AssessmentUseCase) GetSession(ctx context.Context, id uuid.UUID) (*AssessmentSession, error) {
	return uc.sessionRepo.FindByID(ctx, id)
}

// ListSessions mengembalikan semua sesi (atau filter per tenant jika tenantID diberikan).
func (uc *AssessmentUseCase) ListSessions(ctx context.Context, tenantID *uuid.UUID) ([]*AssessmentSession, error) {
	if tenantID != nil {
		return uc.sessionRepo.FindByTenantID(ctx, *tenantID)
	}
	return uc.sessionRepo.FindAll(ctx)
}

// UpdateSession memperbarui sesi assessment.
func (uc *AssessmentUseCase) UpdateSession(ctx context.Context, session *AssessmentSession) (*AssessmentSession, error) {
	return uc.sessionRepo.Update(ctx, session)
}

// DeleteSession menghapus sesi assessment.
func (uc *AssessmentUseCase) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return uc.sessionRepo.Delete(ctx, id)
}

// --- AssessmentResult Use Cases ---

// SubmitResult menyimpan atau memperbarui hasil assessment untuk satu item regulasi.
// Setelah upsert, recalculate ringkasan per regulasi.
func (uc *AssessmentUseCase) SubmitResult(ctx context.Context, result *AssessmentResult) (*AssessmentResult, error) {
	if result.SessionID == uuid.Nil || result.RegulationItemID == uuid.Nil {
		return nil, fmt.Errorf("session_id and regulation_item_id are required")
	}

	validStatuses := map[string]bool{"YES": true, "NO": true, "N/A": true}
	if !validStatuses[result.ComplianceStatus] {
		return nil, fmt.Errorf("compliance_status must be YES, NO, or N/A")
	}

	if result.ID == uuid.Nil {
		result.ID = uuid.New()
	}

	saved, err := uc.resultRepo.Upsert(ctx, result)
	if err != nil {
		return nil, err
	}

	// Ambil regulation_id dari item untuk recalculate summary.
	item, err := uc.itemRepo.FindByID(ctx, result.RegulationItemID)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("could not find regulation item for recalculation: %v", err)
		return saved, nil
	}

	if _, err := uc.regulationAssRepo.RecalculateForSession(ctx, result.SessionID, item.RegulationID); err != nil {
		uc.log.WithContext(ctx).Warnf("recalculation failed: %v", err)
	}

	return saved, nil
}

// GetResult mengambil hasil assessment berdasarkan ID.
func (uc *AssessmentUseCase) GetResult(ctx context.Context, id uuid.UUID) (*AssessmentResult, error) {
	return uc.resultRepo.FindByID(ctx, id)
}

// ListResultsBySession mengembalikan semua hasil dalam satu sesi.
func (uc *AssessmentUseCase) ListResultsBySession(ctx context.Context, sessionID uuid.UUID) ([]*AssessmentResult, error) {
	return uc.resultRepo.FindBySessionID(ctx, sessionID)
}

// DeleteResult menghapus hasil assessment.
func (uc *AssessmentUseCase) DeleteResult(ctx context.Context, id uuid.UUID) error {
	return uc.resultRepo.Delete(ctx, id)
}

// --- RegulationAssessment Use Cases ---

// GetRegulationAssessments mengembalikan ringkasan assessment per regulasi dalam satu sesi.
func (uc *AssessmentUseCase) GetRegulationAssessments(ctx context.Context, sessionID uuid.UUID) ([]*RegulationAssessment, error) {
	return uc.regulationAssRepo.FindBySessionID(ctx, sessionID)
}
