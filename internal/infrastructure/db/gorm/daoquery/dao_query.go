package daoquery

import (
	"auth/internal/infrastructure/db/gorm/dao"

	"gorm.io/gorm"
)

func NewDaoQuery(db *gorm.DB) *dao.Query {
	return dao.Use(db)
}
