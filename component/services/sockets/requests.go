package sockets

type Request interface {
	GetCommand() string
	GetData() interface{}
	GetParams() interface{}
}

type JSONRequest struct {
	Command string
	Data map[string]interface{}
	Params map[string]interface{}
}

type CommandRequest struct {
	Command string
	Data string
}

func (R JSONRequest) New() *JSONRequest {
	r := new(JSONRequest)
	r.Data = make(map[string]interface{})
	r.Params = make(map[string]interface{})
	return r
}

func (R *JSONRequest) GetCommand() string {
	return R.Command
}

func (R *JSONRequest) GetData() interface{} {
	return R.Data
}

func (R *JSONRequest) GetParams() interface{} {
	return R.Params
}

func (R *CommandRequest) GetCommand() string {
	return R.Command
}

func (R *CommandRequest) GetData() interface{} {
	return R.Data
}

func (R *CommandRequest) GetParams() interface{} {
	return nil
}
