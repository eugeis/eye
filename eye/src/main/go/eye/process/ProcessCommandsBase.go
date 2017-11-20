package process

import (
    "encoding/json"
    "fmt"
    "github.com/eugeis/gee/enum"
    "github.com/looplab/eventhorizon"
    "gopkg.in/mgo.v2/bson"
)
const (
     CreateCommand eventhorizon.CommandType = "Create"
     DeleteCommand eventhorizon.CommandType = "Delete"
     UpdateCommand eventhorizon.CommandType = "Update"
)




        
type Create struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Create) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Create) AggregateType() eventhorizon.AggregateType  { return ProcessServiceAggregateType }
func (o *Create) CommandType() eventhorizon.CommandType      { return CreateCommand }



        
type Delete struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Delete) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Delete) AggregateType() eventhorizon.AggregateType  { return ProcessServiceAggregateType }
func (o *Delete) CommandType() eventhorizon.CommandType      { return DeleteCommand }



        
type Update struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Update) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Update) AggregateType() eventhorizon.AggregateType  { return ProcessServiceAggregateType }
func (o *Update) CommandType() eventhorizon.CommandType      { return UpdateCommand }





type ProcessServiceCommandType struct {
	name  string
	ordinal int
}

func (o *ProcessServiceCommandType) Name() string {
    return o.name
}

func (o *ProcessServiceCommandType) Ordinal() int {
    return o.ordinal
}

func (o ProcessServiceCommandType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *ProcessServiceCommandType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := ProcessServiceCommandTypes().ParseProcessServiceCommandType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ProcessServiceCommandType %q", lit.Name)
        }
	}
	return
}

func (o ProcessServiceCommandType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *ProcessServiceCommandType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := ProcessServiceCommandTypes().ParseProcessServiceCommandType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ProcessServiceCommandType %q", lit)
        }
    }
    return
}

func (o *ProcessServiceCommandType) IsCreate() bool {
    return o == _processServiceCommandTypes.Create()
}

func (o *ProcessServiceCommandType) IsDelete() bool {
    return o == _processServiceCommandTypes.Delete()
}

func (o *ProcessServiceCommandType) IsUpdate() bool {
    return o == _processServiceCommandTypes.Update()
}

type processServiceCommandTypes struct {
	values []*ProcessServiceCommandType
    literals []enum.Literal
}

var _processServiceCommandTypes = &processServiceCommandTypes{values: []*ProcessServiceCommandType{
    {name: "Create", ordinal: 0},
    {name: "Delete", ordinal: 1},
    {name: "Update", ordinal: 2}},
}

func ProcessServiceCommandTypes() *processServiceCommandTypes {
	return _processServiceCommandTypes
}

func (o *processServiceCommandTypes) Values() []*ProcessServiceCommandType {
	return o.values
}

func (o *processServiceCommandTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *processServiceCommandTypes) Create() *ProcessServiceCommandType {
    return _processServiceCommandTypes.values[0]
}

func (o *processServiceCommandTypes) Delete() *ProcessServiceCommandType {
    return _processServiceCommandTypes.values[1]
}

func (o *processServiceCommandTypes) Update() *ProcessServiceCommandType {
    return _processServiceCommandTypes.values[2]
}

func (o *processServiceCommandTypes) ParseProcessServiceCommandType(name string) (ret *ProcessServiceCommandType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*ProcessServiceCommandType), ok
	}
	return
}



