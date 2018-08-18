package main

import (
	"time"
	"github.com/labstack/echo"
	"strconv"
)

func timestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func ParamToInt64(c echo.Context, param string) (int64, bool) {
	if str := c.Param(param); str != "" {
		p, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return p, true
		}
	}

	return 0, false
}