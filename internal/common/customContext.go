package common

import (
	"bytes"
	"io"

	"pr-reviewer/internal/modules/log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const REQUEST_ID_KEY = "request_id_key"
const LOGGER = "logger"
const IS_LOGGER_SET = "is_logger_set"

type Context struct {
	echo.Context
}

func NewContext(c echo.Context) *Context {
	cc := &Context{c}
	cc.Set(REQUEST_ID_KEY, uuid.New().String())
	cc.Set(LOGGER, log.Logger())
	return cc
}

func Convert(c echo.Context) *Context {
	cc := c.(*Context)
	if cc.Get(IS_LOGGER_SET) == nil || !cc.Get(IS_LOGGER_SET).(bool) {
		cc.Logrus().Init()
		cc.Set(IS_LOGGER_SET, true)
	}
	return cc
}

func (c *Context) Logrus() *log.Logrus {
	return c.Get(LOGGER).(*log.Logrus)
}

func (c *Context) RequestID() string {
	return c.Get(REQUEST_ID_KEY).(string)
}

// we need to read get data from body twice, once in that endpoint and another time in function depend onm signup type
// since Body is a reader we can't read twice from the same reader so we read to the byte[], and create new reader and reset body twice
func (c *Context) BindAndReset(dest interface{}) error {
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request().Body)
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := c.Bind(dest); err != nil {
		return err
	}

	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}
