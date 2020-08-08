package common

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"
)

//convert sql data to struct
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		//get sql values
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		//get field name
		name := objValue.Type().Field(i).Name
		//get field type
		structFieldType := objValue.Field(i).Type()
		//get variable type
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			//convert types
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {
				log.Fatal(err)
			}
		}
		//set values
		objValue.FieldByName(name).Set(val)
	}
}

//convert type
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//can add more type conversions

	return reflect.ValueOf(value), errors.New("unknown type：" + ntype)
}
