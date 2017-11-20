package process

import (
    "eye/shared"
    "github.com/looplab/eventhorizon"
)
        
type ProcessService struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
    shared.Service
}

func NewProcessService() (ret *ProcessService) {
    service := shared.NewService()
    ret = &ProcessService{
        Service: service,
    }
    return
}
func (o *ProcessService) EntityID() eventhorizon.UUID { return o.Id }










