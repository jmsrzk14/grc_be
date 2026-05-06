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

func (r *regulationRepo) FindByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*biz.Regulation, error) {
	type resultWithStats struct {
		RegulationModel
		AmountPass int
		AmountFail int
		AmountNA   int
	}
	var res resultWithStats

	db := r.data.db.WithContext(ctx).Table("regulations")

	if tenantID != uuid.Nil {
		// Ambil session terbaru milik tenant
		type sessionRow struct{ ID uuid.UUID }
		var latestSession sessionRow
		r.data.db.WithContext(ctx).
			Table("assessment_sessions").
			Select("id").
			Where("tenant_id = ?", tenantID).
			Order("created_at DESC").
			Limit(1).
			Scan(&latestSession)

		if latestSession.ID != uuid.Nil {
			// Hitung amount dari tabel assessment_results berdasarkan session terbaru
			subQuery := r.data.db.Table("assessment_results").
				Select(
					"regulation_items.regulation_id,"+
						" COUNT(CASE WHEN assessment_results.compliance_status = 'YES' THEN 1 END) AS amount_pass,"+
						" COUNT(CASE WHEN assessment_results.compliance_status = 'NO' THEN 1 END) AS amount_fail,"+
						" COUNT(CASE WHEN assessment_results.compliance_status = 'N/A' THEN 1 END) AS amount_na",
				).
				Joins("JOIN regulation_items ON regulation_items.id = assessment_results.regulation_item_id").
				Where("assessment_results.session_id = ?", latestSession.ID).
				Group("regulation_items.regulation_id")

			db = db.
				Select("regulations.*, COALESCE(stats.amount_pass, 0) AS amount_pass, COALESCE(stats.amount_fail, 0) AS amount_fail, COALESCE(stats.amount_na, 0) AS amount_na").
				Joins("LEFT JOIN (?) AS stats ON regulations.id = stats.regulation_id", subQuery)
		} else {
			db = db.Select("regulations.*, 0 AS amount_pass, 0 AS amount_fail, 0 AS amount_na")
		}
	} else {
		db = db.Select("regulations.*, 0 AS amount_pass, 0 AS amount_fail, 0 AS amount_na")
	}

	if err := db.Where("regulations.id = ?", id).Scan(&res).Error; err != nil {
		return nil, err
	}
	if res.ID == uuid.Nil {
		return nil, biz.ErrNotFound
	}

	reg := toRegulationDomain(&res.RegulationModel)
	reg.AmountPass = res.AmountPass
	reg.AmountFail = res.AmountFail
	reg.AmountNA = res.AmountNA
	return reg, nil
}

