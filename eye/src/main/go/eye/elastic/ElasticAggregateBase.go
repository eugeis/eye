package elastic

import (
    "errors"
    "fmt"
    "github.com/eugeis/gee/eh"
    "github.com/looplab/eventhorizon"
    "github.com/looplab/eventhorizon/commandhandler/bus"
    "time"
)
type CommandHandler struct {
    CreateHandler func (*Create, *ElasticService, eh.AggregateStoreEvent) (err error)  `json:"createHandler" eh:"optional"`
    DeleteHandler func (*Delete, *ElasticService, eh.AggregateStoreEvent) (err error)  `json:"deleteHandler" eh:"optional"`
    UpdateHandler func (*Update, *ElasticService, eh.AggregateStoreEvent) (err error)  `json:"updateHandler" eh:"optional"`
}

func (o *CommandHandler) AddCreatePreparer(preparer func (*Create, *ElasticService) (err error) ) {
    prevHandler := o.CreateHandler
	o.CreateHandler = func(command *Create, entity *ElasticService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) AddDeletePreparer(preparer func (*Delete, *ElasticService) (err error) ) {
    prevHandler := o.DeleteHandler
	o.DeleteHandler = func(command *Delete, entity *ElasticService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) AddUpdatePreparer(preparer func (*Update, *ElasticService) (err error) ) {
    prevHandler := o.UpdateHandler
	o.UpdateHandler = func(command *Update, entity *ElasticService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) Execute(cmd eventhorizon.Command, entity eventhorizon.Entity, store eh.AggregateStoreEvent) (err error) {
    switch cmd.CommandType() {
    case CreateCommand:
        err = o.CreateHandler(cmd.(*Create), entity.(*ElasticService), store)
    case DeleteCommand:
        err = o.DeleteHandler(cmd.(*Delete), entity.(*ElasticService), store)
    case UpdateCommand:
        err = o.UpdateHandler(cmd.(*Update), entity.(*ElasticService), store)
    default:
		err = errors.New(fmt.Sprintf("Not supported command type '%v' for entity '%v", cmd.CommandType(), entity))
	}
    return
}

func (o *CommandHandler) SetupCommandHandler() (err error) {
    o.CreateHandler = func(command *Create, entity *ElasticService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateNewId(entity.Id, command.Id, ElasticServiceAggregateType); err == nil {
            store.StoreEvent(createdEvent, &Created{
                Name: command.Name,
                Id: command.Id,}, time.Now())
        }
        return
    }
    o.DeleteHandler = func(command *Delete, entity *ElasticService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, command.Id, ElasticServiceAggregateType); err == nil {
            store.StoreEvent(deletedEvent, &Deleted{
                Id: command.Id,}, time.Now())
        }
        return
    }
    o.UpdateHandler = func(command *Update, entity *ElasticService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, command.Id, ElasticServiceAggregateType); err == nil {
            store.StoreEvent(updatedEvent, &Updated{
                Name: command.Name,
                Id: command.Id,}, time.Now())
        }
        return
    }
    return
}


type EventHandler struct {
    CreateHandler func (*Create, *ElasticService) (err error)  `json:"createHandler" eh:"optional"`
    CreatedHandler func (*Created, *ElasticService) (err error)  `json:"createdHandler" eh:"optional"`
    DeleteHandler func (*Delete, *ElasticService) (err error)  `json:"deleteHandler" eh:"optional"`
    DeletedHandler func (*Deleted, *ElasticService) (err error)  `json:"deletedHandler" eh:"optional"`
    UpdateHandler func (*Update, *ElasticService) (err error)  `json:"updateHandler" eh:"optional"`
    UpdatedHandler func (*Updated, *ElasticService) (err error)  `json:"updatedHandler" eh:"optional"`
}

func (o *EventHandler) Apply(event eventhorizon.Event, entity eventhorizon.Entity) (err error) {
    switch event.EventType() {
    case CreateEvent:
        err = o.CreateHandler(event.Data().(*Create), entity.(*ElasticService))
    case CreatedEvent:
        err = o.CreatedHandler(event.Data().(*Created), entity.(*ElasticService))
    case DeleteEvent:
        err = o.DeleteHandler(event.Data().(*Delete), entity.(*ElasticService))
    case DeletedEvent:
        err = o.DeletedHandler(event.Data().(*Deleted), entity.(*ElasticService))
    case UpdateEvent:
        err = o.UpdateHandler(event.Data().(*Update), entity.(*ElasticService))
    case UpdatedEvent:
        err = o.UpdatedHandler(event.Data().(*Updated), entity.(*ElasticService))
    default:
		err = errors.New(fmt.Sprintf("Not supported event type '%v' for entity '%v", event.EventType(), entity))
	}
    return
}

func (o *EventHandler) SetupEventHandler() (err error) {

    //register event object factory
    eventhorizon.RegisterEventData(CreateEvent, func() eventhorizon.EventData {
		return &Create{}
	})

    //default handler implementation
    o.CreateHandler = func(event *Create, entity *ElasticService) (err error) {
        //err = eh.EventHandlerNotImplemented(CreateEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(CreatedEvent, func() eventhorizon.EventData {
		return &Created{}
	})

    //default handler implementation
    o.CreatedHandler = func(event *Created, entity *ElasticService) (err error) {
        if err = eh.ValidateNewId(entity.Id, event.Id, ElasticServiceAggregateType); err == nil {
            entity.Name = event.Name
            entity.Id = event.Id
        }
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(DeleteEvent, func() eventhorizon.EventData {
		return &Delete{}
	})

    //default handler implementation
    o.DeleteHandler = func(event *Delete, entity *ElasticService) (err error) {
        //err = eh.EventHandlerNotImplemented(DeleteEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(DeletedEvent, func() eventhorizon.EventData {
		return &Deleted{}
	})

    //default handler implementation
    o.DeletedHandler = func(event *Deleted, entity *ElasticService) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, event.Id, ElasticServiceAggregateType); err == nil {
            *entity = *NewElasticService()
        }
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(UpdateEvent, func() eventhorizon.EventData {
		return &Update{}
	})

    //default handler implementation
    o.UpdateHandler = func(event *Update, entity *ElasticService) (err error) {
        //err = eh.EventHandlerNotImplemented(UpdateEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(UpdatedEvent, func() eventhorizon.EventData {
		return &Updated{}
	})

    //default handler implementation
    o.UpdatedHandler = func(event *Updated, entity *ElasticService) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, event.Id, ElasticServiceAggregateType); err == nil {
            entity.Name = event.Name
        }
        return
    }
    return
}


