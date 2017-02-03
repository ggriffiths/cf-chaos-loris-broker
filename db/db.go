package db

import (
	"fmt"
	"github.com/Altoros/cf-chaos-loris-broker/model"
	"github.com/jinzhu/gorm"
)

func New() (*gorm.DB, error) {

	creds, err := LoadServiceCredentials("p-mysql")

	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			creds.GetUsername(),
			creds.GetPassword(),
			creds.GetHost(),
			creds.GetPort(),
			creds.GetDBName()))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.ServiceInstance{}, &model.ServiceBinding{}).Error

	return db, err
}
