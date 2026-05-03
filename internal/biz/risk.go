package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type RiskUseCase struct {
	riskRepo        RiskRepo
	categoryRepo    RiskCategoryRepo
	tenantRepo      RiskCategoryTenantRepo
	log             *log.Helper
}

func NewRiskUseCase(riskRepo RiskRepo, categoryRepo RiskCategoryRepo, tenantRepo RiskCategoryTenantRepo, logger log.Logger) *RiskUseCase {
	return &RiskUseCase{
		riskRepo:        riskRepo,
		categoryRepo:    categoryRepo,
		tenantRepo:      tenantRepo,
		log:             log.NewHelper(logger),
	}
}

// --- Risk Category Use Cases ---

func (uc *RiskUseCase) CreateCategory(ctx context.Context, category *RiskCategory) (*RiskCategory, error) {
	if category.Title == "" {
		return nil, fmt.Errorf("category title is required")
	}
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return uc.categoryRepo.Create(ctx, category)
}

func (uc *RiskUseCase) GetCategory(ctx context.Context, id uuid.UUID) (*RiskCategory, error) {
	return uc.categoryRepo.FindByID(ctx, id)
}

func (uc *RiskUseCase) ListCategories(ctx context.Context) ([]*RiskCategory, error) {
	return uc.categoryRepo.FindAll(ctx)
}

func (uc *RiskUseCase) UpdateCategory(ctx context.Context, category *RiskCategory) (*RiskCategory, error) {
	return uc.categoryRepo.Update(ctx, category)
}

func (uc *RiskUseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return uc.categoryRepo.Delete(ctx, id)
}

// --- Risk Category Tenant Use Cases ---

func (uc *RiskUseCase) GetCategoryTenant(ctx context.Context, tenantID, categoryID uuid.UUID) (*RiskCategoryTenant, error) {
	tenantSetting, err := uc.tenantRepo.FindByTenantAndCategory(ctx, tenantID, categoryID)
	if err != nil {
		// If not found, return empty setting instead of error
		return &RiskCategoryTenant{
			RiskCategoryID: categoryID,
			TenantID:       tenantID,
		}, nil
	}
	return tenantSetting, nil
}

func (uc *RiskUseCase) SaveCategoryTenant(ctx context.Context, s *RiskCategoryTenant) (*RiskCategoryTenant, error) {
	if s.TenantID == uuid.Nil || s.RiskCategoryID == uuid.Nil {
		return nil, fmt.Errorf("tenant_id and risk_category_id are required")
	}
	
	existing, err := uc.tenantRepo.FindByTenantAndCategory(ctx, s.TenantID, s.RiskCategoryID)
	if err != nil {
		if s.ID == uuid.Nil {
			s.ID = uuid.New()
		}
		return uc.tenantRepo.Create(ctx, s)
	}
	
	s.ID = existing.ID
	return uc.tenantRepo.Update(ctx, s)
}

// --- Risk Use Cases ---

func (uc *RiskUseCase) CreateRisk(ctx context.Context, risk *Risk, tenantID uuid.UUID) (*Risk, error) {
	if risk.RiskTitle == "" {
		return nil, fmt.Errorf("risk title is required")
	}
	if risk.CategoryID == uuid.Nil {
		return nil, fmt.Errorf("category_id is required")
	}

	if risk.ID == uuid.Nil {
		risk.ID = uuid.New()
	}
	if risk.MitigationStatus == "" {
		risk.MitigationStatus = "belum direncanakan"
	}
	
	// 1. Simpan Risk ke tabel risks
	createdRisk, err := uc.riskRepo.Create(ctx, risk)
	if err != nil {
		return nil, err
	}

	// 2. Hubungkan ke tenant melalui RiskCategoryTenant
	// Cek apakah sudah ada setting appetite/tolerance untuk kategori ini di tenant ini
	// (untuk inheritance default value)
	var appetite, tolerance string
	existingSettings, err := uc.tenantRepo.FindByTenantAndCategory(ctx, tenantID, risk.CategoryID)
	if err == nil {
		appetite = existingSettings.Appetite
		tolerance = existingSettings.Tolerance
	}

	_, err = uc.tenantRepo.Create(ctx, &RiskCategoryTenant{
		ID:             uuid.New(),
		RiskID:         createdRisk.ID,
		RiskCategoryID: risk.CategoryID,
		TenantID:       tenantID,
		Appetite:       appetite,
		Tolerance:      tolerance,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungkan risiko ke tenant: %v", err)
	}

	return createdRisk, nil
}

func (uc *RiskUseCase) GetRisk(ctx context.Context, id uuid.UUID) (*Risk, error) {
	return uc.riskRepo.FindByID(ctx, id)
}

func (uc *RiskUseCase) ListRisks(ctx context.Context, tenantID uuid.UUID) ([]*Risk, error) {
	return uc.riskRepo.FindAll(ctx, tenantID)
}

func (uc *RiskUseCase) UpdateRisk(ctx context.Context, risk *Risk) (*Risk, error) {
	return uc.riskRepo.Update(ctx, risk)
}

func (uc *RiskUseCase) DeleteRisk(ctx context.Context, id uuid.UUID) error {
	return uc.riskRepo.Delete(ctx, id)
}
