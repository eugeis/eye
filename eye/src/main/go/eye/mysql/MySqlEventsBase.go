package mysql

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




type MySqlServiceEventType struct {
	name  string
	ordinal int
}

func (o *MySqlServiceEventType) Name() string {
    return o.name
}

func (o *MySqlServiceEventType) Ordinal() int {
    return o.ordinal
}

func (o MySqlServiceEventType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *MySqlServiceEventType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := MySqlServiceEventTypes().ParseMySqlServiceEventType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid MySqlServiceEventType %q", lit.Name)
        }
	}
	return
}

func (o MySqlServiceEventType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *MySqlServiceEventType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := MySqlServiceEventTypes().ParseMySqlServiceEventType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid MySqlServiceEventType %q", lit)
        }
    }
    return
}

func (o *MySqlServiceEventType) IsCreate() bool {
    return o == _mySqlServiceEventTypes.Create()
}

func (o *MySqlServiceEventType) IsCreated() bool {
    return o == _mySqlServiceEventTypes.Created()
}

func (o *MySqlServiceEventType) IsDelete() bool {
    return o == _mySqlServiceEventTypes.Delete()
}

func (o *MySqlServiceEventType) IsDeleted() bool {
    return o == _mySqlServiceEventTypes.Deleted()
}

func (o *MySqlServiceEventType) IsUpdate() bool {
    return o == _mySqlServiceEventTypes.Update()
}

func (o *MySqlServiceEventType) IsUpdated() bool {
    return o == _mySqlServiceEventTypes.Updated()
}

type mySqlServiceEventTypes struct {
	values []*MySqlServiceEventType
    literals []enum.Literal
}

var _mySqlServiceEventTypes = &mySqlServiceEventTypes{values: []*MySqlServiceEventType{
    {name: "Create", ordinal: 0},
    {name: "Created", ordinal: 1},
    {name: "Delete", ordinal: 2},
    {name: "Deleted", ordinal: 3},
    {name: "Update", ordinal: 4},
    {name: "Updated", ordinal: 5}},
}

func MySqlServiceEventTypes() *mySqlServiceEventTypes {
	return _mySqlServiceEventTypes
}

func (o *mySqlServiceEventTypes) Values() []*MySqlServiceEventType {
	return o.values
}

func (o *mySqlServiceEventTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *mySqlServiceEventTypes) Create() *MySqlServiceEventType {
    return _mySqlServiceEventTypes.values[0]
}

func (o *mySqlServiceEventTypes) Created() *MySqlServiceEventType {
    return _mySqlServiceEventTypes.values[1]
}

func (o *mySqlServiceEventTypes) Delete() *MySqlServiceEventType {
    return _mySqlServiceEventTypes.values[2]
}

func (o *mySqlServiceEventTypes) Deleted() *MySqlServiceEventType {
    return _mySqlServiceEventTypes.values[3]
}

func (o *mySqlServiceEventTypes) Update() *MySqlServiceEventType {
    return _mySqlServiceEventTypes.values[4]
}

func (o *mySqlServiceEventTypes) Updated() *MySqlServiceEventType {
    return _mySqlServiceEventTypes.values[5]
}

func (o *mySqlServiceEventTypes) ParseMySqlServiceEventType(name string) (ret *MySqlServiceEventType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*MySqlServiceEventType), ok
	}
	return
}



