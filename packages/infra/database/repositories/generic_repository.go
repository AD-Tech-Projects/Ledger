package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type GormRepository[M any, E any] struct {
	db         *gorm.DB
	toDomain   func(*M) *E
	fromDomain func(*E) *M
}

func NewRepository[M any, E any](
	db *gorm.DB,
	toDomain func(*M) *E,
	fromDomain func(*E) *M,
) *GormRepository[M, E] {
	return &GormRepository[M, E]{
		db:         db,
		toDomain:   toDomain,
		fromDomain: fromDomain,
	}
}

func (r *GormRepository[M, E]) Create(ctx context.Context, domainEntity *E) error {
	model := r.fromDomain(domainEntity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*domainEntity = *r.toDomain(model)
	return nil
}

func (r *GormRepository[M, E]) GetByID(ctx context.Context, id uint) (*E, error) {
	var model M
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *GormRepository[M, E]) Update(ctx context.Context, domainEntity *E) error {
	model := r.fromDomain(domainEntity)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *GormRepository[M, E]) Delete(ctx context.Context, id uint) error {
	var model M
	return r.db.WithContext(ctx).Delete(&model, id).Error
}
