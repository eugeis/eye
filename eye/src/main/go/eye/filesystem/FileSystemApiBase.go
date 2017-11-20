package filesystem

import (
    "eye/shared"
    "github.com/looplab/eventhorizon"
)
        
type FileSystemService struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
    shared.Service
}

func NewFileSystemService() (ret *FileSystemService) {
    service := shared.NewService()
    ret = &FileSystemService{
        Service: service,
    }
    return
}
func (o *FileSystemService) EntityID() eventhorizon.UUID { return o.Id }










