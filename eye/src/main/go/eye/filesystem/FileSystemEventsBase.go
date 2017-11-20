package filesystem

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




type FileSystemServiceEventType struct {
	name  string
	ordinal int
}

func (o *FileSystemServiceEventType) Name() string {
    return o.name
}

func (o *FileSystemServiceEventType) Ordinal() int {
    return o.ordinal
}

func (o FileSystemServiceEventType) MarshalJSON() (ret []byte, err error) {
	return json.Marshal(&enum.EnumBaseJson{Name: o.name})
}

func (o *FileSystemServiceEventType) UnmarshalJSON(data []byte) (err error) {
	lit := enum.EnumBaseJson{}
	if err = json.Unmarshal(data, &lit); err == nil {
		if v, ok := FileSystemServiceEventTypes().ParseFileSystemServiceEventType(lit.Name); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid FileSystemServiceEventType %q", lit.Name)
        }
	}
	return
}

func (o FileSystemServiceEventType) GetBSON() (ret interface{}, err error) {
	return o.name, nil
}

func (o *FileSystemServiceEventType) SetBSON(raw bson.Raw) (err error) {
	var lit string
    if err = raw.Unmarshal(&lit); err == nil {
		if v, ok := FileSystemServiceEventTypes().ParseFileSystemServiceEventType(lit); ok {
            *o = *v
        } else {
            err = fmt.Errorf("invalid FileSystemServiceEventType %q", lit)
        }
    }
    return
}

func (o *FileSystemServiceEventType) IsCreate() bool {
    return o == _fileSystemServiceEventTypes.Create()
}

func (o *FileSystemServiceEventType) IsCreated() bool {
    return o == _fileSystemServiceEventTypes.Created()
}

func (o *FileSystemServiceEventType) IsDelete() bool {
    return o == _fileSystemServiceEventTypes.Delete()
}

func (o *FileSystemServiceEventType) IsDeleted() bool {
    return o == _fileSystemServiceEventTypes.Deleted()
}

func (o *FileSystemServiceEventType) IsUpdate() bool {
    return o == _fileSystemServiceEventTypes.Update()
}

func (o *FileSystemServiceEventType) IsUpdated() bool {
    return o == _fileSystemServiceEventTypes.Updated()
}

type fileSystemServiceEventTypes struct {
	values []*FileSystemServiceEventType
    literals []enum.Literal
}

var _fileSystemServiceEventTypes = &fileSystemServiceEventTypes{values: []*FileSystemServiceEventType{
    {name: "Create", ordinal: 0},
    {name: "Created", ordinal: 1},
    {name: "Delete", ordinal: 2},
    {name: "Deleted", ordinal: 3},
    {name: "Update", ordinal: 4},
    {name: "Updated", ordinal: 5}},
}

func FileSystemServiceEventTypes() *fileSystemServiceEventTypes {
	return _fileSystemServiceEventTypes
}

func (o *fileSystemServiceEventTypes) Values() []*FileSystemServiceEventType {
	return o.values
}

func (o *fileSystemServiceEventTypes) Literals() []enum.Literal {
	if o.literals == nil {
		o.literals = make([]enum.Literal, len(o.values))
		for i, item := range o.values {
			o.literals[i] = item
		}
	}
	return o.literals
}

func (o *fileSystemServiceEventTypes) Create() *FileSystemServiceEventType {
    return _fileSystemServiceEventTypes.values[0]
}

func (o *fileSystemServiceEventTypes) Created() *FileSystemServiceEventType {
    return _fileSystemServiceEventTypes.values[1]
}

func (o *fileSystemServiceEventTypes) Delete() *FileSystemServiceEventType {
    return _fileSystemServiceEventTypes.values[2]
}

func (o *fileSystemServiceEventTypes) Deleted() *FileSystemServiceEventType {
    return _fileSystemServiceEventTypes.values[3]
}

func (o *fileSystemServiceEventTypes) Update() *FileSystemServiceEventType {
    return _fileSystemServiceEventTypes.values[4]
}

func (o *fileSystemServiceEventTypes) Updated() *FileSystemServiceEventType {
    return _fileSystemServiceEventTypes.values[5]
}

func (o *fileSystemServiceEventTypes) ParseFileSystemServiceEventType(name string) (ret *FileSystemServiceEventType, ok bool) {
	if item, ok := enum.Parse(name, o.Literals()); ok {
		return item.(*FileSystemServiceEventType), ok
	}
	return
}



