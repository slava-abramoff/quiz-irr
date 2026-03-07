package storage

import "gorm.io/gorm"

type rawRepo struct {
	db *gorm.DB
}

func NewRawRepo(db *gorm.DB) *rawRepo {
	return &rawRepo{db: db}
}

func (r *rawRepo) Create()      {}
func (r *rawRepo) SavePayload() {}
func (r *rawRepo) Load()        {}
func (r *rawRepo) Delete()      {}
