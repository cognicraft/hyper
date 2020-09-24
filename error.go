package hyper

// Error .
type Error struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	Message     string `json:"message"`
	Code        string `json:"code,omitempty"`
}

// Errors .
type Errors []Error

// Find
func (es Errors) Find(accept func(e Error) bool) (Error, bool) {
	for _, e := range es {
		if accept(e) {
			return e, true
		}
	}
	return Error{}, false
}

// Filter
func (es Errors) Filter(accept func(e Error) bool) Errors {
	var res Errors
	for _, e := range es {
		if accept(e) {
			res = append(res, e)
		}
	}
	return res
}

func ErrorItem(errs ...error) Item {
	res := Item{}
	for _, err := range errs {
		e := Error{Message: err.Error()}
		if errC, ok := err.(errorCoder); ok {
			e.Code = errC.Code()
		}
		res.Errors = append(res.Errors, e)
	}
	return res
}

type errorCoder interface {
	Code() string
}
