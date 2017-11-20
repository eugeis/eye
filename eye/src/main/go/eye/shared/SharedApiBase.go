package shared




type QueryResult struct {
    Info string `json:"info" eh:"optional"`
}

func NewQueryResult() (ret *QueryResult) {
    ret = &QueryResult{}
    return
}




type CommandRequest struct {
}

func NewCommandRequest() (ret *CommandRequest) {
    ret = &CommandRequest{}
    return
}


type ExportRequest struct {
    Query string `json:"query" eh:"optional"`
    EvalExpr string `json:"evalExpr" eh:"optional"`
}

func NewExportRequest() (ret *ExportRequest) {
    ret = &ExportRequest{}
    return
}


type ValidationRequest struct {
    RegExpr string `json:"regExpr" eh:"optional"`
    All bool `json:"all" eh:"optional"`
    *ExportRequest
}

func NewValidationRequest() (ret *ValidationRequest) {
    exportRequest := NewExportRequest()
    ret = &ValidationRequest{
        ExportRequest: exportRequest,
    }
    return
}





