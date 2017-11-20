package mysql

import (
    "errors"
    "fmt"
    "github.com/eugeis/gee/eh"
    "github.com/looplab/eventhorizon"
    "github.com/looplab/eventhorizon/commandhandler/bus"
    "time"
)
type CommandHandler struct {
    CreateHandler func (*Create, *MySqlService, eh.AggregateStoreEvent) (err error)  `json:"createHandler" eh:"optional"`
    DeleteHandler func (*Delete, *MySqlService, eh.AggregateStoreEvent) (err error)  `json:"deleteHandler" eh:"optional"`
    UpdateHandler func (*Update, *MySqlService, eh.AggregateStoreEvent) (err error)  `json:"updateHandler" eh:"optional"`
}

func (o *CommandHandler) AddCreatePreparer(preparer func (*Create, *MySqlService) (err error) ) {
    prevHandler := o.CreateHandler
	o.CreateHandler = func(command *Create, entity *MySqlService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) AddDeletePreparer(preparer func (*Delete, *MySqlService) (err error) ) {
    prevHandler := o.DeleteHandler
	o.DeleteHandler = func(command *Delete, entity *MySqlService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) AddUpdatePreparer(preparer func (*Update, *MySqlService) (err error) ) {
    prevHandler := o.UpdateHandler
	o.UpdateHandler = func(command *Update, entity *MySqlService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) Execute(cmd eventhorizon.Command, entity eventhorizon.Entity, store eh.AggregateStoreEvent) (err error) {
    switch cmd.CommandType() {
    case CreateCommand:
        err = o.CreateHandler(cmd.(*Create), entity.(*MySqlService), store)
    case DeleteCommand:
        err = o.DeleteHandler(cmd.(*Delete), entity.(*MySqlService), store)
    case UpdateCommand:
        err = o.UpdateHandler(cmd.(*Update), entity.(*MySqlService), store)
    default:
		err = errors.New(fmt.Sprintf("Not supported command type '%v' for entity '%v", cmd.CommandType(), entity))
	}
    return
}

func (o *CommandHandler) SetupCommandHandler() (err error) {
    o.CreateHandler = func(command *Create, entity *MySqlService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateNewId(entity.Id, command.Id, MySqlServiceAggregateType); err == nil {
            store.StoreEvent(createdEvent, &Created{
                Name: command.Name,
                Id: command.Id,}, time.Now())
        }
        return
    }
    o.DeleteHandler = func(command *Delete, entity *MySqlService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, command.Id, MySqlServiceAggregateType); err == nil {
            store.StoreEvent(deletedEvent, &Deleted{
                Id: command.Id,}, time.Now())
        }
        return
    }
    o.UpdateHandler = func(command *Update, entity *MySqlService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, command.Id, MySqlServiceAggregateType); err == nil {
            store.StoreEvent(updatedEvent, &Updated{
                Name: command.Name,
                Id: command.Id,}, time.Now())
        }
        return
    }
    return
}


type EventHandler struct {
    CreateHandler func (*Create, *MySqlService) (err error)  `json:"createHandler" eh:"optional"`
    CreatedHandler func (*Created, *MySqlService) (err error)  `json:"createdHandler" eh:"optional"`
    DeleteHandler func (*Delete, *MySqlService) (err error)  `json:"deleteHandler" eh:"optional"`
    DeletedHandler func (*Deleted, *MySqlService) (err error)  `json:"deletedHandler" eh:"optional"`
    UpdateHandler func (*Update, *MySqlService) (err error)  `json:"updateHandler" eh:"optional"`
    UpdatedHandler func (*Updated, *MySqlService) (err error)  `json:"updatedHandler" eh:"optional"`
}

func (o *EventHandler) Apply(event eventhorizon.Event, entity eventhorizon.Entity) (err error) {
    switch event.EventType() {
    case CreateEvent:
        err = o.CreateHandler(event.Data().(*Create), entity.(*MySqlService))
    case CreatedEvent:
        err = o.CreatedHandler(event.Data().(*Created), entity.(*MySqlService))
    case DeleteEvent:
        err = o.DeleteHandler(event.Data().(*Delete), entity.(*MySqlService))
    case DeletedEvent:
        err = o.DeletedHandler(event.Data().(*Deleted), entity.(*MySqlService))
    case UpdateEvent:
        err = o.UpdateHandler(event.Data().(*Update), entity.(*MySqlService))
    case UpdatedEvent:
        err = o.UpdatedHandler(event.Data().(*Updated), entity.(*MySqlService))
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
    o.CreateHandler = func(event *Create, entity *MySqlService) (err error) {
        //err = eh.EventHandlerNotImplemented(CreateEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(CreatedEvent, func() eventhorizon.EventData {
		return &Created{}
	})

    //default handler implementation
    o.CreatedHandler = func(event *Created, entity *MySqlService) (err error) {
        if err = eh.ValidateNewId(entity.Id, event.Id, MySqlServiceAggregateType); err == nil {
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
    o.DeleteHandler = func(event *Delete, entity *MySqlService) (err error) {
        //err = eh.EventHandlerNotImplemented(DeleteEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(DeletedEvent, func() eventhorizon.EventData {
		return &Deleted{}
	})

    //default handler implementation
    o.DeletedHandler = func(event *Deleted, entity *MySqlService) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, event.Id, MySqlServiceAggregateType); err == nil {
            *entity = *NewMySqlService()
        }
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(UpdateEvent, func() eventhorizon.EventData {
		return &Update{}
	})

    //default handler implementation
    o.UpdateHandler = func(event *Update, entity *MySqlService) (err error) {
        //err = eh.EventHandlerNotImplemented(UpdateEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(UpdatedEvent, func() eventhorizon.EventData {
		return &Updated{}
	})

    //default handler implementation
    o.UpdatedHandler = func(event *Updated, entity *MySqlService) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, event.Id, MySqlServiceAggregateType); err == nil {
            entity.Name = event.Name
        }
        return
    }
    return
}


