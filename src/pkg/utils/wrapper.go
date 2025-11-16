package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpError "order-service/src/pkg/http-error"
	"order-service/src/pkg/log"

	"github.com/gofiber/fiber/v2"
)

type Result struct {
	Data  interface{}
	Error interface{}
}

type BaseWrapperModel struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Meta struct {
	Method        string    `json:"method"`
	Url           string    `json:"url"`
	Code          string    `json:"code"`
	ContentLength int       `json:"content_length"`
	Date          time.Time `json:"date"`
	Ip            string    `json:"ip"`
}

func Response(data interface{}, message string, code int, c *fiber.Ctx) error {
	success := code < http.StatusBadRequest

	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Method(),
		Code:          fmt.Sprintf("%v", http.StatusOK),
		ContentLength: len(c.Body()),
		Ip:            c.IP(),
	}
	byteMeta, _ := json.Marshal(meta)

	log.GetLogger().Info("service-info", "Logging service...", "audit-log", string(byteMeta))

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	return c.Status(code).JSON(result)
}

func ResponseError(err interface{}, c *fiber.Ctx) error {
	errObj := getErrorStatusCode(err)

	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Method(),
		Code:          fmt.Sprintf("%v", errObj.ResponseCode),
		ContentLength: len(c.Body()),
		Ip:            c.IP(),
	}
	byteMeta, _ := json.Marshal(meta)

	log.GetLogger().Error("service-error", "Logging service...", "audit-log", string(byteMeta))

	result := BaseWrapperModel{
		Success: false,
		Data:    errObj.Data,
		Message: errObj.Message,
		Code:    errObj.Code,
	}

	return c.Status(errObj.ResponseCode).JSON(result)
}

func getErrorStatusCode(err interface{}) httpError.CommonErrorData {
	errData := httpError.CommonErrorData{}

	switch obj := err.(type) {
	case httpError.BadRequestData:
		errData.ResponseCode = http.StatusBadRequest
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
	case httpError.UnauthorizedData:
		errData.ResponseCode = http.StatusUnauthorized
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
	case httpError.ConflictData:
		errData.ResponseCode = http.StatusConflict
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
	case httpError.NotFoundData:
		errData.ResponseCode = http.StatusNotFound
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
	case httpError.InternalServerErrorData:
		errData.ResponseCode = http.StatusInternalServerError
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
	default:
		errData.ResponseCode = http.StatusConflict
		errData.Code = http.StatusConflict
	}

	return errData
}
