package main

import (
	"github.com/labstack/echo"
	"database/sql"
	"time"
	"strconv"
	"fmt"
	"github.com/fatih/color"
	"reflect"
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

func ValidateRequest(request *APIRequest, required []string, disallowed []string) (bool, []string) {
	var result []string

	if required != nil {
		for _, r := range required {
			if _, ok := request.ApiArguments[r]; !ok {
				result = append(result, "missing " + r)
			}
		}
	}

	if disallowed != nil {
		for _, r := range disallowed {
			if _, ok := request.ApiArguments[r]; ok {
				result = append(result, r + " not allowed")
			}
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

func ValidateAndSelectAll(request *APIRequest, dest interface{}, qs string, required []string, disallowed ...string) *APIResponse {
	if ok, err := ValidateRequest(request, required, disallowed); !ok {
		return APIFail(err)
	}

	if err := SelectAll(dest, qs, request.ApiArguments); err != nil {
		PrintError(fmt.Sprintf("Error trying to select all. %+v", err))
		return APIFail(SelectFail)
	}

	return APISuccess(dest)
}

func ValidateAndSelectOne(request *APIRequest, dest interface{}, qs string, required []string, disallowed ...string) *APIResponse {
	if ok, err := ValidateRequest(request, required, disallowed); !ok {
		return APIFail(err)
	}

	if err := SelectOne(dest, qs, request.ApiArguments); err != nil {
		PrintError(fmt.Sprintf("Error trying to select one. %+v", err))
		return APIFail(SelectFail)
	}

	return APISuccess(dest)
}

func ValidateAndDelete(request *APIRequest, qs string, required []string, disallowed ...string) *APIResponse {
	if ok, err := ValidateRequest(request, required, disallowed); !ok {
		return APIFail(err)
	}

	result, err := Delete(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Error trying to delete. %+v", err))
		return APIFail(DeleteFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return APIFail(NoRowsAffected)
	}

	return APISuccess(rowsAffected)
}

func ValidateAndInsert(request *APIRequest, qs string, required []string, disallowed ...string) *APIResponse {
	if ok, err := ValidateRequest(request, required, disallowed); !ok {
		return APIFail(err)
	}

	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Error trying to insert. %+v", err))
		return APIFail(InsertFail)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return APIFail(NoInsertID)
	}

	return APISuccess(lastInsertId)
}

func ValidateAndUpdateFromStruct(request *APIRequest, dest interface{}, table string, qs string, qs2 string, required []string, disallowed ...string) *APIResponse {
	if ok, err := ValidateRequest(request, required, disallowed); !ok {
		return APIFail(err)
	}

	keys := ``
	elem := reflect.ValueOf(dest).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		if tag, ok := field.Tag.Lookup("static"); ok && tag == "true" {
			continue
		}

		key := field.Name
		if _, ok := request.ApiArguments[key]; ok {
			if len(keys) != 0 {
				keys += `,`
			}
			keys += ` ` + table + `.` + key + ` = :` + key
		}
	}

	result, err := Update(qs + keys + qs2, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Error trying to update from struct. %+v", err))
		return APIFail(UpdateFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return APIFail(NoRowsAffected)
	}

	return APISuccess(rowsAffected)
}

func ValidateAndUpdate(request *APIRequest, qs string, required []string, disallowed ...string) *APIResponse {
	if ok, err := ValidateRequest(request, required, disallowed); !ok {
		return APIFail(err)
	}

	result, err := Update(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Error trying to update. %+v", err))
		return APIFail(UpdateFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return APIFail(NoRowsAffected)
	}

	return APISuccess(rowsAffected)
}