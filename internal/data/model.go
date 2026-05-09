package data

import (
	"time"

	"github.com/google/uuid"
)

// --- GORM Models (representasi tabel di database) ---

// TenantModel adalah model GORM untuk tabel tenants.
type TenantModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Status    string    `gorm:"not null;default:'Active'"`
	CreatedAt time.Time
}

func (TenantModel) TableName() string { return "tenants" }

// PropertyModel adalah model GORM untuk tabel properties.
type PropertyModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null;uniqueIndex"`
	Description string    `gorm:"type:text"`
}

func (PropertyModel) TableName() string { return "properties" }

// TenantPropertyModel adalah model GORM untuk tabel tenants_properties.
type TenantPropertyModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index"`
	PropertyID uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (TenantPropertyModel) TableName() string { return "tenants_properties" }

// RegulationModel adalah model GORM untuk tabel regulations.
type RegulationModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title          string    `gorm:"not null"`
	RegulationType string    `gorm:"not null"` // POJK, SEOJK, UU
	IssuedDate     time.Time
	Status         string `gorm:"not null;default:'Active'"` // Active, Revoked
	Category       string `gorm:"not null;default:'External'"`
	CreatedAt      time.Time
}

func (RegulationModel) TableName() string { return "regulations" }

// RegulationItemModel adalah model GORM untuk tabel regulation_items.
type RegulationItemModel struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey"`
	RegulationID    uuid.UUID       `gorm:"type:uuid;not null;index"`
	Properties      []PropertyModel `gorm:"many2many:regulation_item_properties;foreignKey:ID;joinForeignKey:regulation_item_id;References:ID;joinReferences:property_id"`
	ItemCode        int             `gorm:"type:int"` // e.g., 1
	ReferenceNumber string          `gorm:"not null"` // e.g., 'Pasal 1 ayat 1'
	Content         string          `gorm:"type:text"`
}

func (RegulationItemModel) TableName() string { return "regulation_items" }

// RegulationPropertyMappingModel adalah model GORM untuk tabel regulation_properties_mapping.
type RegulationPropertyMappingModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RegulationID uuid.UUID `gorm:"type:uuid;not null;index"`
	PropertyID   uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (RegulationPropertyMappingModel) TableName() string { return "regulation_properties_mapping" }

// AssessmentSessionModel adalah model GORM untuk tabel assessment_sessions.
type AssessmentSessionModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Title      string    `gorm:"not null"`
	PeriodYear int       `gorm:"not null"`
	Status     string    `gorm:"not null;default:'Draft'"` // Draft, In_Progress, Completed
	CreatedAt  time.Time
}

func (AssessmentSessionModel) TableName() string { return "assessment_sessions" }

// AssessmentResultModel adalah model GORM untuk tabel assessment_results.
type AssessmentResultModel struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionID        uuid.UUID `gorm:"type:uuid;not null;index"`
	RegulationItemID uuid.UUID `gorm:"type:uuid;not null;index"`
	ComplianceStatus string    `gorm:"not null"` // YES, NO, N/A
	EvidenceLink     string
	Remarks          string `gorm:"type:text"`
	UpdatedAt        time.Time
}

func (AssessmentResultModel) TableName() string { return "assessment_results" }

// RegulationAssessmentModel adalah model GORM untuk tabel regulation_assesments.
type RegulationAssessmentModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RegulationID uuid.UUID `gorm:"type:uuid;not null;index"`
	SessionID    uuid.UUID `gorm:"type:uuid;not null;index"`
	AmountPass   int       `gorm:"default:0"`
	AmountFail   int       `gorm:"default:0"`
	AmountNA     int       `gorm:"default:0"`
}

// UserModel adalah model GORM untuk tabel users.
type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"not null;uniqueIndex"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"not null;uniqueIndex"`
	FullName  string    `gorm:"not null"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Role      string    `gorm:"not null;default:'User'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string { return "users" }

// TenantRegulationModel adalah model GORM untuk tabel tenant_regulations.
type TenantRegulationModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index"`
	RegulationID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time
}

func (TenantRegulationModel) TableName() string { return "tenant_regulations" }

// RiskCategoryModel adalah model GORM untuk tabel risk_categories.
type RiskCategoryModel struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title string    `gorm:"not null"`
}

func (RiskCategoryModel) TableName() string { return "risk_categories" }

// RiskCategoryTenantModel adalah model GORM untuk tabel risk_category_tenants.
type RiskCategoryTenantModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RiskID         *uuid.UUID `gorm:"type:uuid;index"`
	RiskCategoryID uuid.UUID  `gorm:"type:uuid;not null;index"`
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	Appetite       string     `gorm:"type:text"`
	Tolerance      string     `gorm:"type:text"`
}

func (RiskCategoryTenantModel) TableName() string { return "risk_category_tenants" }

// RiskModel adalah model GORM untuk tabel risks.
type RiskModel struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	RiskTitle          string    `gorm:"not null"`
	RiskDescription    string    `gorm:"type:text"`
	CategoryID         uuid.UUID `gorm:"type:uuid;not null;index"`
	LikelihoodInherent int       `gorm:"default:0"`
	ImpactInherent     int       `gorm:"default:0"`
	LikelihoodResidual int       `gorm:"default:0"`
	ImpactResidual     int       `gorm:"default:0"`
	MitigationPlan     string    `gorm:"type:text"`
	MitigationStatus   string    `gorm:"not null;default:'belum direncanakan'"`
}

func (RiskModel) TableName() string { return "risks" }
