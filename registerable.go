package gormcallback

import "gorm.io/gorm"

type Registerable interface {
	Register(name string, fn func(*gorm.DB)) error
}

func registerCallback(db *gorm.DB, cb func(*gorm.DB)) error {
	registerables := []Registerable{
		db.Callback().Create(),
		db.Callback().Query(),
		db.Callback().Update(),
		db.Callback().Delete(),
		db.Callback().Raw(),
		db.Callback().Row(),
	}
	for _, registerable := range registerables {
		if err := registerable.Register("gormcallback", cb); err != nil {
			return err
		}
	}
	return nil
}
