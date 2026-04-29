package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type RiskUseCase struct {
	riskRepo     RiskRepo
	categoryRepo RiskCategoryRepo
	log          *log.Helper
}

func NewRiskUseCase(riskRepo RiskRepo, categoryRepo RiskCategoryRepo, logger log.Logger) *RiskUseCase {
	return &RiskUseCase{
		riskRepo:     riskRepo,
		categoryRepo: categoryRepo,
		log:          log.NewHelper(logger),
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

// --- Risk Use Cases ---

func (uc *RiskUseCase) CreateRisk(ctx context.Context, risk *Risk) (*Risk, error) {
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
	return uc.riskRepo.Create(ctx, risk)
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
