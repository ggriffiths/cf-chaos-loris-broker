package model

import (
	"github.com/jinzhu/gorm"
)

type ServiceInstance struct {
	gorm.Model
	InstanceId  string
	ScheduleUrl string
	Probability float64
}

type ServiceBinding struct {
	gorm.Model
	BindingId      string
	InstanceId     string
	ApplicationUrl string
	ChaosUrl       string
}
