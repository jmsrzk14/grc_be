package biz

import (
	"context"

	"github.com/google/uuid"
)

// --- Tenant Repository Interface ---

// TenantRepo mendefinisikan kontrak akses data untuk Tenant.
type TenantRepo interface {
	Create(ctx context.Context, tenant *Tenant) (*Tenant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	FindAll(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) (*Tenant, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PropertyRepo mendefinisikan kontrak akses data untuk Property.
type PropertyRepo interface {
	Create(ctx context.Context, property *Property) (*Property, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Property, error)
	FindAll(ctx context.Context) ([]*Property, error)
	Update(ctx context.Context, property *Property) (*Property, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// TenantPropertyRepo mendefinisikan kontrak untuk mapping Tenant-Property.
type TenantPropertyRepo interface {
	Create(ctx context.Context, tp *TenantProperty) (*TenantProperty, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*TenantProperty, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// --- Regulation Repository Interface ---

// RegulationRepo mendefinisikan kontrak akses data untuk Regulation.
type RegulationRepo interface {
	Create(ctx context.Context, r *Regulation) (*Regulation, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Regulation, error)
	FindAll(ctx context.Context) ([]*Regulation, error)
	Update(ctx context.Context, r *Regulation) (*Regulation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// RegulationItemRepo mendefinisikan kontrak akses data untuk RegulationItem.
type RegulationItemRepo interface {
	Create(ctx context.Context, item *RegulationItem) (*RegulationItem, error)
	FindByID(ctx context.Context, id uuid.UUID) (*RegulationItem, error)
	FindByRegulationID(ctx context.Context, regulationID uuid.UUID) ([]*RegulationItem, error)
	Update(ctx context.Context, item *RegulationItem) (*RegulationItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// RegulationPropertyMappingRepo mendefinisikan kontrak untuk mapping Regulasi-Property.
type RegulationPropertyMappingRepo interface {
	Create(ctx context.Context, mapping *RegulationPropertyMapping) (*RegulationPropertyMapping, error)
	FindByRegulationID(ctx context.Context, regulationID uuid.UUID) ([]*RegulationPropertyMapping, error)
	FindByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*RegulationPropertyMapping, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// --- Assessment Repository Interface ---

// AssessmentSessionRepo mendefinisikan kontrak akses data untuk AssessmentSession.
type AssessmentSessionRepo interface {
	Create(ctx context.Context, session *AssessmentSession) (*AssessmentSession, error)
	FindByID(ctx context.Context, id uuid.UUID) (*AssessmentSession, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*AssessmentSession, error)
	FindAll(ctx context.Context) ([]*AssessmentSession, error)
	Update(ctx context.Context, session *AssessmentSession) (*AssessmentSession, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// AssessmentResultRepo mendefinisikan kontrak akses data untuk AssessmentResult.
type AssessmentResultRepo interface {
	Create(ctx context.Context, result *AssessmentResult) (*AssessmentResult, error)
	FindByID(ctx context.Context, id uuid.UUID) (*AssessmentResult, error)
	FindBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*AssessmentResult, error)
	Upsert(ctx context.Context, result *AssessmentResult) (*AssessmentResult, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// RegulationAssessmentRepo mendefinisikan kontrak untuk ringkasan assessment.
type RegulationAssessmentRepo interface {
	Create(ctx context.Context, ra *RegulationAssessment) (*RegulationAssessment, error)
	FindBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*RegulationAssessment, error)
	FindByID(ctx context.Context, id uuid.UUID) (*RegulationAssessment, error)
	Update(ctx context.Context, ra *RegulationAssessment) (*RegulationAssessment, error)
	RecalculateForSession(ctx context.Context, sessionID uuid.UUID, regulationID uuid.UUID) (*RegulationAssessment, error)
}
