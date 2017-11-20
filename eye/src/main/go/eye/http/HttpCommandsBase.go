package http

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
func (o *Create) AggregateType() eventhorizon.AggregateType  { return HttpServiceAggregateType }
func (o *Create) CommandType() eventhorizon.CommandType      { return CreateCommand }



        
type Delete struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Delete) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Delete) AggregateType() eventhorizon.AggregateType  { return HttpServiceAggregateType }
func (o *Delete) CommandType() eventhorizon.CommandType      { return DeleteCommand }



        
type Update struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Update) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Update) AggregateType() eventhorizon.AggregateType  { return HttpServiceAggregateType }
func (o *Update) CommandType() eventhorizon.CommandType      { return UpdateCommand }





type HttpServiceCommandType struct {
	name  string
	ordinal int
}

func (o *HttpServiceCommandType) Name() string {
    return o.name
}

func (o *HttpServiceCommandType) Ordinal() int {
    return o.ordinal
}

func (o HttpServiceCommandType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *HttpServiceCommandType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := HttpServiceCommandTypes().ParseHttpServiceCommandType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid HttpServiceCommandType %q", lit.Name)
        }
	}
	return
}

func (o HttpServiceCommandType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *HttpServiceCommandType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := HttpServiceCommandTypes().ParseHttpServiceCommandType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid HttpServiceCommandType %q", lit)
        }
    }
    return
}

func (o *HttpServiceCommandType) IsCreate() bool {
    return o == _httpServiceCommandTypes.Create()
}

func (o *HttpServiceCommandType) IsDelete() bool {
    return o == _httpServiceCommandTypes.Delete()
}

func (o *HttpServiceCommandType) IsUpdate() bool {
    return o == _httpServiceCommandTypes.Update()
}

type httpServiceCommandTypes struct {
	values []*HttpServiceCommandType
    literals []enum.Literal
}

var _httpServiceCommandTypes = &httpServiceCommandTypes{values: []*HttpServiceCommandType{
    {name: "Create", ordinal: 0},
    {name: "Delete", ordinal: 1},
    {name: "Update", ordinal: 2}},
}

func HttpServiceCommandTypes() *httpServiceCommandTypes {
	return _httpServiceCommandTypes
}

func (o *httpServiceCommandTypes) Values() []*HttpServiceCommandType {
	return o.values
}

func (o *httpServiceCommandTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *httpServiceCommandTypes) Create() *HttpServiceCommandType {
    return _httpServiceCommandTypes.values[0]
}

func (o *httpServiceCommandTypes) Delete() *HttpServiceCommandType {
    return _httpServiceCommandTypes.values[1]
}

func (o *httpServiceCommandTypes) Update() *HttpServiceCommandType {
    return _httpServiceCommandTypes.values[2]
}

func (o *httpServiceCommandTypes) ParseHttpServiceCommandType(name string) (ret *HttpServiceCommandType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*HttpServiceCommandType), ok
	}
	return
}



