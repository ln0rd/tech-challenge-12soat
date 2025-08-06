package interfaces

import "gorm.io/gorm"

// Database define a interface para operações de banco de dados
type Database interface {
	Create(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Updates(values interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
}
