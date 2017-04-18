package elssrv

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"errors"
	"github.com/galo/els-go/pkg/api"
	"github.com/galo/els-go/pkg/dynamodb/routingkeys"
	"golang.org/x/net/context"
	"github.com/hpcwp/els-go/dynamodb/routingkeys"
)

// Service describes a service that adds things together.
type ElsService interface {
	GetServiceInstanceByKey(ctx context.Context, *api.RoutingKey) (*api.ServiceInstance, error)

	// Add a routingKey to a service
	AddRoutingKey(context.Context, *api.AddRoutingKeyRequest) (*api.ServiceInstance, error)
}

type ServiceInstance struct {
	Url      string `json:"url"`
	Metadata string `json:"metadata"`
}

type basicElsService struct{
	rksrv *routingkeys.Service
}

// Errors
var (
	ErrNotFound = errors.New("ServiceInstance not found ")
)

// The implementation of the service
func (bs basicElsService) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (*api.ServiceInstance, error) {

	if routingKey.Id == "" {
		return &api.ServiceInstance{}, ErrInvalid
	}

	serviceInstance := bs.rksrv.Get(routingKey.Id)


	if serviceInstance == nil {
		return nil, ErrNotFound
	}
	if len(serviceInstance.ServiceInstances) == 0 {
		return nil, ErrNotFound
	}

	// We just return the first service url
	serviceUrl := serviceInstance.ServiceInstances[0].Uri
	if serviceUrl == nil {
		return nil, ErrNotFound
	}

	srvInstance := api.ServiceInstance{*serviceUrl, "rw"}
	return &srvInstance, nil
}

// The implementation of teh service
func (bs basicElsService) AddRoutingKey(ctx context.Context, addRoutingKeyRequest *api.AddRoutingKeyRequest) (*api.ServiceInstance, error) {
	if addRoutingKeyRequest.ServiceUri== "" {
		return &api.ServiceInstance{}, ErrInvalid
	}
	if addRoutingKeyRequest.RoutingKey== "" {
		return &api.ServiceInstance{}, ErrNotFound
	}


	instance := &ServiceInstance{addRoutingKeyRequest.ServiceUri, addRoutingKeyRequest.Tags}

	bs.rksrv.Add(instance, addRoutingKeyRequest.RoutingKey)

	return &api.ServiceInstance{instance.Url,instance.Metadata}, nil

}


const RoutingKeyTableName  = "routingKeys"


// NewBasicService returns a naïve dynamoDb implementation of Service.
func NewBasicService(tableName string, id string , secret string , token string) ElsService {
	rk := routingkeys.New(tableName, id, secret, token)

	return basicElsService{rk}
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(ElsService) ElsService

// ErrEmpty is returned when input is invalid
var ErrInvalid = errors.New("Invalid routing key")
