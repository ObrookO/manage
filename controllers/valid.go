package controllers

import (
	"errors"
	"fmt"
	"manage/models"
	"reflect"

	"github.com/astaxie/beego/validation"
)

// validData 表单验证
func validData(data interface{}) error {
	v := validation.Validation{}
	var b bool
	var err error

	switch reflect.TypeOf(data).String() {
	case "models.Article":
		b, err = v.Valid(data.(models.Article))

	case "models.Manager":
		b, err = v.Valid(data.(models.Manager))

	default:
		return errors.New("unrecognized type of data")
	}

	if err != nil {
		return err
	}

	if !b {
		return errors.New(fmt.Sprintf("%s %s", v.Errors[0].Field, v.Errors[0].Message))
	}

	return nil
}
