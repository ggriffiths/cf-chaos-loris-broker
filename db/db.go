package db

import (
	binded_service "github.com/18F/cf-service-connect/models"
	"os"
)

func New() (gorm.DB, error) {
	creds := binded_service.CredentialsFromJSON(os.Env("VCAP_SERVICES"))
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%d?charset=utf8&parseTime=True&loc=Local",
			creds.GetUsername(),
			creds.GetPassword(),
			creds.GetHost(),
			creds.GetPort(),
			creds.GetDBName()))

	// db.AutoMigrate(&User{}, &Product{}, &Order{})

	return db, err
}
