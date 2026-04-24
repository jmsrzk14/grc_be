package data

import (
	"context"
	"errors"
	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *authRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	var m UserModel
	if err := r.data.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	return &biz.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		Email:     m.Email,
		FullName:  m.FullName,
		TenantID:  m.TenantID,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *authRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	m := UserModel{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		FullName:  u.FullName,
		TenantID:  u.TenantID,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if err := r.data.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return u, nil
}
