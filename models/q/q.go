package q

import (
	"github.com/yunling101/ControllerManager/models"
)

type M map[string]interface{}

type modal struct {
	tableName string
}

func Table(tableName string) *modal {
	return &modal{tableName: tableName}
}

func (m *modal) WhereMany(result, query interface{}, args ...interface{}) error {
	return models.DB.Table(m.tableName).Where(query, args...).Find(result).Error
}

func (m *modal) WhereOne(result, query interface{}, args ...interface{}) error {
	return models.DB.Table(m.tableName).Where(query, args...).First(result).Error
}

func (m *modal) QueryMany(selector map[string]interface{}, result interface{}) error {
	return models.DB.Table(m.tableName).Where(selector).Find(result).Error
}

func (m *modal) QueryOne(selector map[string]interface{}, result interface{}) error {
	return models.DB.Table(m.tableName).Where(selector).First(result).Error
}

func (m *modal) InsertOne(value interface{}) error {
	return models.DB.Table(m.tableName).Create(value).Error
}

func (m *modal) Count(selector map[string]interface{}) (count int64, err error) {
	err = models.DB.Table(m.tableName).Where(selector).Count(&count).Error
	return
}

func (m *modal) UpdateOne(selector, update map[string]interface{}) error {
	return models.DB.Table(m.tableName).Where(selector).Updates(update).Error
}
