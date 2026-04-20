package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// TenantUseCase menangani logika bisnis untuk Tenant.
type TenantUseCase struct {
	repo TenantRepo
	log  *log.Helper
}

// NewTenantUseCase membuat instance TenantUseCase baru.
func NewTenantUseCase(repo TenantRepo, logger log.Logger) *TenantUseCase {
	return &TenantUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// CreateTenant membuat tenant baru.
func (uc *TenantUseCase) CreateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	if tenant.Name == "" {
		return nil, fmt.Errorf("tenant name is required")
	}
	tenant.ID = uuid.New()
	result, err := uc.repo.Create(ctx, tenant)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("CreateTenant failed: %v", err)
		return nil, err
	}
	return result, nil
}

// GetTenant mengambil tenant berdasarkan ID.
func (uc *TenantUseCase) GetTenant(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return uc.repo.FindByID(ctx, id)
}

// ListTenants mengembalikan semua tenant.
func (uc *TenantUseCase) ListTenants(ctx context.Context) ([]*Tenant, error) {
	return uc.repo.FindAll(ctx)
}

// UpdateTenant memperbarui data tenant.
func (uc *TenantUseCase) UpdateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	if tenant.ID == uuid.Nil {
		return nil, fmt.Errorf("tenant id is required")
	}
	return uc.repo.Update(ctx, tenant)
}

// DeleteTenant menghapus tenant berdasarkan ID.
func (uc *TenantUseCase) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// --- PropertyUseCase ---

// PropertyUseCase menangani logika bisnis untuk Property.
type PropertyUseCase struct {
	repo TenantRepo
	pr   PropertyRepo
	log  *log.Helper
}

// NewPropertyUseCase membuat instance baru.
func NewPropertyUseCase(pr PropertyRepo, logger log.Logger) *PropertyUseCase {
	return &PropertyUseCase{
		pr:  pr,
		log: log.NewHelper(logger),
	}
}

// CreateProperty membuat property baru.
func (uc *PropertyUseCase) CreateProperty(ctx context.Context, p *Property) (*Property, error) {
	if p.Name == "" {
		return nil, fmt.Errorf("property name is required")
	}
	p.ID = uuid.New()
	return uc.pr.Create(ctx, p)
}

// GetProperty mengambil property berdasarkan ID.
func (uc *PropertyUseCase) GetProperty(ctx context.Context, id uuid.UUID) (*Property, error) {
	return uc.pr.FindByID(ctx, id)
}

// ListProperties mengembalikan semua property.
func (uc *PropertyUseCase) ListProperties(ctx context.Context) ([]*Property, error) {
	return uc.pr.FindAll(ctx)
}

// UpdateProperty memperbarui property.
func (uc *PropertyUseCase) UpdateProperty(ctx context.Context, p *Property) (*Property, error) {
	return uc.pr.Update(ctx, p)
}

// DeleteProperty menghapus property berdasarkan ID.
func (uc *PropertyUseCase) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	return uc.pr.Delete(ctx, id)
}
