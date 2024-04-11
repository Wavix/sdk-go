package utils

type RequiredDateParams struct {
	From string `validate:"required,datetime=2006-01-02" url:"from,omitempty"`
	To   string `validate:"required,datetime=2006-01-02" url:"to,omitempty"`
}

type RequiredDatePayload struct {
	From string `validate:"required,datetime=2006-01-02" json:"from,omitempty"`
	To   string `validate:"required,datetime=2006-01-02" json:"to,omitempty"`
}

type OptionalDateParams struct {
	From string `validate:"omitempty,datetime=2006-01-02" url:"from,omitempty"`
	To   string `validate:"omitempty,datetime=2006-01-02" url:"to,omitempty"`
}

type OptionalDatePayload struct {
	From string `validate:"omitempty,datetime=2006-01-02" json:"from,omitempty"`
	To   string `validate:"omitempty,datetime=2006-01-02" json:"to,omitempty"`
}
