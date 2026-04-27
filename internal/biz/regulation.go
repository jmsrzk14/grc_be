package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// RegulationUseCase menangani logika bisnis untuk Regulation.
type RegulationUseCase struct {
	repo        RegulationRepo
	itemRepo    RegulationItemRepo
	mappingRepo RegulationPropertyMappingRepo
	log         *log.Helper
}

// NewRegulationUseCase membuat instance baru.
func NewRegulationUseCase(
	repo RegulationRepo,
	itemRepo RegulationItemRepo,
	mappingRepo RegulationPropertyMappingRepo,
	logger log.Logger,
) *RegulationUseCase {
	return &RegulationUseCase{
		repo:        repo,
		itemRepo:    itemRepo,
		mappingRepo: mappingRepo,
		log:         log.NewHelper(logger),
	}
}

// CreateRegulation membuat regulasi baru.
func (uc *RegulationUseCase) CreateRegulation(ctx context.Context, r *Regulation) (*Regulation, error) {
	if r.Title == "" {
		return nil, fmt.Errorf("regulation title is required")
	}
	r.ID = uuid.New()
	return uc.repo.Create(ctx, r)
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

// --- RegulationItem Use Cases ---

// CreateRegulationItem membuat item/pasal baru dalam regulasi.
func (uc *RegulationUseCase) CreateRegulationItem(ctx context.Context, item *RegulationItem) (*RegulationItem, error) {
	if item.RegulationID == uuid.Nil {
		return nil, fmt.Errorf("regulation_id is required")
	}
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
