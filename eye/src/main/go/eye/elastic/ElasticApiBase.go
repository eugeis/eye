package elastic

import (
    "eye/shared"
    "github.com/looplab/eventhorizon"
)
        
type ElasticService struct {
    Id eventhorizon.UUID `json:"id" eh:"optional"`
    shared.Service
}

func NewElasticService() (ret *ElasticService) {
    service := shared.NewService()
    ret = &ElasticService{
        Service: service,
    }
    return
}
func (o *ElasticService) EntityID() eventhorizon.UUID { return o.Id }










