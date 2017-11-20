package elastic

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
func (o *Create) AggregateType() eventhorizon.AggregateType  { return ElasticServiceAggregateType }
func (o *Create) CommandType() eventhorizon.CommandType      { return CreateCommand }



        
type Delete struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Delete) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Delete) AggregateType() eventhorizon.AggregateType  { return ElasticServiceAggregateType }
func (o *Delete) CommandType() eventhorizon.CommandType      { return DeleteCommand }



        
type Update struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Update) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Update) AggregateType() eventhorizon.AggregateType  { return ElasticServiceAggregateType }
func (o *Update) CommandType() eventhorizon.CommandType      { return UpdateCommand }





type ElasticServiceCommandType struct {
	name  string
	ordinal int
}

func (o *ElasticServiceCommandType) Name() string {
    return o.name
}

func (o *ElasticServiceCommandType) Ordinal() int {
    return o.ordinal
}

func (o ElasticServiceCommandType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *ElasticServiceCommandType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := ElasticServiceCommandTypes().ParseElasticServiceCommandType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ElasticServiceCommandType %q", lit.Name)
        }
	}
	return
}

func (o ElasticServiceCommandType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *ElasticServiceCommandType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := ElasticServiceCommandTypes().ParseElasticServiceCommandType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid ElasticServiceCommandType %q", lit)
        }
    }
    return
}

func (o *ElasticServiceCommandType) IsCreate() bool {
    return o == _elasticServiceCommandTypes.Create()
}

func (o *ElasticServiceCommandType) IsDelete() bool {
    return o == _elasticServiceCommandTypes.Delete()
}

func (o *ElasticServiceCommandType) IsUpdate() bool {
    return o == _elasticServiceCommandTypes.Update()
}

type elasticServiceCommandTypes struct {
	values []*ElasticServiceCommandType
    literals []enum.Literal
}

var _elasticServiceCommandTypes = &elasticServiceCommandTypes{values: []*ElasticServiceCommandType{
    {name: "Create", ordinal: 0},
    {name: "Delete", ordinal: 1},
    {name: "Update", ordinal: 2}},
}

func ElasticServiceCommandTypes() *elasticServiceCommandTypes {
	return _elasticServiceCommandTypes
}

func (o *elasticServiceCommandTypes) Values() []*ElasticServiceCommandType {
	return o.values
}

func (o *elasticServiceCommandTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *elasticServiceCommandTypes) Create() *ElasticServiceCommandType {
    return _elasticServiceCommandTypes.values[0]
}

func (o *elasticServiceCommandTypes) Delete() *ElasticServiceCommandType {
    return _elasticServiceCommandTypes.values[1]
}

func (o *elasticServiceCommandTypes) Update() *ElasticServiceCommandType {
    return _elasticServiceCommandTypes.values[2]
}

func (o *elasticServiceCommandTypes) ParseElasticServiceCommandType(name string) (ret *ElasticServiceCommandType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*ElasticServiceCommandType), ok
	}
	return
}



