package repository

import "github.com/akmalfairuz/finance/module/database"

type AnalyticsRepository struct {
	db *database.DB
}

func NewAnalyticsRepository(db *database.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// TODO: implement
