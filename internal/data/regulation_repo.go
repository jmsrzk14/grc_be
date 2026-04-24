package data

import (
	"context"
	"errors"

	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Regulation Repository Implementation ---

type regulationRepo struct {
	data *Data
	log  *log.Helper
}

// NewRegulationRepo membuat instance repository Regulation.
func NewRegulationRepo(data *Data, logger log.Logger) biz.RegulationRepo {
	return &regulationRepo{data: data, log: log.NewHelper(logger)}
}

func (r *regulationRepo) Create(ctx context.Context, reg *biz.Regulation) (*biz.Regulation, error) {
	m := &RegulationModel{
		ID:             reg.ID,
		Title:          reg.Title,
		RegulationType: reg.RegulationType,
		IssuedDate:     reg.IssuedDate,
		Status:         reg.Status,
		Category:       reg.Category,
	}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	return reg, nil
}

func (r *regulationRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Regulation, error) {
	type resultWithStats struct {
		RegulationModel
		AmountPass int
		AmountFail int
		AmountNA   int
	}
	var res resultWithStats
	err := r.data.db.WithContext(ctx).Table("regulations").
		Select("regulations.*, COALESCE(SUM(ra.amount_pass), 0) as amount_pass, COALESCE(SUM(ra.amount_fail), 0) as amount_fail, COALESCE(SUM(ra.amount_na), 0) as amount_na").
		Joins("LEFT JOIN regulation_assesments ra ON regulations.id = ra.regulation_id").
		Where("regulations.id = ?", id).
		Group("regulations.id").
		Scan(&res).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, err
	}
	// GORM Scan might not return ErrRecordNotFound with Group/Select
	if res.ID == uuid.Nil {
		return nil, biz.ErrNotFound
	}

	reg := toRegulationDomain(&res.RegulationModel)
	reg.AmountPass = res.AmountPass
	reg.AmountFail = res.AmountFail
	reg.AmountNA = res.AmountNA
	return reg, nil
}

func (r *regulationRepo) FindAll(ctx context.Context) ([]*biz.Regulation, error) {
	type resultWithStats struct {
		RegulationModel
		AmountPass int
		AmountFail int
		AmountNA   int
	}
	var results []*resultWithStats

	err := r.data.db.WithContext(ctx).Table("regulations").
		Select("regulations.*, COALESCE(SUM(ra.amount_pass), 0) as amount_pass, COALESCE(SUM(ra.amount_fail), 0) as amount_fail, COALESCE(SUM(ra.amount_na), 0) as amount_na").
		Joins("LEFT JOIN regulation_assesments ra ON regulations.id = ra.regulation_id").
		Group("regulations.id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	regs := make([]*biz.Regulation, 0, len(results))
	for _, res := range results {
		reg := toRegulationDomain(&res.RegulationModel)
		reg.AmountPass = res.AmountPass
		reg.AmountFail = res.AmountFail
		reg.AmountNA = res.AmountNA
		regs = append(regs, reg)
	}
	return regs, nil
}

func (r *regulationRepo) Update(ctx context.Context, reg *biz.Regulation) (*biz.Regulation, error) {
	m := &RegulationModel{
		ID:             reg.ID,
		Title:          reg.Title,
		RegulationType: reg.RegulationType,
		IssuedDate:     reg.IssuedDate,
		Status:         reg.Status,
		Category:       reg.Category,
	}
	if result := r.data.db.WithContext(ctx).Save(m); result.Error != nil {
		return nil, result.Error
	}
	return toRegulationDomain(m), nil
}

func (r *regulationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&RegulationModel{}, "id = ?", id).Error
}

func toRegulationDomain(m *RegulationModel) *biz.Regulation {
	return &biz.Regulation{
		ID:             m.ID,
		Title:          m.Title,
		RegulationType: m.RegulationType,
		IssuedDate:     m.IssuedDate,
		Status:         m.Status,
		Category:       m.Category,
	}
}

// --- RegulationItem Repository Implementation ---

type regulationItemRepo struct {
	data *Data
	log  *log.Helper
}

// NewRegulationItemRepo membuat instance repository RegulationItem.
func NewRegulationItemRepo(data *Data, logger log.Logger) biz.RegulationItemRepo {
	return &regulationItemRepo{data: data, log: log.NewHelper(logger)}
}

