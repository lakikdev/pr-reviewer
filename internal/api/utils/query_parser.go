package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type parser struct {
	Query url.Values
}

func NewParser(query url.Values) *parser {
	return &parser{query}
}

//TimeParam - get time value from request query
func (p *parser) TimeParam(name string) (time.Time, error) {
	//return NOW as default time
	t := time.Now()
	value := p.Query.Get(name)
	if value == "" {
		return t, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return t, err
	}

	return parsed, nil
}

//GetLimitOffset
func (p *parser) GetLimitOffset() (int64, int64, error) {
	limit, err := p.IntParam("limit")
	if err != nil {
		return 0, 0, err
	} else if limit == 0 {
		limit = 50
	}
	offset, err := p.IntParam("offset")
	if err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}

//IntParam
func (p *parser) IntParam(name string) (int64, error) {
	value := p.Query.Get(name)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing parameter: %s", name)
	}
	return parsed, nil
}

//FloatParam
func (p *parser) FloatParam(name string) (float64, error) {
	value := p.Query.Get(name)
	if value == "" {
		return 0, fmt.Errorf("missing parameter: %s", name)
	}
	return strconv.ParseFloat(value, 64)
}

//StringParam
func (p *parser) StringParam(name string) (string, error) {
	return p.Query.Get(name), nil
}

//BoolParam
func (p *parser) BoolParam(name string) (bool, error) {
	return strconv.ParseBool(p.Query.Get(name))
}
