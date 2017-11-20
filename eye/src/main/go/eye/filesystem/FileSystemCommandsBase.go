package filesystem

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
func (o *Create) AggregateType() eventhorizon.AggregateType  { return FileSystemServiceAggregateType }
func (o *Create) CommandType() eventhorizon.CommandType      { return CreateCommand }



        
type Delete struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Delete) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Delete) AggregateType() eventhorizon.AggregateType  { return FileSystemServiceAggregateType }
func (o *Delete) CommandType() eventhorizon.CommandType      { return DeleteCommand }



        
type Update struct {
    Name string `json:"name" eh:"optional"`
    Id eventhorizon.UUID `json:"id" eh:"optional"`
}
func (o *Update) AggregateID() eventhorizon.UUID            { return o.Id }
func (o *Update) AggregateType() eventhorizon.AggregateType  { return FileSystemServiceAggregateType }
func (o *Update) CommandType() eventhorizon.CommandType      { return UpdateCommand }





type FileSystemServiceCommandType struct {
	name  string
	ordinal int
}

func (o *FileSystemServiceCommandType) Name() string {
    return o.name
}

func (o *FileSystemServiceCommandType) Ordinal() int {
    return o.ordinal
}

func (o FileSystemServiceCommandType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *FileSystemServiceCommandType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := FileSystemServiceCommandTypes().ParseFileSystemServiceCommandType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid FileSystemServiceCommandType %q", lit.Name)
        }
	}
	return
}

func (o FileSystemServiceCommandType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *FileSystemServiceCommandType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := FileSystemServiceCommandTypes().ParseFileSystemServiceCommandType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid FileSystemServiceCommandType %q", lit)
        }
    }
    return
}

func (o *FileSystemServiceCommandType) IsCreate() bool {
    return o == _fileSystemServiceCommandTypes.Create()
}

func (o *FileSystemServiceCommandType) IsDelete() bool {
    return o == _fileSystemServiceCommandTypes.Delete()
}

func (o *FileSystemServiceCommandType) IsUpdate() bool {
    return o == _fileSystemServiceCommandTypes.Update()
}

type fileSystemServiceCommandTypes struct {
	values []*FileSystemServiceCommandType
    literals []enum.Literal
}

var _fileSystemServiceCommandTypes = &fileSystemServiceCommandTypes{values: []*FileSystemServiceCommandType{
    {name: "Create", ordinal: 0},
    {name: "Delete", ordinal: 1},
    {name: "Update", ordinal: 2}},
}

func FileSystemServiceCommandTypes() *fileSystemServiceCommandTypes {
	return _fileSystemServiceCommandTypes
}

func (o *fileSystemServiceCommandTypes) Values() []*FileSystemServiceCommandType {
	return o.values
}

func (o *fileSystemServiceCommandTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *fileSystemServiceCommandTypes) Create() *FileSystemServiceCommandType {
    return _fileSystemServiceCommandTypes.values[0]
}

func (o *fileSystemServiceCommandTypes) Delete() *FileSystemServiceCommandType {
    return _fileSystemServiceCommandTypes.values[1]
}

func (o *fileSystemServiceCommandTypes) Update() *FileSystemServiceCommandType {
    return _fileSystemServiceCommandTypes.values[2]
}

func (o *fileSystemServiceCommandTypes) ParseFileSystemServiceCommandType(name string) (ret *FileSystemServiceCommandType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*FileSystemServiceCommandType), ok
	}
	return
}



