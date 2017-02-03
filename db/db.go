package db

import (
	"errors"
	"fmt"
	binded_service "github.com/18F/cf-service-connect/models"
	"github.com/Altoros/cf-chaos-loris-broker/model"
	"github.com/jinzhu/gorm"
	"os"
)

func New() (*gorm.DB, error) {
	if os.Getenv("VCAP_SERVICES") == "" {
		return nil, errors.New("no p-mysql services is not added")
	}

	creds, err := binded_service.CredentialsFromJSON(os.Getenv("VCAP_SERVICES"))
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%d?charset=utf8&parseTime=True&loc=Local",
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
