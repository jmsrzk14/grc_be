package data

import (
	"context"
	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type riskRepo struct {
	data *Data
	log  *log.Helper
}

func NewRiskRepo(data *Data, logger log.Logger) biz.RiskRepo {
	return &riskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *riskRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Risk, error) {
	var model RiskModel
	if err := r.data.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &biz.Risk{
		ID:                 model.ID,
		RiskTitle:          model.RiskTitle,
		RiskDescription:    model.RiskDescription,
		CategoryID:         model.CategoryID,
		LikelihoodInherent: model.LikelihoodInherent,
		ImpactInherent:     model.ImpactInherent,
		LikelihoodResidual: model.LikelihoodResidual,
		ImpactResidual:     model.ImpactResidual,
		MitigationPlan:     model.MitigationPlan,
		MitigationStatus:   model.MitigationStatus,
	}, nil
}

func (r *riskRepo) FindAll(ctx context.Context, tenantID uuid.UUID) ([]*biz.Risk, error) {
	var models []RiskModel
	// Join with risk_category_tenants to filter by tenant and get appetite/tolerance if needed
	if err := r.data.db.WithContext(ctx).
		Joins("JOIN risk_category_tenants s ON risks.id = s.risk_id").
		Where("s.tenant_id = ?", tenantID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	risks := make([]*biz.Risk, 0, len(models))
	for _, model := range models {
		risks = append(risks, &biz.Risk{
			ID:                 model.ID,
			RiskTitle:          model.RiskTitle,
			RiskDescription:    model.RiskDescription,
			CategoryID:         model.CategoryID,
			LikelihoodInherent: model.LikelihoodInherent,
			ImpactInherent:     model.ImpactInherent,
			LikelihoodResidual: model.LikelihoodResidual,
			ImpactResidual:     model.ImpactResidual,
			MitigationPlan:     model.MitigationPlan,
			MitigationStatus:   model.MitigationStatus,
		})
	}
	return risks, nil
}

func (r *riskRepo) Create(ctx context.Context, risk *biz.Risk) (*biz.Risk, error) {
	model := &RiskModel{
		ID:                 risk.ID,
		RiskTitle:          risk.RiskTitle,
		RiskDescription:    risk.RiskDescription,
		CategoryID:         risk.CategoryID,
		LikelihoodInherent: risk.LikelihoodInherent,
		ImpactInherent:     risk.ImpactInherent,
		LikelihoodResidual: risk.LikelihoodResidual,
		ImpactResidual:     risk.ImpactResidual,
		MitigationPlan:     risk.MitigationPlan,
		MitigationStatus:   risk.MitigationStatus,
	}
	if err := r.data.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return risk, nil
}

func (r *riskRepo) Update(ctx context.Context, risk *biz.Risk) (*biz.Risk, error) {
	model := &RiskModel{
		ID:                 risk.ID,
		RiskTitle:          risk.RiskTitle,
		RiskDescription:    risk.RiskDescription,
		CategoryID:         risk.CategoryID,
		LikelihoodInherent: risk.LikelihoodInherent,
		ImpactInherent:     risk.ImpactInherent,
		LikelihoodResidual: risk.LikelihoodResidual,
		ImpactResidual:     risk.ImpactResidual,
		MitigationPlan:     risk.MitigationPlan,
		MitigationStatus:   risk.MitigationStatus,
	}
	if err := r.data.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}
	return risk, nil
}

func (r *riskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.data.db.WithContext(ctx).Delete(&RiskModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

type riskCategoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewRiskCategoryRepo(data *Data, logger log.Logger) biz.RiskCategoryRepo {
	return &riskCategoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *riskCategoryRepo) Create(ctx context.Context, category *biz.RiskCategory) (*biz.RiskCategory, error) {
	model := &RiskCategoryModel{
		ID:    category.ID,
		Title: category.Title,
	}
	if err := r.data.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *riskCategoryRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.RiskCategory, error) {
	var model RiskCategoryModel
	if err := r.data.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &biz.RiskCategory{
		ID:    model.ID,
		Title: model.Title,
	}, nil
}

func (r *riskCategoryRepo) FindAll(ctx context.Context) ([]*biz.RiskCategory, error) {
	var models []RiskCategoryModel
	if err := r.data.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	categories := make([]*biz.RiskCategory, 0, len(models))
	for _, model := range models {
		categories = append(categories, &biz.RiskCategory{
			ID:    model.ID,
			Title: model.Title,
		})
	}
	return categories, nil
}

func (r *riskCategoryRepo) Update(ctx context.Context, category *biz.RiskCategory) (*biz.RiskCategory, error) {
	model := &RiskCategoryModel{
		ID:    category.ID,
		Title: category.Title,
	}
	if err := r.data.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *riskCategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.data.db.WithContext(ctx).Delete(&RiskCategoryModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

type riskCategoryTenantRepo struct {
	data *Data
	log  *log.Helper
}

func NewRiskCategoryTenantRepo(data *Data, logger log.Logger) biz.RiskCategoryTenantRepo {
	return &riskCategoryTenantRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *riskCategoryTenantRepo) Create(ctx context.Context, s *biz.RiskCategoryTenant) (*biz.RiskCategoryTenant, error) {
	var riskID *uuid.UUID
	if s.RiskID != uuid.Nil {
		riskID = &s.RiskID
	}
	model := &RiskCategoryTenantModel{
		ID:             s.ID,
		RiskID:         riskID,
		RiskCategoryID: s.RiskCategoryID,
		TenantID:       s.TenantID,
		Appetite:       s.Appetite,
		Tolerance:      s.Tolerance,
	}
	if err := r.data.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *riskCategoryTenantRepo) FindByTenantAndCategory(ctx context.Context, tenantID, categoryID uuid.UUID) (*biz.RiskCategoryTenant, error) {
	var model RiskCategoryTenantModel
	if err := r.data.db.WithContext(ctx).Where("tenant_id = ? AND risk_category_id = ?", tenantID, categoryID).First(&model).Error; err != nil {
		return nil, err
	}
	riskID := uuid.Nil
	if model.RiskID != nil {
		riskID = *model.RiskID
	}
	return &biz.RiskCategoryTenant{
		ID:             model.ID,
		RiskID:         riskID,
		RiskCategoryID: model.RiskCategoryID,
		TenantID:       model.TenantID,
		Appetite:       model.Appetite,
		Tolerance:      model.Tolerance,
	}, nil
}

func (r *riskCategoryTenantRepo) Update(ctx context.Context, s *biz.RiskCategoryTenant) (*biz.RiskCategoryTenant, error) {
	var riskID *uuid.UUID
	if s.RiskID != uuid.Nil {
		riskID = &s.RiskID
	}
	model := &RiskCategoryTenantModel{
		ID:             s.ID,
		RiskID:         riskID,
		RiskCategoryID: s.RiskCategoryID,
		TenantID:       s.TenantID,
		Appetite:       s.Appetite,
		Tolerance:      s.Tolerance,
	}
	if err := r.data.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}
	return s, nil
}
