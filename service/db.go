package service

import (
	"time"

	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/module/db"
	"gorm.io/gorm"
)

type DBService struct{}

func (d *DBService) DB() *gorm.DB {
	return db.DB
}

func (*DBService) Insert(data model.Table) error {
	return db.DB.Create(data).Error
}

func (*DBService) First(data model.Table) error {
	return db.DB.Where(data).First(data).Error
}

func (*DBService) Count(data model.Table, wheres ...string) (count int64, err error) {
	m := db.DB.Model(data).Where(data)
	for _, where := range wheres {
		m = m.Where(where)
	}
	err = m.Count(&count).Error
	return
}

func (*DBService) Find(data model.Table, list any, wheres []string, args ...int) (err error) {
	var offset, limit int
	if len(args) == 1 {
		limit = args[0]
	} else if len(args) >= 2 {
		offset = args[0]
		limit = args[1]
	}
	var startedAt, ended_at int
	if len(args) >= 3 {
		startedAt = args[2]
	}
	if len(args) == 4 {
		ended_at = args[3]
	}

	m := db.DB.Where(data)
	for _, where := range wheres {
		m = m.Where(where)
	}
	if startedAt > 0 {
		m = m.Where("created_at >= ?", time.Unix(int64(startedAt), 0))
	}
	if ended_at > 0 {
		m = m.Where("created_at < ?", time.Unix(int64(ended_at), 0))
	}
	m = m.Order("created_at desc")
	if offset != 0 {
		m = m.Offset(offset)
	}
	if limit > 0 {
		m = m.Limit(limit)
	}
	return m.Find(list).Error
}

func (*DBService) FindWhere(data model.Table, wheres []string, list any, args ...int) (err error) {
	var offset, limit int
	if len(args) == 1 {
		limit = args[0]
	} else if len(args) == 2 {
		offset = args[0]
		limit = args[1]
	}
	m := db.DB.Where(data)
	for _, where := range wheres {
		m = m.Where(where)
	}
	if offset != 0 {
		m = m.Offset(offset)
	}
	if limit > 0 {
		m = m.Limit(limit)
	}
	return m.Find(&list).Error
}

func (*DBService) Updates(data model.Table, cols ...string) error {
	return db.DB.Model(data).Select(cols).Updates(data).Error
}

func (*DBService) Delete(data model.Table, wheres ...string) error {
	m := db.DB.Model(data)
	for _, where := range wheres {
		m = m.Where(where)
	}
	return m.Delete(data).Error
}

func (*DBService) RawSQL(data any, sql string, args ...any) error {
	return db.DB.Raw(sql, args...).Scan(data).Error
}

type DB = gorm.DB

func (*DBService) Transaction(fc func(tx *DB) error) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		return fc(tx)
	})
}

func (*DBService) User(uid int, checkKey ...bool) (*model.User, error) {
	var user model.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		return nil, err
	}
	if len(checkKey) > 0 && checkKey[0] {
		if user.APIKey == "" {
			return nil, errNotFoundKey
		}
	}
	return &user, nil
}
