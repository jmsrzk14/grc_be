package data

import (
	"context"
	"errors"

	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Tenant Repository Implementation ---

type tenantRepo struct {
	data *Data
	log  *log.Helper
}

// NewTenantRepo membuat instance repository Tenant.
func NewTenantRepo(data *Data, logger log.Logger) biz.TenantRepo {
	return &tenantRepo{data: data, log: log.NewHelper(logger)}
}

func (r *tenantRepo) Create(ctx context.Context, tenant *biz.Tenant) (*biz.Tenant, error) {
	m := &TenantModel{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Type:      tenant.Type,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt,
	}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	tenant.CreatedAt = m.CreatedAt
	return tenant, nil
}

func (r *tenantRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Tenant, error) {
	var m TenantModel
	if result := r.data.db.WithContext(ctx).First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toTenantDomain(&m), nil
}

func (r *tenantRepo) FindAll(ctx context.Context) ([]*biz.Tenant, error) {
	var models []*TenantModel
	if result := r.data.db.WithContext(ctx).Find(&models); result.Error != nil {
		return nil, result.Error
	}
	tenants := make([]*biz.Tenant, 0, len(models))
	for _, m := range models {
		tenants = append(tenants, toTenantDomain(m))
	}
	return tenants, nil
}

func (r *tenantRepo) Update(ctx context.Context, tenant *biz.Tenant) (*biz.Tenant, error) {
	m := &TenantModel{
		ID:     tenant.ID,
		Name:   tenant.Name,
		Type:   tenant.Type,
		Status: tenant.Status,
	}
	if result := r.data.db.WithContext(ctx).Save(m); result.Error != nil {
		return nil, result.Error
	}
	return toTenantDomain(m), nil
}

func (r *tenantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&TenantModel{}, "id = ?", id).Error
}

func toTenantDomain(m *TenantModel) *biz.Tenant {
	return &biz.Tenant{
		ID:        m.ID,
		Name:      m.Name,
		Type:      m.Type,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}
}

// --- Property Repository Implementation ---

type propertyRepo struct {
	data *Data
	log  *log.Helper
}

// NewPropertyRepo membuat instance repository Property.
func NewPropertyRepo(data *Data, logger log.Logger) biz.PropertyRepo {
	return &propertyRepo{data: data, log: log.NewHelper(logger)}
}

func (r *propertyRepo) Create(ctx context.Context, p *biz.Property) (*biz.Property, error) {
	m := &PropertyModel{ID: p.ID, Name: p.Name, Description: p.Description}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	return p, nil
}

func (r *propertyRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Property, error) {
	var m PropertyModel
	if result := r.data.db.WithContext(ctx).First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return &biz.Property{ID: m.ID, Name: m.Name, Description: m.Description}, nil
}

func (r *propertyRepo) FindAll(ctx context.Context) ([]*biz.Property, error) {
	var models []*PropertyModel
	if result := r.data.db.WithContext(ctx).Find(&models); result.Error != nil {
		return nil, result.Error
	}
	props := make([]*biz.Property, 0, len(models))
	for _, m := range models {
		props = append(props, &biz.Property{ID: m.ID, Name: m.Name, Description: m.Description})
	}
	return props, nil
}

func (r *propertyRepo) Update(ctx context.Context, p *biz.Property) (*biz.Property, error) {
	m := &PropertyModel{ID: p.ID, Name: p.Name, Description: p.Description}
	if result := r.data.db.WithContext(ctx).Save(m); result.Error != nil {
		return nil, result.Error
	}
	return p, nil
}

func (r *propertyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&PropertyModel{}, "id = ?", id).Error
}

// --- TenantProperty Repository Implementation ---

type tenantPropertyRepo struct {
	data *Data
	log  *log.Helper
}

// NewTenantPropertyRepo membuat instance repository TenantProperty.
func NewTenantPropertyRepo(data *Data, logger log.Logger) biz.TenantPropertyRepo {
	return &tenantPropertyRepo{data: data, log: log.NewHelper(logger)}
}

func (r *tenantPropertyRepo) Create(ctx context.Context, tp *biz.TenantProperty) (*biz.TenantProperty, error) {
	m := &TenantPropertyModel{ID: tp.ID, TenantID: tp.TenantID, PropertyID: tp.PropertyID}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	return tp, nil
}

func (r *tenantPropertyRepo) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*biz.TenantProperty, error) {
	var models []*TenantPropertyModel
	if result := r.data.db.WithContext(ctx).Find(&models, "tenant_id = ?", tenantID); result.Error != nil {
		return nil, result.Error
	}
	tps := make([]*biz.TenantProperty, 0, len(models))
	for _, m := range models {
		tps = append(tps, &biz.TenantProperty{ID: m.ID, TenantID: m.TenantID, PropertyID: m.PropertyID})
	}
	return tps, nil
}

func (r *tenantPropertyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&TenantPropertyModel{}, "id = ?", id).Error
}
