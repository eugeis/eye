package mysql

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
func (o *Create) AggregateType() eventhorizon.AggregateType  { return MySqlServiceAggregateType }
func (o *Create) CommandType() eventhorizon.CommandType      { return CreateCommand }



        
type Delete struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Delete) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Delete) AggregateType() eventhorizon.AggregateType  { return MySqlServiceAggregateType }
func (o *Delete) CommandType() eventhorizon.CommandType      { return DeleteCommand }



        
type Update struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Update) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Update) AggregateType() eventhorizon.AggregateType  { return MySqlServiceAggregateType }
func (o *Update) CommandType() eventhorizon.CommandType      { return UpdateCommand }





type MySqlServiceCommandType struct {
	name  string
	ordinal int
}

func (o *MySqlServiceCommandType) Name() string {
    return o.name
}

func (o *MySqlServiceCommandType) Ordinal() int {
    return o.ordinal
}

func (o MySqlServiceCommandType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *MySqlServiceCommandType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := MySqlServiceCommandTypes().ParseMySqlServiceCommandType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid MySqlServiceCommandType %q", lit.Name)
        }
	}
	return
}

func (o MySqlServiceCommandType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *MySqlServiceCommandType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := MySqlServiceCommandTypes().ParseMySqlServiceCommandType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid MySqlServiceCommandType %q", lit)
        }
    }
    return
}

func (o *MySqlServiceCommandType) IsCreate() bool {
    return o == _mySqlServiceCommandTypes.Create()
}

func (o *MySqlServiceCommandType) IsDelete() bool {
    return o == _mySqlServiceCommandTypes.Delete()
}

func (o *MySqlServiceCommandType) IsUpdate() bool {
    return o == _mySqlServiceCommandTypes.Update()
}

type mySqlServiceCommandTypes struct {
	values []*MySqlServiceCommandType
    literals []enum.Literal
}

var _mySqlServiceCommandTypes = &mySqlServiceCommandTypes{values: []*MySqlServiceCommandType{
    {name: "Create", ordinal: 0},
    {name: "Delete", ordinal: 1},
    {name: "Update", ordinal: 2}},
}

func MySqlServiceCommandTypes() *mySqlServiceCommandTypes {
	return _mySqlServiceCommandTypes
}

func (o *mySqlServiceCommandTypes) Values() []*MySqlServiceCommandType {
	return o.values
}

func (o *mySqlServiceCommandTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *mySqlServiceCommandTypes) Create() *MySqlServiceCommandType {
    return _mySqlServiceCommandTypes.values[0]
}

func (o *mySqlServiceCommandTypes) Delete() *MySqlServiceCommandType {
    return _mySqlServiceCommandTypes.values[1]
}

func (o *mySqlServiceCommandTypes) Update() *MySqlServiceCommandType {
    return _mySqlServiceCommandTypes.values[2]
}

func (o *mySqlServiceCommandTypes) ParseMySqlServiceCommandType(name string) (ret *MySqlServiceCommandType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*MySqlServiceCommandType), ok
	}
	return
}



