package http

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




type HttpServiceEventType struct {
	name  string
	ordinal int
}

func (o *HttpServiceEventType) Name() string {
    return o.name
}

func (o *HttpServiceEventType) Ordinal() int {
    return o.ordinal
}

func (o HttpServiceEventType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *HttpServiceEventType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := HttpServiceEventTypes().ParseHttpServiceEventType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid HttpServiceEventType %q", lit.Name)
        }
	}
	return
}

func (o HttpServiceEventType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *HttpServiceEventType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := HttpServiceEventTypes().ParseHttpServiceEventType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid HttpServiceEventType %q", lit)
        }
    }
    return
}

func (o *HttpServiceEventType) IsCreate() bool {
    return o == _httpServiceEventTypes.Create()
}

func (o *HttpServiceEventType) IsCreated() bool {
    return o == _httpServiceEventTypes.Created()
}

func (o *HttpServiceEventType) IsDelete() bool {
    return o == _httpServiceEventTypes.Delete()
}

func (o *HttpServiceEventType) IsDeleted() bool {
    return o == _httpServiceEventTypes.Deleted()
}

func (o *HttpServiceEventType) IsUpdate() bool {
    return o == _httpServiceEventTypes.Update()
}

func (o *HttpServiceEventType) IsUpdated() bool {
    return o == _httpServiceEventTypes.Updated()
}

type httpServiceEventTypes struct {
	values []*HttpServiceEventType
    literals []enum.Literal
}

var _httpServiceEventTypes = &httpServiceEventTypes{values: []*HttpServiceEventType{
    {name: "Create", ordinal: 0},
    {name: "Created", ordinal: 1},
    {name: "Delete", ordinal: 2},
    {name: "Deleted", ordinal: 3},
    {name: "Update", ordinal: 4},
    {name: "Updated", ordinal: 5}},
}

func HttpServiceEventTypes() *httpServiceEventTypes {
	return _httpServiceEventTypes
}

func (o *httpServiceEventTypes) Values() []*HttpServiceEventType {
	return o.values
}

func (o *httpServiceEventTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *httpServiceEventTypes) Create() *HttpServiceEventType {
    return _httpServiceEventTypes.values[0]
}

func (o *httpServiceEventTypes) Created() *HttpServiceEventType {
    return _httpServiceEventTypes.values[1]
}

func (o *httpServiceEventTypes) Delete() *HttpServiceEventType {
    return _httpServiceEventTypes.values[2]
}

func (o *httpServiceEventTypes) Deleted() *HttpServiceEventType {
    return _httpServiceEventTypes.values[3]
}

func (o *httpServiceEventTypes) Update() *HttpServiceEventType {
    return _httpServiceEventTypes.values[4]
}

func (o *httpServiceEventTypes) Updated() *HttpServiceEventType {
    return _httpServiceEventTypes.values[5]
}

func (o *httpServiceEventTypes) ParseHttpServiceEventType(name string) (ret *HttpServiceEventType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*HttpServiceEventType), ok
	}
	return
}



