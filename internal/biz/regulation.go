package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// RegulationUseCase menangani logika bisnis untuk Regulation.
type RegulationUseCase struct {
	repo              RegulationRepo
	itemRepo          RegulationItemRepo
	mappingRepo       RegulationPropertyMappingRepo
	tenantRegRepo     TenantRegulationRepo
	sessionRepo       AssessmentSessionRepo
	resultRepo        AssessmentResultRepo
	regulationAssRepo RegulationAssessmentRepo
	log               *log.Helper
}

// NewRegulationUseCase membuat instance baru.
func NewRegulationUseCase(
	repo RegulationRepo,
	itemRepo RegulationItemRepo,
	mappingRepo RegulationPropertyMappingRepo,
	tenantRegRepo TenantRegulationRepo,
	sessionRepo AssessmentSessionRepo,
	resultRepo AssessmentResultRepo,
	regulationAssRepo RegulationAssessmentRepo,
	logger log.Logger,
) *RegulationUseCase {
	return &RegulationUseCase{
		repo:              repo,
		itemRepo:          itemRepo,
		mappingRepo:       mappingRepo,
		tenantRegRepo:     tenantRegRepo,
		sessionRepo:       sessionRepo,
		resultRepo:        resultRepo,
		regulationAssRepo: regulationAssRepo,
		log:               log.NewHelper(logger),
	}
}

// CreateRegulation membuat regulasi baru.
func (uc *RegulationUseCase) CreateRegulation(ctx context.Context, r *Regulation) (*Regulation, error) {
	if r.Title == "" {
		return nil, fmt.Errorf("regulation title is required")
	}
	r.ID = uuid.New()
	created, err := uc.repo.Create(ctx, r)
	if err != nil {
		return nil, err
	}

	// Jika kategori internal dan ada tenant_id, simpan mapping ke tenant_regulation
	if created.Category == "Internal" && created.TenantID != uuid.Nil {
		_, err = uc.tenantRegRepo.Create(ctx, &TenantRegulation{
			ID:           uuid.New(),
			TenantID:     created.TenantID,
			RegulationID: created.ID,
		})
		if err != nil {
			uc.log.Errorf("failed to create tenant regulation mapping: %v", err)
		}
	}

	return created, nil
}

// UpsertRegulation mencari regulasi berdasarkan judul, jika ada diupdate, jika tidak dicreate.
func (uc *RegulationUseCase) UpsertRegulation(ctx context.Context, r *Regulation) (*Regulation, error) {
	if r.Title == "" {
		return nil, fmt.Errorf("regulation title is required")
	}

	existing, err := uc.repo.FindByTitle(ctx, r.Title)
	if err == nil {
		// Update existing
		r.ID = existing.ID
		return uc.repo.Update(ctx, r)
	}

	// Create new
	r.ID = uuid.New()
	created, err := uc.repo.Create(ctx, r)
	if err != nil {
		return nil, err
	}

	// Jika kategori internal dan ada tenant_id, simpan mapping ke tenant_regulation
	if created.Category == "Internal" && created.TenantID != uuid.Nil {
		_, err = uc.tenantRegRepo.Create(ctx, &TenantRegulation{
			ID:           uuid.New(),
			TenantID:     created.TenantID,
			RegulationID: created.ID,
		})
		if err != nil {
			uc.log.Errorf("failed to create tenant regulation mapping: %v", err)
		}
	}

	return created, nil
}

// GetRegulation mengambil regulasi berdasarkan ID.
func (uc *RegulationUseCase) GetRegulation(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*Regulation, error) {
	return uc.repo.FindByID(ctx, id, tenantID)
}

// ListRegulations mengembalikan semua regulasi.
func (uc *RegulationUseCase) ListRegulations(ctx context.Context, tenantID uuid.UUID) ([]*Regulation, error) {
	return uc.repo.FindAll(ctx, tenantID)
}

// UpdateRegulation memperbarui data regulasi.
func (uc *RegulationUseCase) UpdateRegulation(ctx context.Context, r *Regulation) (*Regulation, error) {
	return uc.repo.Update(ctx, r)
}

// DeleteRegulation menghapus regulasi berdasarkan ID.
func (uc *RegulationUseCase) DeleteRegulation(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// AssignTenantToRegulation menghubungkan tenant ke suatu regulasi.
func (uc *RegulationUseCase) AssignTenantToRegulation(ctx context.Context, regulationID, tenantID uuid.UUID) error {
	_, err := uc.tenantRegRepo.Create(ctx, &TenantRegulation{
		ID:           uuid.New(),
		TenantID:     tenantID,
		RegulationID: regulationID,
	})
	return err
}

// RevokeTenantFromRegulation mencopot hak akses tenant ke regulasi dan menghapus assessment periode berjalan.
func (uc *RegulationUseCase) RevokeTenantFromRegulation(ctx context.Context, regulationID, tenantID uuid.UUID) error {
	// 1. Hapus mapping akses
	if err := uc.tenantRegRepo.Delete(ctx, tenantID, regulationID); err != nil {
		return err
	}

	// 2. Cari assessment session untuk periode berjalan (tahun ini)
	currentYear := time.Now().Year()
	sessions, err := uc.sessionRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("failed to find sessions for tenant %s: %v", tenantID, err)
		return nil // Tetap sukses hapus mapping
	}

	for _, session := range sessions {
		if session.PeriodYear == currentYear {
			// 3. Hapus assessment results untuk regulasi ini di session tersebut
			if err := uc.resultRepo.DeleteBySessionAndRegulation(ctx, session.ID, regulationID); err != nil {
				uc.log.WithContext(ctx).Errorf("failed to delete results for session %s, reg %s: %v", session.ID, regulationID, err)
			}
			// 4. Hapus summary assessment untuk regulasi ini di session tersebut
			if err := uc.regulationAssRepo.DeleteBySessionAndRegulation(ctx, session.ID, regulationID); err != nil {
				uc.log.WithContext(ctx).Errorf("failed to delete summary for session %s, reg %s: %v", session.ID, regulationID, err)
			}
		}
	}

	return nil
}

