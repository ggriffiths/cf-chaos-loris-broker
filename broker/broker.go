package broker

import (
	"code.cloudfoundry.org/lager"
	"context"
	"errors"
	"github.com/Altoros/cf-chaos-loris-broker/client"
	"github.com/Altoros/cf-chaos-loris-broker/cmd"
	"github.com/Altoros/cf-chaos-loris-broker/config"
	"github.com/Altoros/cf-chaos-loris-broker/model"
	"github.com/jinzhu/gorm"
	"github.com/pivotal-cf/brokerapi"
)

type serviceBroker struct {
	Client *client.Client
	Config config.Config
	Db     *gorm.DB
	Opts   cmd.CommandOpts
	Logger lager.Logger
}

func NewServiceBroker(
	client *client.Client,
	opts cmd.CommandOpts,
	config config.Config,
	db *gorm.DB,
	logger lager.Logger) *serviceBroker {

	return &serviceBroker{
		Client: client,
		Config: config,
		Opts:   opts,
		Db:     db,
		Logger: logger,
	}
}

func (b *serviceBroker) Services(context context.Context) []brokerapi.Service {
	planList := []brokerapi.ServicePlan{}

	for _, plan := range b.Config.Plans {
		planList = append(planList, brokerapi.ServicePlan{
			ID:          plan.Name,
			Name:        plan.Name,
			Description: plan.Description,
		})
	}
	b.Logger.Info("Serving a catalog request")
	return []brokerapi.Service{
		{
			ID:            b.Opts.ServiceID,
			Name:          b.Opts.Name,
			Description:   b.Opts.Description,
			Bindable:      true,
			Tags:          []string{"chaos-loris", "test"},
			Plans:         planList,
			PlanUpdatable: true,
		},
	}
}

func (b *serviceBroker) Provision(context context.Context, instanceId string, provisionDetails brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	planName := provisionDetails.PlanID
	plan, err := b.Config.PlanByName(planName)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
	}
	b.Logger.Info("Starting provisioning a service instance", lager.Data{
		"instance-id":       instanceId,
		"plan-id":           plan.Name,
		"plan-desctiption":  plan.Description,
		"organization-guid": provisionDetails.OrganizationGUID,
		"space-guid":        provisionDetails.SpaceGUID,
	})

	scheduleUrl, err := b.Client.CreateSchedule(instanceId, plan.Schedule)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
	}

	serivceInstance := model.ServiceInstance{
		InstanceId:  instanceId,
		ScheduleUrl: scheduleUrl,
		Probability: plan.Probability,
	}

	err = b.Db.Create(&serivceInstance).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
	}

	b.Logger.Info("service instance created", lager.Data{
		"instance-id": instanceId,
		"scheduleUrl": scheduleUrl,
	})

	return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
}

func (b *serviceBroker) Update(context context.Context, instanceId string, updateDetails brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{IsAsync: false}, errors.New("Not implemented")
}

func (b *serviceBroker) Deprovision(context context.Context, instanceId string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	var serivceInstance model.ServiceInstance
	err := b.Db.First(&serivceInstance, "instance_id = ?", instanceId).Error
	err = b.Client.Delete(serivceInstance.ScheduleUrl)
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{IsAsync: false}, err
	}
	err = b.Db.Delete(&serivceInstance).Error
	b.Logger.Info("service instance is removed", lager.Data{
		"instance-id": instanceId,
	})
	return brokerapi.DeprovisionServiceSpec{IsAsync: false}, err
}

func (b *serviceBroker) Bind(context context.Context, instanceId, bindingId string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	var serivceInstance model.ServiceInstance
	err := b.Db.First(&serivceInstance, "instance_id = ?", instanceId).Error
	if err != nil {
		return brokerapi.Binding{}, err
	}
	appUrl, err := b.Client.CreateApp(details.AppGUID)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	chaosUrl, err := b.Client.CreateChaos(appUrl, serivceInstance.ScheduleUrl, serivceInstance.Probability)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	serviceBinding := model.ServiceBinding{
		BindingId:      bindingId,
		ChaosUrl:       chaosUrl,
		ApplicationUrl: appUrl,
		InstanceId:     instanceId,
	}
	err = b.Db.Create(&serviceBinding).Error
	if err != nil {
		return brokerapi.Binding{}, err
	}
	b.Logger.Info("service binding created", lager.Data{
		"instance-id": instanceId,
		"binding-id":  bindingId,
	})
	return brokerapi.Binding{}, err
}

func (b *serviceBroker) Unbind(context context.Context, instanceId, bindingId string, details brokerapi.UnbindDetails) error {
	var serivceBinding model.ServiceBinding
	err := b.Db.First(&serivceBinding, "binding_id = ?", bindingId).Error
	if err != nil {
		return err
	}

	err = b.Client.Delete(serivceBinding.ChaosUrl)
	if err != nil {
		return err
	}

	err = b.Client.Delete(serivceBinding.ApplicationUrl)
	if err != nil {
		return err
	}

	err = b.Db.Delete(&serivceBinding).Error

	return err
}

func (b *serviceBroker) LastOperation(context context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