const ElasticServiceAggregateType eventhorizon.AggregateType = "ElasticService"

type AggregateInitializer struct {
    *eh.AggregateInitializer
    *CommandHandler
    *EventHandler
    ProjectorHandler *EventHandler `json:"projectorHandler" eh:"optional"`
}



func New@@EMPTY@@(eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus, eventPublisher eventhorizon.EventPublisher, 
                commandBus *bus.CommandHandler, 
                readRepos func (string, func () (ret eventhorizon.Entity) ) (ret eventhorizon.ReadWriteRepo) ) (ret *AggregateInitializer) {
    
    commandHandler := &ElasticServiceCommandHandler{}
    eventHandler := &ElasticServiceEventHandler{}
    entityFactory := func() eventhorizon.Entity { return NewElasticService() }
    ret = &AggregateInitializer{AggregateInitializer: eh.NewAggregateInitializer(ElasticServiceAggregateType,
        func(id eventhorizon.UUID) eventhorizon.Aggregate {
            return eh.NewAggregateBase(ElasticServiceAggregateType, id, commandHandler, eventHandler, entityFactory())
        }, entityFactory,
        ElasticServiceCommandTypes().Literals(), ElasticServiceEventTypes().Literals(), eventHandler,
        []func() error{commandHandler.SetupCommandHandler, eventHandler.SetupEventHandler},
        eventStore, eventBus, eventPublisher, commandBus, readRepos), ElasticServiceCommandHandler: commandHandler, ElasticServiceEventHandler: eventHandler, ProjectorHandler: eventHandler,
    }

    return
}


type ElasticEventhorizonInitializer struct {
    eventStore eventhorizon.EventStore `json:"eventStore" eh:"optional"`
    eventBus eventhorizon.EventBus `json:"eventBus" eh:"optional"`
    eventPublisher eventhorizon.EventPublisher `json:"eventPublisher" eh:"optional"`
    commandBus *bus.CommandHandler `json:"commandBus" eh:"optional"`
    ElasticServiceAggregateInitializer *AggregateInitializer `json:"elasticServiceAggregateInitializer" eh:"optional"`
}

func New@@EMPTY@@(eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus, eventPublisher eventhorizon.EventPublisher, 
                commandBus *bus.CommandHandler, 
                readRepos func (string, func () (ret eventhorizon.Entity) ) (ret eventhorizon.ReadWriteRepo) ) (ret *ElasticEventhorizonInitializer) {
    elasticServiceAggregateInitializer := New@@EMPTY@@(eventStore, eventBus, eventPublisher, commandBus, readRepos)
    ret = &ElasticEventhorizonInitializer{
        eventStore: eventStore,
        eventBus: eventBus,
        eventPublisher: eventPublisher,
        commandBus: commandBus,
        ElasticServiceAggregateInitializer: elasticServiceAggregateInitializer,
    }
    return
}

func (o *ElasticEventhorizonInitializer) Setup() (err error) {
    
    if err = o.ElasticServiceAggregateInitializer.Setup(); err != nil {
        return
    }

    return
}