const MySqlServiceAggregateType eventhorizon.AggregateType = "MySqlService"

type AggregateInitializer struct {
    *eh.AggregateInitializer
    *CommandHandler
    *EventHandler
    ProjectorHandler *EventHandler `json:"projectorHandler" eh:"optional"`
}



func New@@EMPTY@@(eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus, eventPublisher eventhorizon.EventPublisher, 
                commandBus *bus.CommandHandler, 
                readRepos func (string, func () (ret eventhorizon.Entity) ) (ret eventhorizon.ReadWriteRepo) ) (ret *AggregateInitializer) {
    
    commandHandler := &MySqlServiceCommandHandler{}
    eventHandler := &MySqlServiceEventHandler{}
    entityFactory := func() eventhorizon.Entity { return NewMySqlService() }
    ret = &AggregateInitializer{AggregateInitializer: eh.NewAggregateInitializer(MySqlServiceAggregateType,
        func(id eventhorizon.UUID) eventhorizon.Aggregate {
            return eh.NewAggregateBase(MySqlServiceAggregateType, id, commandHandler, eventHandler, entityFactory())
        }, entityFactory,
        MySqlServiceCommandTypes().Literals(), MySqlServiceEventTypes().Literals(), eventHandler,
        []func() error{commandHandler.SetupCommandHandler, eventHandler.SetupEventHandler},
        eventStore, eventBus, eventPublisher, commandBus, readRepos), MySqlServiceCommandHandler: commandHandler, MySqlServiceEventHandler: eventHandler, ProjectorHandler: eventHandler,
    }

    return
}


type MySqlEventhorizonInitializer struct {
    eventStore eventhorizon.EventStore `json:"eventStore" eh:"optional"`
    eventBus eventhorizon.EventBus `json:"eventBus" eh:"optional"`
    eventPublisher eventhorizon.EventPublisher `json:"eventPublisher" eh:"optional"`
    commandBus *bus.CommandHandler `json:"commandBus" eh:"optional"`
    MySqlServiceAggregateInitializer *AggregateInitializer `json:"mySqlServiceAggregateInitializer" eh:"optional"`
}

func New@@EMPTY@@(eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus, eventPublisher eventhorizon.EventPublisher, 
                commandBus *bus.CommandHandler, 
                readRepos func (string, func () (ret eventhorizon.Entity) ) (ret eventhorizon.ReadWriteRepo) ) (ret *MySqlEventhorizonInitializer) {
    mySqlServiceAggregateInitializer := New@@EMPTY@@(eventStore, eventBus, eventPublisher, commandBus, readRepos)
    ret = &MySqlEventhorizonInitializer{
        eventStore: eventStore,
        eventBus: eventBus,
        eventPublisher: eventPublisher,
        commandBus: commandBus,
        MySqlServiceAggregateInitializer: mySqlServiceAggregateInitializer,
    }
    return
}

func (o *MySqlEventhorizonInitializer) Setup() (err error) {
    
    if err = o.MySqlServiceAggregateInitializer.Setup(); err != nil {
        return
    }

    return
}









