package middlewares

import (
	"fmt"
	"strconv"
	"time"

	"pr-reviewer/internal/common"

	"github.com/labstack/echo/v4"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*common.Context)

			req := c.Request()
			res := c.Response()
			start := time.Now()

			err := next(c)
			if err != nil {
				HTTPErrorHandler(err, cc)
			}

			stop := time.Now()

			p := req.URL.Path
			// skip /metric endpoint from logging to avoid logs every 5 seconds that endpoint for prometheus
			if p == "/metrics" {
				return nil
			}

			bytesIn := req.Header.Get(echo.HeaderContentLength)
			bytesOut := strconv.FormatInt(res.Size, 10)

			// check if we have userID in logrus entry data
			// if we have it, we need to remove it from logrus entry data
			// and add it to logrus entry fields as UserID (to make sure it will be in the beginning of the log)
			// logrus prints the fields in alphabetical order
			// first special characters, then numbers, then capital letters, then small letters
			userID := cc.Logrus().Entry.Data["userID"]
			if userID != nil {
				cc.Logrus().WithField("UserID", userID)
				cc.Logrus().Entry.Data["userID"] = nil
			}

			cc.Logrus().WithFields(map[string]interface{}{
				"uri":           req.RequestURI,
				"method":        req.Method,
				"path":          p,
				"user_agent":    req.UserAgent(),
				"status":        res.Status,
				"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"latency_human": stop.Sub(start).String(),
				"bytes_in":      bytesIn,
				"bytes_out":     bytesOut,
				"withError":     err != nil,
			})

			if err != nil {
				cc.Logrus().Error(err.Error())
			} else {
				funcName := cc.Logrus().Entry.Data["func"]
				successMsg := "request success"
				if funcName != nil {
					successMsg = fmt.Sprintf("%s - Request Success", funcName.(string))
				}
				cc.Logrus().Info(successMsg)
			}

			return err
		}
	}
}
