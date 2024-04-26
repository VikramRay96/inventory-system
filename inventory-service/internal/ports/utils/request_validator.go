package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"inventory-system/common/pkg/dto"
	ConfigurationServiceDto "inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	"reflect"
	"strings"
)

type IRequestValidator interface {
	ValidateCreateConfigurationRequest(requestBody ConfigurationServiceDto.CreateNewConfigurationRequestBody) (*dto.ErrorResponseDto, bson.M)
	ValidationErrors(err error) *dto.ErrorResponseDto
	ValidateStruct(interfaceData interface{}) *dto.ErrorResponseDto
}
type RequestValidator struct {
	Validator *validator.Validate
}

func NewRequestValidator() *RequestValidator {
	validatorObject := validator.New()
	validatorObject.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	requestValidator := &RequestValidator{
		Validator: validatorObject,
	}
	return requestValidator
}

func (rv RequestValidator) ValidateCreateConfigurationRequest(requestBody ConfigurationServiceDto.CreateNewConfigurationRequestBody) (*dto.ErrorResponseDto, bson.M) {
	err := rv.Validator.Struct(&requestBody)
	if err != nil {
		return rv.ValidationErrors(err), nil
	}
	return nil, nil
}

func (rv RequestValidator) ValidationErrors(err error) *dto.ErrorResponseDto {
	var errorList string
	var validationErrors validator.ValidationErrors
	var errDto dto.ErrorResponseDto

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			errorList = errorList + fieldError.Field() + " : " + ValidationErrorMsg(fieldError) + ", "
		}
		if errorList != "" {
			errorList = strings.TrimRight(errorList, ", ")
		}
	}
	if err != nil {
		errDto.SetError(status_code.IMS400)
		returnValue := dto.ErrorResponseDto{
			StatusCode: errDto.StatusCode,
			Message:    errorList,
		}
		return &returnValue
	}
	return nil
}
func (rv RequestValidator) ValidateStruct(interfaceData interface{}) *dto.ErrorResponseDto {
	return rv.ValidationErrors(rv.Validator.Struct(interfaceData))
}

func ValidationErrorMsg(f validator.FieldError) string {
	switch f.Tag() {
	case "required":
		return "Field is missing"

	}
	return f.Error()
}
