package elastic

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




type ElasticServiceEventType struct {
	name  string
	ordinal int
}

func (o *ElasticServiceEventType) Name() string {
    return o.name
}

func (o *ElasticServiceEventType) Ordinal() int {
    return o.ordinal
}

func (o ElasticServiceEventType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *ElasticServiceEventType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := ElasticServiceEventTypes().ParseElasticServiceEventType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ElasticServiceEventType %q", lit.Name)
        }
	}
	return
}

func (o ElasticServiceEventType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *ElasticServiceEventType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := ElasticServiceEventTypes().ParseElasticServiceEventType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ElasticServiceEventType %q", lit)
        }
    }
    return
}

func (o *ElasticServiceEventType) IsCreate() bool {
    return o == _elasticServiceEventTypes.Create()
}

func (o *ElasticServiceEventType) IsCreated() bool {
    return o == _elasticServiceEventTypes.Created()
}

func (o *ElasticServiceEventType) IsDelete() bool {
    return o == _elasticServiceEventTypes.Delete()
}

func (o *ElasticServiceEventType) IsDeleted() bool {
    return o == _elasticServiceEventTypes.Deleted()
}

func (o *ElasticServiceEventType) IsUpdate() bool {
    return o == _elasticServiceEventTypes.Update()
}

func (o *ElasticServiceEventType) IsUpdated() bool {
    return o == _elasticServiceEventTypes.Updated()
}

type elasticServiceEventTypes struct {
	values []*ElasticServiceEventType
    literals []enum.Literal
}

var _elasticServiceEventTypes = &elasticServiceEventTypes{values: []*ElasticServiceEventType{
    {name: "Create", ordinal: 0},
    {name: "Created", ordinal: 1},
    {name: "Delete", ordinal: 2},
    {name: "Deleted", ordinal: 3},
    {name: "Update", ordinal: 4},
    {name: "Updated", ordinal: 5}},
}

func ElasticServiceEventTypes() *elasticServiceEventTypes {
	return _elasticServiceEventTypes
}

func (o *elasticServiceEventTypes) Values() []*ElasticServiceEventType {
	return o.values
}

func (o *elasticServiceEventTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *elasticServiceEventTypes) Create() *ElasticServiceEventType {
    return _elasticServiceEventTypes.values[0]
}

func (o *elasticServiceEventTypes) Created() *ElasticServiceEventType {
    return _elasticServiceEventTypes.values[1]
}

func (o *elasticServiceEventTypes) Delete() *ElasticServiceEventType {
    return _elasticServiceEventTypes.values[2]
}

func (o *elasticServiceEventTypes) Deleted() *ElasticServiceEventType {
    return _elasticServiceEventTypes.values[3]
}

func (o *elasticServiceEventTypes) Update() *ElasticServiceEventType {
    return _elasticServiceEventTypes.values[4]
}

func (o *elasticServiceEventTypes) Updated() *ElasticServiceEventType {
    return _elasticServiceEventTypes.values[5]
}

func (o *elasticServiceEventTypes) ParseElasticServiceEventType(name string) (ret *ElasticServiceEventType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*ElasticServiceEventType), ok
	}
	return
}



