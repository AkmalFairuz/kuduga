package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
)

type MetaRepository struct {
	db *database.DB
}

func NewMetaRepository(db *database.DB) *MetaRepository {
	return &MetaRepository{db}
}

func (r *MetaRepository) Get(name string) (string, error) {
	var ret model.Meta
	if err := r.db.Get(&ret, "SELECT * FROM meta WHERE name = ?", name); err != nil {
		return "", err
	}
	return ret.Value, nil
}

func (r *MetaRepository) GetAll() ([]model.Meta, error) {
	var ret []model.Meta
	if err := r.db.Select(&ret, "SELECT * FROM meta"); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *MetaRepository) Set(name string, value string) error {
	_, err := r.db.Exec("REPLACE INTO meta (name, value) VALUES (?, ?)", name, value)
	return err
}
