package process

import (
    "encoding/json"
    "fmt"
    "github.com/eugeis/gee/enum"
    "github.com/looplab/eventhorizon"
    "gopkg.in/mgo.v2/bson"
)
const (
     CreateEvent eventhorizon.EventType = "Create"
     CreatedEvent eventhorizon.EventType = "Created"
     DeleteEvent eventhorizon.EventType = "Delete"
     DeletedEvent eventhorizon.EventType = "Deleted"
     UpdateEvent eventhorizon.EventType = "Update"
     UpdatedEvent eventhorizon.EventType = "Updated"
)




type Create struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}


type Created struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}


type Delete struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}


type Deleted struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}


type Update struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}


type Updated struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}




type ProcessServiceEventType struct {
	name  string
	ordinal int
}

func (o *ProcessServiceEventType) Name() string {
    return o.name
}

func (o *ProcessServiceEventType) Ordinal() int {
    return o.ordinal
}

func (o ProcessServiceEventType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *ProcessServiceEventType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := ProcessServiceEventTypes().ParseProcessServiceEventType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ProcessServiceEventType %q", lit.Name)
        }
	}
	return
}

func (o ProcessServiceEventType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *ProcessServiceEventType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := ProcessServiceEventTypes().ParseProcessServiceEventType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ProcessServiceEventType %q", lit)
        }
    }
    return
}

func (o *ProcessServiceEventType) IsCreate() bool {
    return o == _processServiceEventTypes.Create()
}

func (o *ProcessServiceEventType) IsCreated() bool {
    return o == _processServiceEventTypes.Created()
}

func (o *ProcessServiceEventType) IsDelete() bool {
    return o == _processServiceEventTypes.Delete()
}

func (o *ProcessServiceEventType) IsDeleted() bool {
    return o == _processServiceEventTypes.Deleted()
}

func (o *ProcessServiceEventType) IsUpdate() bool {
    return o == _processServiceEventTypes.Update()
}

func (o *ProcessServiceEventType) IsUpdated() bool {
    return o == _processServiceEventTypes.Updated()
}

type processServiceEventTypes struct {
	values []*ProcessServiceEventType
    literals []enum.Literal
}

var _processServiceEventTypes = &processServiceEventTypes{values: []*ProcessServiceEventType{
    {name: "Create", ordinal: 0},
    {name: "Created", ordinal: 1},
    {name: "Delete", ordinal: 2},
    {name: "Deleted", ordinal: 3},
    {name: "Update", ordinal: 4},
    {name: "Updated", ordinal: 5}},
}

func ProcessServiceEventTypes() *processServiceEventTypes {
	return _processServiceEventTypes
}

func (o *processServiceEventTypes) Values() []*ProcessServiceEventType {
	return o.values
}

func (o *processServiceEventTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *processServiceEventTypes) Create() *ProcessServiceEventType {
    return _processServiceEventTypes.values[0]
}

func (o *processServiceEventTypes) Created() *ProcessServiceEventType {
    return _processServiceEventTypes.values[1]
}

func (o *processServiceEventTypes) Delete() *ProcessServiceEventType {
    return _processServiceEventTypes.values[2]
}

func (o *processServiceEventTypes) Deleted() *ProcessServiceEventType {
    return _processServiceEventTypes.values[3]
}

func (o *processServiceEventTypes) Update() *ProcessServiceEventType {
    return _processServiceEventTypes.values[4]
}

func (o *processServiceEventTypes) Updated() *ProcessServiceEventType {
    return _processServiceEventTypes.values[5]
}

func (o *processServiceEventTypes) ParseProcessServiceEventType(name string) (ret *ProcessServiceEventType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*ProcessServiceEventType), ok
	}
	return
}



