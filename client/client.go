package client

import (
	"crypto/tls"
	"github.com/davecgh/go-spew/spew"
	"github.com/jagregory/halgo"
	// "io/ioutil"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Client struct {
	Navigator halgo.Navigator
	Logger    lager.Logger
}

type Application struct {
	halgo.Links
	ApplicationId string `json:"applicationId"`
}

type Schedule struct {
	halgo.Links
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

type Chaos struct {
	halgo.Links
	ApplicationUrl Application
	ScheduleUrl    Schedule
	Probobility    float64
}

type Event struct {
	halgo.Links
	ExecutedAt              string
	TerminatedInstances     []int
	TotalInstanceCount      int
	TerminatedInstanceCount int
	Chaos                   Chaos
}

func New(host string, logger lager.Logger) Client {
	navigator := halgo.NewNavigator(host)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	navigator.HttpClient = &http.Client{Transport: tr}
	client := Client{Navigator: navigator, logger}
	return client
}

func (c Client) CreateApp(applicationId string) (string, error) {
	obj := Application{ApplicationId: applicationId}
	c.Logger.Info("Creating an app: " + json.Marshal(obj))
	return c.create("applications", obj)
}

func (c Client) CreateSchedule(name string, expression string) (string, error) {
	obj := Schedule{Name: name, Expression: expression}
	c.Logger.Info("Creating a schedule: " + json.Marshal(obj))
	return c.create("schedules", obj)
}

func (c Client) CreateChaos(applicationLink string, scheduleLink string, probability float64) (string, error) {
	obj := Chaos{
		ApplicationUrl: applicationLink,
		ScheduleUrl:    scheduleLink,
		Probobility:    probability,
	}
	c.Logger.Info("Creating a chaos: " + json.Marshal(obj))
	return c.create("chaoses", obj)
}

func (c Client) create(resource string, obj interface{}) (string, error) {
	jsonString, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	response, err := c.Navigator.Follow(resource).Post("application/json", bytes.NewReader(jsonString))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", errors.New("Server error")
	}
	resounceUrl, err := response.Location()
	if err != nil {
		return "", err
	}
	return resounceUrl.String(), err
}

func (c Client) Delete(url string) error {
	response, err := c.Navigator.Follow(url).Delete()
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return errors.New("Server error")
	}
	return err
}
