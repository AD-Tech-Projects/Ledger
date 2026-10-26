package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository[M any, E any] struct {
	db         *gorm.DB
	toDomain   func(*M) *E
	fromDomain func(*E) *M
}

type QueryOption[M any] func(*gorm.DB) *gorm.DB

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

func (r *GormRepository[M, E]) WithDB(db *gorm.DB) *GormRepository[M, E] {
	return &GormRepository[M, E]{
		db:         db,
		toDomain:   r.toDomain,
		fromDomain: r.fromDomain,
	}
}

func (r *GormRepository[M, E]) Create(ctx context.Context, domainEntity *E) error {
	model := r.fromDomain(domainEntity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormRepository[M, E]) GetByID(
	ctx context.Context,
	id string,
	options ...QueryOption[M],
) (*E, error) {
	var model M
	query := r.db.WithContext(ctx)
	for _, option := range options {
		query = option(query)
	}

	err := query.Where(clause.Eq{
		Column: clause.Column{
			Table: clause.CurrentTable,
			Name:  "id",
		},
		Value: id,
	}).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *GormRepository[M, E]) GetAll(
	ctx context.Context,
	filter *E,
	options ...QueryOption[M],
) ([]*E, error) {
	var models []M
	query := r.db.WithContext(ctx)
	for _, option := range options {
		query = option(query)
	}
	if filter != nil {
		query = query.Where(r.fromDomain(filter))
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*E, 0, len(models))
	for index := range models {
		entities = append(entities, r.toDomain(&models[index]))
	}

	return entities, nil
}

func (r *GormRepository[M, E]) Update(ctx context.Context, domainEntity *E) error {
	model := r.fromDomain(domainEntity)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *GormRepository[M, E]) Delete(ctx context.Context, id string) error {
	var model M
	return r.db.WithContext(ctx).Delete(&model, id).Error
}

func (r *GormRepository[M, E]) Transaction(
	ctx context.Context,
	fn func(*gorm.DB) error,
) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func Join[M any](association string) QueryOption[M] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(association)
	}
}