func (r *regulationItemRepo) Create(ctx context.Context, item *biz.RegulationItem) (*biz.RegulationItem, error) {
	m := &RegulationItemModel{
		ID:              item.ID,
		RegulationID:    item.RegulationID,
		ReferenceNumber: item.ReferenceNumber,
		Content:         item.Content,
	}

	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		if len(item.TenantPropertyIDs) > 0 {
			tenantProps := make([]TenantPropertyModel, 0, len(item.TenantPropertyIDs))
			for _, id := range item.TenantPropertyIDs {
				tenantProps = append(tenantProps, TenantPropertyModel{ID: id})
			}
			return tx.Model(m).Association("TenantProperties").Append(tenantProps)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *regulationItemRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.RegulationItem, error) {
	var m RegulationItemModel
	if result := r.data.db.WithContext(ctx).Preload("TenantProperties").First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toRegulationItemDomain(&m), nil
}

func (r *regulationItemRepo) FindByRegulationID(ctx context.Context, regulationID uuid.UUID) ([]*biz.RegulationItem, error) {
	var models []*RegulationItemModel
	if result := r.data.db.WithContext(ctx).Preload("TenantProperties").Order("reference_number ASC").Find(&models, "regulation_id = ?", regulationID); result.Error != nil {
		return nil, result.Error
	}
	items := make([]*biz.RegulationItem, 0, len(models))
	for _, m := range models {
		items = append(items, toRegulationItemDomain(m))
	}
	return items, nil
}

func (r *regulationItemRepo) Update(ctx context.Context, item *biz.RegulationItem) (*biz.RegulationItem, error) {
	tenantProps := make([]TenantPropertyModel, 0, len(item.TenantPropertyIDs))
	for _, id := range item.TenantPropertyIDs {
		tenantProps = append(tenantProps, TenantPropertyModel{ID: id})
	}
	m := &RegulationItemModel{
		ID:               item.ID,
		RegulationID:     item.RegulationID,
		ReferenceNumber:  item.ReferenceNumber,
		Content:          item.Content,
	}

	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		return tx.Model(m).Association("TenantProperties").Replace(tenantProps)
	})

	if err != nil {
		return nil, err
	}
	return toRegulationItemDomain(m), nil
}

func (r *regulationItemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&RegulationItemModel{}, "id = ?", id).Error
}

func toRegulationItemDomain(m *RegulationItemModel) *biz.RegulationItem {
	ids := make([]uuid.UUID, 0, len(m.TenantProperties))
	for _, tp := range m.TenantProperties {
		ids = append(ids, tp.ID)
	}
	return &biz.RegulationItem{
		ID:                m.ID,
		RegulationID:      m.RegulationID,
		TenantPropertyIDs: ids,
		ReferenceNumber:   m.ReferenceNumber,
		Content:           m.Content,
	}
}

// --- RegulationPropertyMapping Repository Implementation ---

type regulationPropertyMappingRepo struct {
	data *Data
	log  *log.Helper
}

// NewRegulationPropertyMappingRepo membuat instance repository RegulationPropertyMapping.
func NewRegulationPropertyMappingRepo(data *Data, logger log.Logger) biz.RegulationPropertyMappingRepo {
	return &regulationPropertyMappingRepo{data: data, log: log.NewHelper(logger)}
}

func (r *regulationPropertyMappingRepo) Create(ctx context.Context, mapping *biz.RegulationPropertyMapping) (*biz.RegulationPropertyMapping, error) {
	m := &RegulationPropertyMappingModel{
		ID:           mapping.ID,
		RegulationID: mapping.RegulationID,
		PropertyID:   mapping.PropertyID,
	}
	if result := r.data.db.WithContext(ctx).Create(m); result.Error != nil {
		return nil, result.Error
	}
	return mapping, nil
}

func (r *regulationPropertyMappingRepo) FindByRegulationID(ctx context.Context, regulationID uuid.UUID) ([]*biz.RegulationPropertyMapping, error) {
	var models []*RegulationPropertyMappingModel
	if result := r.data.db.WithContext(ctx).Find(&models, "regulation_id = ?", regulationID); result.Error != nil {
		return nil, result.Error
	}
	mappings := make([]*biz.RegulationPropertyMapping, 0, len(models))
	for _, m := range models {
		mappings = append(mappings, &biz.RegulationPropertyMapping{ID: m.ID, RegulationID: m.RegulationID, PropertyID: m.PropertyID})
	}
	return mappings, nil
}

func (r *regulationPropertyMappingRepo) FindByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*biz.RegulationPropertyMapping, error) {
	var models []*RegulationPropertyMappingModel
	if result := r.data.db.WithContext(ctx).Find(&models, "property_id = ?", propertyID); result.Error != nil {
		return nil, result.Error
	}
	mappings := make([]*biz.RegulationPropertyMapping, 0, len(models))
	for _, m := range models {
		mappings = append(mappings, &biz.RegulationPropertyMapping{ID: m.ID, RegulationID: m.RegulationID, PropertyID: m.PropertyID})
	}
	return mappings, nil
}

func (r *regulationPropertyMappingRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&RegulationPropertyMappingModel{}, "id = ?", id).Error
}
