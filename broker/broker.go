package broker

import (
	"github.com/Altoros/cf-chaos-loris-broker/config"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

type ServiceInstanceCreator interface {
	Create(instanceID string, settings map[string]interface{}) error
	Update(instanceID string, params map[string]interface{}) error
	Destroy(instanceID string) error
	InstanceExists(instanceID string) (bool, error)
}

type ServiceInstanceBinder interface {
	Bind(instanceID string, bindingID string) (interface{}, error)
	Unbind(instanceID string, bindingID string) error
	InstanceExists(instanceID string) (bool, error)
}

type serviceBroker struct {
	InstanceCreator ServiceInstanceCreator
	InstanceBinder  ServiceInstanceBinder
	// StatePersister  persisters.StatePersister
	Config config.Config
	Logger lager.Logger
}

func NewServiceBroker(
	// instanceCreator ServiceInstanceCreator,
	// instanceBinder ServiceInstanceBinder,
	// state,
	conf config.Config,
	logger lager.Logger) *serviceBroker {

	return &serviceBroker{
		// InstanceCreator: instanceCreator,
		// InstanceBinder:  instanceBinder,
		// StatePersister:  statePersister,
		Config: conf,
		Logger: logger,
	}
}

func (b *serviceBroker) Services() []brokerapi.Service {
	planList := []brokerapi.ServicePlan{}
	b.Logger.Info("Serving a catalog request")
	return []brokerapi.Service{
		{
			ID:            b.Config.ServiceBroker.ServiceID,
			Name:          b.Config.ServiceBroker.Name,
			Description:   b.Config.ServiceBroker.Description,
			Bindable:      true,
			Tags:          []string{"broker"},
			Plans:         planList,
			PlanUpdatable: true,
		},
	}
}

func (b *serviceBroker) Provision(instanceID string, provisionDetails brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	settings := map[string]interface{}{}
	return brokerapi.ProvisionedServiceSpec{IsAsync: false}, b.InstanceCreator.Create(instanceID, settings)
}

func (b *serviceBroker) Update(instanceID string, updateDetails brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	params := map[string]interface{}{}
	return brokerapi.UpdateServiceSpec{IsAsync: false}, b.InstanceCreator.Update(instanceID, params)
}

func (b *serviceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{IsAsync: false}, b.InstanceCreator.Destroy(instanceID)
}

func (b *serviceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	b.Logger.Info("Looking for the service credentials", lager.Data{
		"instance-id": instanceID,
		"binding-id":  bindingID,
		"details":     details,
	})
	creds, err := b.InstanceBinder.Bind(instanceID, bindingID)
	return brokerapi.Binding{Credentials: creds}, err
}

func (b *serviceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}

func (b *serviceBroker) LastOperation(instanceID string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
