package render

import (
	"west.garden/template/common/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RespJsonData struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

var errMap = map[int]string{
	ErrParams:    "There's some error with your information. Please check and try agian.",
	ServiceError: "There's some error with your information. Please check and try agian.",
}

type TipError struct {
	Err error
}

func NewTipErr(err string) *TipError {
	return &TipError{
		Err: fmt.Errorf(err),
	}
}

func (t *TipError) Error() string {
	return t.Err.Error()
}

func Json(c *gin.Context, code int, data interface{}) {
	msg := ""
	if code != 0 {
		t, ok := data.(*TipError)
		if ok {
			msg = t.Err.Error()
		} else {
			if val, ok := errMap[code]; ok {
				log.Log.Error("//discover error//:", data)
				msg = val
			} else {
				msg = fmt.Sprintf("%v", data)
			}
		}
		data = nil
	} else {
		msg = "success"
	}
	result := &RespJsonData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, result)
}

func AbortJson(c *gin.Context, code int, data interface{}) {
	msg := http.StatusText(code)
	if code != http.StatusOK && code != Ok {
		if data != nil {
			msg = fmt.Sprintf("%v", data)
		} else {
			data = msg
		}
	}
	result := &RespJsonData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.AbortWithStatusJSON(code, result)
}

type PageResult struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}
