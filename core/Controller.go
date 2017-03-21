package core

import (
	"rest/integ"
	"fmt"
)

type Operator interface {
}

type Controller struct {
	config         *RestConfig
	serviceFactory Factory
	commandCache   integ.Cache
}

func NewController(config *RestConfig) Controller {
	return Controller{config: config, serviceFactory: config.ServiceFactory(), commandCache: integ.NewCache()}
}

func (o Controller) Close() {
	if o.serviceFactory != nil {
		o.serviceFactory.Close()
	}
	o.commandCache = nil
	o.config = nil
	o.serviceFactory = nil
}

func (o Controller) Ping(serviceName string) (err error) {
	service, err := o.serviceFactory.Find(serviceName)
	if err == nil {
		err = service.Ping()
	}
	return
}

func (o Controller) Validate(serviceName string, req *QueryRequest) (err error) {
	if req.Query == "" {
		log.Debug(fmt.Sprintf("ping instead of validate, because no query defined for %s", serviceName))
		return o.Ping(serviceName)
	}

	service, err := o.serviceFactory.Find(serviceName)
	if err == nil {
		var value interface{}
		value, err = o.commandCache.GetOrBuild(req.CommandKey(service.Name()), func() (interface{}, error) {
			return service.NewСheck(req)
		})
		if err == nil {
			command, _ := value.(Check)
			err = command.Validate()
		}
	}
	return
}

func (o Controller) PingAny(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err == nil {
			break
		}
	}
	return
}

func (o Controller) PingAll(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err != nil {
			break
		}
	}
	return
}

func (o Controller) ValidateAny(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err == nil {
			break
		}
	}
	return
}

func (o Controller) ValidateRunning(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err := o.Ping(serviceName)
		if err == nil {
			err = o.Validate(serviceName, req)
			if err != nil {
				break
			}
		}
	}
	return
}

func (o Controller) ValidateAll(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err != nil {
			break
		}
	}
	return
}

func (o Controller) CompareAny(serviceNames []string, req *CompareRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req.QueryRequest)
		if err == nil {
			break
		}
	}
	return
}

func (o Controller) CompareRunning(serviceNames []string, req *CompareRequest) (err error) {
	for _, serviceName := range serviceNames {
		err := o.Ping(serviceName)
		if err == nil {
			//err = o.Validate(serviceName, req)
			if err != nil {
				break
			}
		}
	}
	return
}

func (o Controller) CompareAll(serviceNames []string, req *CompareRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req.QueryRequest)
		if err != nil {
			break
		}
	}
	return
}

func (o Controller) query(serviceName string, req *QueryRequest) (data []byte, err error) {
	service, err := o.serviceFactory.Find(serviceName)
	if err == nil {
		var value interface{}
		value, err = o.commandCache.GetOrBuild(req.CommandKey(service.Name()), func() (interface{}, error) {
			return service.NewСheck(req)
		})
		if err == nil {
			command, _ := value.(Check)
			data, err = command.Query()
		}
	}
	return
}
