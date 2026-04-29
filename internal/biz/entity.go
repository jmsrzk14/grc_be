package biz

import (
	// "context"
	"time"

	"github.com/google/uuid"
)

// --- Tenant & Organization Domain ---

// Tenant merepresentasikan organisasi (BPR/BPRS) yang menggunakan sistem.
type Tenant struct {
	ID        uuid.UUID
	Name      string
	Type      string // e.g., 'BPR', 'BPRS'
	Status    string // e.g., 'Active', 'Inactive'
	CreatedAt time.Time
}

// Property merepresentasikan kategori aset (e.g., 'Teknologi', 'SDM', 'Fisik').
type Property struct {
	ID          uuid.UUID
	Name        string
	Description string
}

// TenantProperty adalah mapping many-to-many antara Tenant dan Property.
type TenantProperty struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	PropertyID uuid.UUID
}

// --- Regulation Domain ---


// Regulation merepresentasikan dokumen regulasi (POJK, SEOJK, UU).
type Regulation struct {
	ID             uuid.UUID
	Title          string
	RegulationType string // POJK, SEOJK, UU
	IssuedDate     time.Time
	Status         string // Active, Revoked
	Category       string // Internal, External
	AmountPass     int
	AmountFail     int
	AmountNA       int
}

// RegulationItem merepresentasikan item/pasal dalam suatu regulasi.
type RegulationItem struct {
	ID              uuid.UUID
	RegulationID    uuid.UUID
	PropertyIDs     []uuid.UUID
	ReferenceNumber string // e.g., 'Pasal 1 ayat 1'
	Content         string
}

// RegulationPropertyMapping adalah mapping antara regulasi dan kategori aset.
type RegulationPropertyMapping struct {
	ID           uuid.UUID
	RegulationID uuid.UUID
	PropertyID   uuid.UUID
}

// --- Assessment Domain ---

// AssessmentSession merepresentasikan satu sesi penilaian compliance.
type AssessmentSession struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	Title      string
	PeriodYear int
	Status     string // Draft, In_Progress, Completed
	CreatedAt  time.Time
}

// AssessmentResult merepresentasikan hasil penilaian per item regulasi.
type AssessmentResult struct {
	ID               uuid.UUID
	SessionID        uuid.UUID
	RegulationItemID uuid.UUID
	ComplianceStatus string // YES, NO, N/A
	EvidenceLink     string
	Remarks          string
	UpdatedAt        time.Time
}

// RegulationAssessment merupakan ringkasan hasil assessment per regulasi dalam satu sesi.
type RegulationAssessment struct {
	ID           uuid.UUID
	RegulationID uuid.UUID
	SessionID    uuid.UUID
	AmountPass   int
	AmountFail   int
	AmountNA     int
}

// User merepresentasikan pengguna sistem.
type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Email     string
	FullName  string
	TenantID  uuid.UUID
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// --- Risk Management Domain ---

// RiskCategory merepresentasikan kategori risiko (e.g., 'Operasional', 'Finansial').
type RiskCategory struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Appetite  string    `json:"appetite"`
	Tolerance string    `json:"tolerance"`
}

// Risk merepresentasikan item risiko.
type Risk struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	RiskTitle          string    `json:"risk_title"`
	RiskDescription    string    `json:"risk_description"`
	CategoryID         uuid.UUID `json:"category_id"`
	LikelihoodInherent int       `json:"likelihood_inherent"`
	ImpactInherent     int       `json:"impact_inherent"`
	LikelihoodResidual int       `json:"likelihood_residual"`
	ImpactResidual     int       `json:"impact_residual"`
	MitigationPlan     string    `json:"mitigation_plan"`
	MitigationStatus   string    `json:"mitigation_status"`
}
