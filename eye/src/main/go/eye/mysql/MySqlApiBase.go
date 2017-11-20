package mysql

import (
    "eye/shared"
    "github.com/looplab/eventhorizon"
)
        
type MySqlService struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
    shared.Service
}

func NewMySqlService() (ret *MySqlService) {
    service := shared.NewService()
    ret = &MySqlService{
        Service: service,
    }
    return
}
func (o *MySqlService) EntityID() eventhorizon.UUID { return o.Id }










