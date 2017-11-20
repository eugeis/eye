package http

import (
    "context"
    "github.com/looplab/eventhorizon"
)
type QueryRepository struct {
    repo eventhorizon.ReadRepo `json:"repo" eh:"optional"`
    context context.Context `json:"context" eh:"optional"`
}

func NewHttpServiceQueryRepository(repo eventhorizon.ReadRepo, context context.Context) (ret *QueryRepository) {
    ret = &QueryRepository{
        repo: repo,
        context: context,
    }
    return
}

func (o *QueryRepository) FindAll() (ret []*HttpService, err error) {
    var result []eventhorizon.Entity
	if result, err = o.repo.FindAll(o.context); err == nil {
        ret = make([]*HttpService, len(result))
		for i, e := range result {
            ret[i] = e.(*HttpService)
		}
    }
        
    var result []eventhorizon.Entity
	if result, err = o.repo.FindAll(o.context); err == nil {
        ret = make([]*HttpService, len(result))
		for i, e := range result {
            ret[i] = e.(*HttpService)
		}
    }
    return
}

func (o *QueryRepository) FindById(id eventhorizon.UUID) (ret *HttpService, err error) {
    var result eventhorizon.Entity
	if result, err = o.repo.Find(o.context, id); err == nil {
        ret = result.(*HttpService)
    }
        
    var result eventhorizon.Entity
	if result, err = o.repo.Find(o.context, id); err == nil {
        ret = result.(*HttpService)
    }
    return
}

func (o *QueryRepository) CountAll() (ret int, err error) {
    var result []*HttpService
	if result, err = o.FindAll(); err == nil {
        ret = len(result)
    }
        
    var result []*HttpService
	if result, err = o.FindAll(); err == nil {
        ret = len(result)
    }
    return
}

func (o *QueryRepository) CountById(id eventhorizon.UUID) (ret int, err error) {
    var result *HttpService
	if result, err = o.FindById(id); err == nil && result != nil {
        ret = 1
    }
        
    var result *HttpService
	if result, err = o.FindById(id); err == nil && result != nil {
        ret = 1
    }
    return
}

func (o *QueryRepository) ExistAll() (ret bool, err error) {
    var result int
	if result, err = o.CountAll(); err == nil {
        ret = result > 0
    }
        
    var result int
	if result, err = o.CountAll(); err == nil {
        ret = result > 0
    }
    return
}

func (o *QueryRepository) ExistById(id eventhorizon.UUID) (ret bool, err error) {
    var result int
	if result, err = o.CountById(id); err == nil {
        ret = result > 0
    }
        
    var result int
	if result, err = o.CountById(id); err == nil {
        ret = result > 0
    }
    return
}