// GetAssignedTenants mengembalikan daftar tenant yang memiliki akses ke regulasi tertentu.
func (uc *RegulationUseCase) GetAssignedTenants(ctx context.Context, regulationID uuid.UUID) ([]*TenantRegulation, error) {
	return uc.tenantRegRepo.FindByRegulationID(ctx, regulationID)
}

// --- RegulationItem Use Cases ---

// CreateRegulationItem membuat item/pasal baru dalam regulasi.
func (uc *RegulationUseCase) CreateRegulationItem(ctx context.Context, item *RegulationItem) (*RegulationItem, error) {
	if item.RegulationID == uuid.Nil {
		return nil, fmt.Errorf("regulation_id is required")
	}

	// Auto-increment if ItemCode is 0 (applies to all categories for manual entry)
	if item.ItemCode == 0 {
		max, err := uc.itemRepo.GetMaxItemCode(ctx, item.RegulationID)
		if err != nil {
			return nil, err
		}
		item.ItemCode = max + 1
	}

	item.ID = uuid.New()
	return uc.itemRepo.Create(ctx, item)
}

// UpsertRegulationItem mencari item berdasarkan item_code dalam regulasi, jika ada diupdate, jika tidak dicreate.
func (uc *RegulationUseCase) UpsertRegulationItem(ctx context.Context, item *RegulationItem) (*RegulationItem, error) {
	if item.RegulationID == uuid.Nil {
		return nil, fmt.Errorf("regulation_id is required")
	}

	if item.ItemCode != 0 {
		existing, err := uc.itemRepo.FindByRegulationIDAndItemCode(ctx, item.RegulationID, item.ItemCode)
		if err == nil {
			// Update existing
			item.ID = existing.ID
			return uc.itemRepo.Update(ctx, item)
		}
	} else {
		// For Upsert (Import), we only auto-increment for Internal regulations if ItemCode is 0.
		// Fetch regulation to check category
		reg, err := uc.repo.FindByID(ctx, item.RegulationID, uuid.Nil)
		if err != nil {
			return nil, err
		}
		
		if reg.Category == "Internal" {
			max, err := uc.itemRepo.GetMaxItemCode(ctx, item.RegulationID)
			if err != nil {
				return nil, err
			}
			item.ItemCode = max + 1
		} else {
			return nil, fmt.Errorf("item_code is required for non-internal regulation imports")
		}
	}

	// Create new
	item.ID = uuid.New()
	return uc.itemRepo.Create(ctx, item)
}

// GetRegulationItem mengambil item berdasarkan ID.
func (uc *RegulationUseCase) GetRegulationItem(ctx context.Context, id uuid.UUID) (*RegulationItem, error) {
	return uc.itemRepo.FindByID(ctx, id)
}

// ListRegulationItems mengembalikan semua item dalam satu regulasi, opsional difilter berdasarkan tenant.
func (uc *RegulationUseCase) ListRegulationItems(ctx context.Context, regulationID uuid.UUID, tenantID uuid.UUID) ([]*RegulationItem, error) {
	return uc.itemRepo.FindByRegulationID(ctx, regulationID, tenantID)
}

// UpdateRegulationItem memperbarui item regulasi.
func (uc *RegulationUseCase) UpdateRegulationItem(ctx context.Context, item *RegulationItem) (*RegulationItem, error) {
	return uc.itemRepo.Update(ctx, item)
}

// DeleteRegulationItem menghapus item regulasi.
func (uc *RegulationUseCase) DeleteRegulationItem(ctx context.Context, id uuid.UUID) error {
	return uc.itemRepo.Delete(ctx, id)
}

// --- RegulationPropertyMapping Use Cases ---

// AddPropertyToRegulation menambahkan mapping antara regulasi dan property.
func (uc *RegulationUseCase) AddPropertyToRegulation(ctx context.Context, mapping *RegulationPropertyMapping) (*RegulationPropertyMapping, error) {
	if mapping.RegulationID == uuid.Nil || mapping.PropertyID == uuid.Nil {
		return nil, fmt.Errorf("regulation_id and property_id are required")
	}
	mapping.ID = uuid.New()
	return uc.mappingRepo.Create(ctx, mapping)
}

// ListRegulationMappings mengembalikan semua mapping untuk suatu regulasi.
func (uc *RegulationUseCase) ListRegulationMappings(ctx context.Context, regulationID uuid.UUID) ([]*RegulationPropertyMapping, error) {
	return uc.mappingRepo.FindByRegulationID(ctx, regulationID)
}

// DeleteRegulationMapping menghapus mapping berdasarkan ID.
func (uc *RegulationUseCase) DeleteRegulationMapping(ctx context.Context, id uuid.UUID) error {
	return uc.mappingRepo.Delete(ctx, id)
}
