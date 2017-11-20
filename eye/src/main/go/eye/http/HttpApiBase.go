package http

import (
    "eye/shared"
    "github.com/looplab/eventhorizon"
)
        
type HttpService struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
    shared.Service
}

func NewHttpService() (ret *HttpService) {
    service := shared.NewService()
    ret = &HttpService{
        Service: service,
    }
    return
}
func (o *HttpService) EntityID() eventhorizon.UUID { return o.Id }










