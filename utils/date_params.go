package utils

import (
	"encoding/json"
	"net/url"
	"time"
)

type PayloadDateParams time.Time
type QueryDateParams time.Time

const timeFormat = "2006-01-02"

func (qdp QueryDateParams) String() string {
	return time.Time(qdp).Format(timeFormat)
}

func (qdp QueryDateParams) EncodeValues(key string, v *url.Values) error {
	v.Add(key, qdp.String())
	return nil
}

func (dp PayloadDateParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dp).Format(timeFormat))
}

func (dp *PayloadDateParams) UnmarshalJSON(data []byte) error {
	str := string(data)
	t, err := time.Parse(`"`+timeFormat+`"`, str)
	if err != nil {
		return err
	}
	*dp = PayloadDateParams(t)
	return nil
}

type RequiredDateParams struct {
	From QueryDateParams `validate:"required" url:"from,omitempty"`
	To   QueryDateParams `validate:"required" url:"to,omitempty"`
}

type RequiredDatePayload struct {
	From PayloadDateParams `validate:"required" json:"from,omitempty"`
	To   PayloadDateParams `validate:"required" json:"to,omitempty"`
}

type OptionalDateParams struct {
	From QueryDateParams `validate:"omitempty" url:"from,omitempty"`
	To   QueryDateParams `validate:"omitempty" url:"to,omitempty"`
}

type OptionalDatePayload struct {
	From PayloadDateParams `validate:"omitempty" json:"from,omitempty"`
	To   PayloadDateParams `validate:"omitempty" json:"to,omitempty"`
}
