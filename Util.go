package main

import (
	"github.com/labstack/echo"
	"database/sql"
	"time"
	"strconv"
	"fmt"
	"github.com/fatih/color"
)

func Timestamp() int64 {
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

func PrintError(str string) {
	color.Red(str)
}

func PrintDebug(str string) {
	if debugging {
		color.Yellow(str)
	}
}

func PrintSuccess(str string) {
	if debugging {
		color.Green(str)
	}
}

func ValidateRequest(request *APIRequest, required []string, disallowed ...string) (bool, []string) {
	var result []string
	for _, r := range required {
		if _, ok := request.ApiArguments[r]; !ok {
			result = append(result, "missing " + r)
		}
	}

	for _, r := range disallowed {
		if _, ok := request.ApiArguments[r]; ok {
			result = append(result, r + " not allowed")
		}
	}

	if len(result) > 0 {
		PrintError(fmt.Sprintf("Missing required api arguments. %+v", result))
	}

	return len(result) == 0, result
}

func Insert(qs string, args map[string]interface{}) (sql.Result, error) {
	return pool.NamedExec(qs, args)
}

func SelectOne(dest interface{}, qs string, args map[string]interface{}) error {
	nstmt, err := pool.PrepareNamed(qs)
	if err != nil {
		return err
	}

	return nstmt.Get(dest, args)
}

func SelectAll(dest interface{}, qs string, args map[string]interface{}) error {
	nstmt, err := pool.PrepareNamed(qs)
	if err != nil {
		return err
	}

	return nstmt.Select(dest, args)
}

func Delete(qs string, args map[string]interface{}) (sql.Result, error) {
	return pool.NamedExec(qs, args)
}

func Update(qs string, args map[string]interface{}) (sql.Result, error) {
	return pool.NamedExec(qs, args)
}

func APIFail(reason interface{}) *APIResponse {
	return &APIResponse{ Success: false, Result: reason }
}

func APISuccess(reason interface{}) *APIResponse {
	return &APIResponse{ Success: true, Result: reason }
}