package util

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"reflect"
)

func Base2String(data interface{}) (string,error) {
	var str string
	if data == nil {
		return "",fmt.Errorf("data is nil")
	}
	switch data.(type) {
	case int,int64,int32 :
		str = strconv.Itoa(data.(int))
	case float64:
		str = strconv.FormatFloat(data.(float64),'f',-1,64)
	case float32:
		str = strconv.FormatFloat(float64(data.(float32)),'f',-1,64)
	case string:
		str = data.(string)
	default:
		return "",fmt.Errorf("data cannot convert")
	}
	return str,nil
}

func Md5(s interface{})(string,error){
	var str string
	str, err := Base2String(s)
	if err != nil {
		return "",err
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr),nil
}

func Base2String2(data interface{}) (string,error) {
	if data == nil {
		return "",fmt.Errorf("data is nil")
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.String:
		return data.(string),nil
	case reflect.Float32,reflect.Float64:
		str := strconv.FormatFloat(data.(float64),'f',-1,64)
		return str,nil
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64 :
		return strconv.Itoa(data.(int)),nil
	default:
		return "",fmt.Errorf("data cannot convert")
	}
}
