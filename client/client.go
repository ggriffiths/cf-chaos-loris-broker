package client

import (
	"code.cloudfoundry.org/lager"
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
	Navigator  *halgo.Navigator
	Host       string
	HttpClient *http.Client
	Logger     lager.Logger
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
	ApplicationUrl string  `json:"application"`
	ScheduleUrl    string  `json:"schedule"`
	Probobility    float64 `json:"probability"`
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
	// navigator := halgo.NewNavigator(host)
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// navigator.HttpClient = &http.Client{Transport: tr}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	client := Client{Logger: logger, Host: host, HttpClient: httpClient}
	return client
}

func (c Client) navigator() halgo.Navigator {
	navigator := halgo.NewNavigator(c.Host)

	navigator.HttpClient = c.HttpClient

	return navigator
}

func (c Client) CreateApp(applicationId string) (string, error) {
	obj := Application{ApplicationId: applicationId}
	c.Logger.Info("Creating an app")
	return c.create("applications", obj)
}

func (c Client) CreateSchedule(name string, expression string) (string, error) {
	obj := Schedule{Name: name, Expression: expression}
	c.Logger.Info("Creating a schedule")
	return c.create("schedules", obj)
}

func (c Client) CreateChaos(applicationLink string, scheduleLink string, probability float64) (string, error) {
	obj := Chaos{
		ApplicationUrl: applicationLink,
		ScheduleUrl:    scheduleLink,
		Probobility:    probability,
	}
	jsonString, err := json.Marshal(obj)
	c.Logger.Info("Creating a resource: " + string(jsonString))
	if err != nil {
		return "", err
	}

	spew.Dump(c.HttpClient)
	response, err := c.HttpClient.Post(c.Host+"/chaoses", "application/json", bytes.NewReader(jsonString))
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	resounceUrl, err := response.Location()
	if err != nil {
		return "", err
	}

	c.Logger.Info("Chaos is Created")
	return resounceUrl.String(), err
}

func (c Client) create(resource string, obj interface{}) (string, error) {
	jsonString, err := json.Marshal(obj)
	c.Logger.Info("Creating a resource: " + string(jsonString))
	if err != nil {
		return "", err
	}
	response, err := c.navigator().Follow(resource).Post("application/json", bytes.NewReader(jsonString))
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
	response, err := c.navigator().Follow(url).Delete()
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return errors.New("Server error")
	}
	return err
}
