package core

import (
	"fmt"
)

func (o *Eye) registerExecutors() {
	for _, item := range o.config.Executor {
		if len(item.Services) > 1 {
			for _, serviceName := range item.Services {
				o.registerSimpleExecutor(fmt.Sprintf("%v-%v", serviceName, item.Name), serviceName, item)
			}
		} else if len(item.Services) > 0 {
			o.registerSimpleExecutor(item.Name, item.Services[0], item)
		} else {
			Log.Info("No service defined for the executor %v", item.Name)
		}
	}
}

func (o *Eye) registerSimpleExecutor(executorFullName string, serviceName string, executor *SimpleExecutor) {
	var err error
	var service Service
	if service, err = o.serviceFactory.Find(serviceName); err == nil {
		var item Executor
		request := &CommandRequest{}
		if item, err = service.NewExecutor(request); err == nil {
			o.executors[executorFullName] = item
		}
	}
}
