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
		_, err = uc.tenantRegRepo.Upsert(ctx, &TenantRegulation{
			ID:           uuid.New(),
			TenantID:     created.TenantID,
			RegulationID: created.ID,
		})
		if err != nil {
			uc.log.Errorf("failed to upsert tenant regulation mapping: %v", err)
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
		_, err = uc.tenantRegRepo.Upsert(ctx, &TenantRegulation{
			ID:           uuid.New(),
			TenantID:     created.TenantID,
			RegulationID: created.ID,
		})
		if err != nil {
			uc.log.Errorf("failed to upsert tenant regulation mapping: %v", err)
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
	// 1. Pastikan mapping ada di tenant_regulations (untuk visibilitas umum)
	_, err := uc.tenantRegRepo.Upsert(ctx, &TenantRegulation{
		ID:           uuid.New(),
		TenantID:     tenantID,
		RegulationID: regulationID,
	})
	if err != nil {
		return err
	}

	// 2. Aktifkan kembali di sesi periode tahun ini (jika ada)
	currentYear := time.Now().Year()
	session, err := uc.sessionRepo.FindByTenantAndYear(ctx, tenantID, currentYear)
	if err == nil {
		// Jika ada sesi, aktifkan status is_active di summary
		_ = uc.regulationAssRepo.Activate(ctx, session.ID, regulationID)

		// Cek apakah session menjadi Completed
		if err := uc.sessionRepo.CheckAndComplete(ctx, session.ID); err != nil {
			uc.log.WithContext(ctx).Warnf("check and complete session failed: %v", err)
		}
	}

	return nil
}

// RevokeTenantFromRegulation mencabut akses regulasi dari tenant secara non-destruktif.
// Pemetaan di tenant_regulations tetap dipertahankan agar tenant masih bisa melihat regulasi,
// namun status aktif di summary penilaian (regulation_assesments) diubah menjadi false.
func (uc *RegulationUseCase) RevokeTenantFromRegulation(ctx context.Context, regulationID, tenantID uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("Starting RevokeTenantFromRegulation: tenantID=%s, regulationID=%s", tenantID, regulationID)
	
	// 1. Cari assessment session untuk tenant ini di periode tahun ini
	currentYear := time.Now().Year()
	session, err := uc.sessionRepo.FindByTenantAndYear(ctx, tenantID, currentYear)
	if err != nil {
		// Jika tidak ada sesi untuk tahun ini, tidak ada yang perlu dinonaktifkan di summary
		uc.log.WithContext(ctx).Warnf("no active session found for tenant %s in year %d, skipping summary deactivation: %v", tenantID, currentYear, err)
		return nil
	}

	// 2. Deactivate is_active di regulation_assesments untuk sesi tahun ini
	// Ini memastikan tombol penilaian hilang di frontend dan data tidak dihitung di dashboard.
	uc.log.WithContext(ctx).Infof("Deactivating assessment summary: tenant %s, reg %s, session %s, year %d", tenantID, regulationID, session.ID, currentYear)
	if err := uc.regulationAssRepo.Deactivate(ctx, session.ID, regulationID); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to deactivate summary for session %s, reg %s: %v", session.ID, regulationID, err)
		return err
	}

	// Cek apakah session menjadi Completed (karena regulasi yang belum dijawab dicabut)
	if err := uc.sessionRepo.CheckAndComplete(ctx, session.ID); err != nil {
		uc.log.WithContext(ctx).Warnf("check and complete session failed: %v", err)
	}

	uc.log.WithContext(ctx).Infof("Successfully deactivated assessment for tenant %s, regulation %s", tenantID, regulationID)
	return nil
}

func (uc *RegulationUseCase) GetAssignedTenants(ctx context.Context, regulationID uuid.UUID) ([]*TenantRegulation, error) {
	mappings, err := uc.tenantRegRepo.FindByRegulationID(ctx, regulationID)
	if err != nil {
		return nil, err
	}

	for _, m := range mappings {
		// Cek status di sesi periode tahun ini
		currentYear := time.Now().Year()
		session, err := uc.sessionRepo.FindByTenantAndYear(ctx, m.TenantID, currentYear)
		if err == nil {
			// Cek apakah aktif di sesi tersebut
			active, err := uc.regulationAssRepo.IsActive(ctx, session.ID, regulationID)
			if err == nil {
				m.IsActive = active
			} else {
				m.IsActive = true // Default jika error
			}
		} else {
			// Jika belum ada sesi untuk tahun ini, defaultnya adalah Aktif (karena sudah di-assign)
			m.IsActive = true
		}
	}

	return mappings, nil
}

func (uc *RegulationUseCase) CreateRegulationItem(ctx context.Context, item *RegulationItem) (*RegulationItem, error) {
	if item.RegulationID == uuid.Nil {
		return nil, fmt.Errorf("regulation_id is required")
	}

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
			item.ID = existing.ID
			return uc.itemRepo.Update(ctx, item)
		}
	} else {
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

	item.ID = uuid.New()
	return uc.itemRepo.Create(ctx, item)
}

func (uc *RegulationUseCase) GetRegulationItem(ctx context.Context, id uuid.UUID) (*RegulationItem, error) {
	return uc.itemRepo.FindByID(ctx, id)
}

func (uc *RegulationUseCase) ListRegulationItems(ctx context.Context, regulationID uuid.UUID, tenantID uuid.UUID) ([]*RegulationItem, error) {
	return uc.itemRepo.FindByRegulationID(ctx, regulationID, tenantID)
}

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
