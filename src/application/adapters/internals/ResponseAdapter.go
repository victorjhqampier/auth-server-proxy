package internals

type ResponseAdapter struct {
	StatusCode int                 `json:"-"`
	Errors     []FieldErrorAdapter `json:"errors,omitempty"`
	Response   interface{}         `json:"response,omitempty"` //`json:"data,omitempty"`
}