func (r *regulationRepo) FindAll(ctx context.Context, tenantID uuid.UUID) ([]*biz.Regulation, error) {
	type resultWithStats struct {
		RegulationModel
		AmountPass int
		AmountFail int
		AmountNA   int
	}
	var results []*resultWithStats

	db := r.data.db.WithContext(ctx).Table("regulations")

	if tenantID != uuid.Nil {
		// Ambil session terbaru milik tenant
		type sessionRow struct{ ID uuid.UUID }
		var latestSession sessionRow
		r.data.db.WithContext(ctx).
			Table("assessment_sessions").
			Select("id").
			Where("tenant_id = ?", tenantID).
			Order("created_at DESC").
			Limit(1).
			Scan(&latestSession)

		if latestSession.ID != uuid.Nil {
			// Hitung amount dari tabel assessment_results berdasarkan session terbaru
			subQuery := r.data.db.Table("assessment_results").
				Select(
					"regulation_items.regulation_id,"+
						" COUNT(CASE WHEN assessment_results.compliance_status = 'YES' THEN 1 END) AS amount_pass,"+
						" COUNT(CASE WHEN assessment_results.compliance_status = 'NO' THEN 1 END) AS amount_fail,"+
						" COUNT(CASE WHEN assessment_results.compliance_status = 'N/A' THEN 1 END) AS amount_na",
				).
				Joins("JOIN regulation_items ON regulation_items.id = assessment_results.regulation_item_id").
				Where("assessment_results.session_id = ?", latestSession.ID).
				Group("regulation_items.regulation_id")

			db = db.
				Select("regulations.*, COALESCE(stats.amount_pass, 0) AS amount_pass, COALESCE(stats.amount_fail, 0) AS amount_fail, COALESCE(stats.amount_na, 0) AS amount_na").
				Joins("LEFT JOIN (?) AS stats ON regulations.id = stats.regulation_id", subQuery)
		} else {
			db = db.Select("regulations.*, 0 AS amount_pass, 0 AS amount_fail, 0 AS amount_na")
		}
	} else {
		db = db.Select("regulations.*, 0 AS amount_pass, 0 AS amount_fail, 0 AS amount_na")
	}

	if err := db.Scan(&results).Error; err != nil {
		r.log.Errorf("failed to list regulations: %v", err)
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

func (r *regulationRepo) FindByTitle(ctx context.Context, title string) (*biz.Regulation, error) {
	var m RegulationModel
	if err := r.data.db.WithContext(ctx).Where("title = ?", title).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, err
	}
	return toRegulationDomain(&m), nil
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
		ItemCode:        item.ItemCode,
		ReferenceNumber: item.ReferenceNumber,
		Content:         item.Content,
	}

	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		if len(item.PropertyIDs) > 0 {
			props := make([]PropertyModel, 0, len(item.PropertyIDs))
			for _, id := range item.PropertyIDs {
				props = append(props, PropertyModel{ID: id})
			}
			return tx.Model(m).Association("Properties").Append(props)
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
	if result := r.data.db.WithContext(ctx).Preload("Properties").First(&m, "id = ?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, result.Error
	}
	return toRegulationItemDomain(&m), nil
}

func (r *regulationItemRepo) FindByRegulationID(ctx context.Context, regulationID uuid.UUID, tenantID uuid.UUID) ([]*biz.RegulationItem, error) {
	var models []*RegulationItemModel
	db := r.data.db.WithContext(ctx).Preload("Properties").Order("reference_number ASC")

	if tenantID != uuid.Nil {
		// Include items with NO properties OR items that match the tenant's properties
		db = db.Where("regulation_items.regulation_id = ?", regulationID).
			Where(`regulation_items.id NOT IN (SELECT regulation_item_id FROM regulation_item_properties) OR 
			       regulation_items.id IN (
				       SELECT rip.regulation_item_id 
				       FROM regulation_item_properties rip 
				       JOIN tenants_properties tp ON rip.property_id = tp.property_id 
				       WHERE tp.tenant_id = ?
			       )`, tenantID)
	} else {
		db = db.Where("regulation_id = ?", regulationID)
	}

	if result := db.Find(&models); result.Error != nil {
		return nil, result.Error
	}
	items := make([]*biz.RegulationItem, 0, len(models))
	for _, m := range models {
		items = append(items, toRegulationItemDomain(m))
	}
	return items, nil
}

// FindExcludedByTenantID mengembalikan semua regulation items yang propertinya
// TIDAK ada dalam daftar properti yang dimiliki tenant.
// Digunakan untuk auto-seed N/A saat session baru dibuat.
func (r *regulationItemRepo) FindExcludedByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*biz.RegulationItem, error) {
	var models []*RegulationItemModel

	// Ambil items yang TIDAK memiliki property yang cocok dengan tenant
	// yaitu items yang property-nya tidak ada di tenants_properties untuk tenant ini
	err := r.data.db.WithContext(ctx).
		Preload("Properties").
		Where(`regulation_items.id NOT IN (
			SELECT DISTINCT rip.regulation_item_id
			FROM regulation_item_properties rip
			JOIN tenants_properties tp ON rip.property_id = tp.property_id
			WHERE tp.tenant_id = ?
		)`, tenantID).
		Find(&models).Error

	if err != nil {
		return nil, err
	}
	items := make([]*biz.RegulationItem, 0, len(models))
	for _, m := range models {
		items = append(items, toRegulationItemDomain(m))
	}
	return items, nil
}

func (r *regulationItemRepo) Update(ctx context.Context, item *biz.RegulationItem) (*biz.RegulationItem, error) {
	props := make([]PropertyModel, 0, len(item.PropertyIDs))
	for _, id := range item.PropertyIDs {
		props = append(props, PropertyModel{ID: id})
	}
	m := &RegulationItemModel{
		ID:              item.ID,
		RegulationID:    item.RegulationID,
		ItemCode:        item.ItemCode,
		ReferenceNumber: item.ReferenceNumber,
		Content:         item.Content,
	}

	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		return tx.Model(m).Association("Properties").Replace(props)
	})

	if err != nil {
		return nil, err
	}
	return toRegulationItemDomain(m), nil
}

func (r *regulationItemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.db.WithContext(ctx).Delete(&RegulationItemModel{}, "id = ?", id).Error
}

func (r *regulationItemRepo) FindByRegulationIDAndItemCode(ctx context.Context, regulationID uuid.UUID, itemCode string) (*biz.RegulationItem, error) {
	var m RegulationItemModel
	if err := r.data.db.WithContext(ctx).Preload("Properties").Where("regulation_id = ? AND item_code = ?", regulationID, itemCode).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, biz.ErrNotFound
		}
		return nil, err
	}
	return toRegulationItemDomain(&m), nil
}

func toRegulationItemDomain(m *RegulationItemModel) *biz.RegulationItem {
	ids := make([]uuid.UUID, 0, len(m.Properties))
	for _, p := range m.Properties {
		ids = append(ids, p.ID)
	}
	return &biz.RegulationItem{
		ID:              m.ID,
		RegulationID:    m.RegulationID,
		PropertyIDs:     ids,
		ItemCode:        m.ItemCode,
		ReferenceNumber: m.ReferenceNumber,
		Content:         m.Content,
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
