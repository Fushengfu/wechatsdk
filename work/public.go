package work

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/**
 *  反射函数
 */
func (w *WorkWechat) callBack(method string, params ...interface{}) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			logs.Critical("callBack系统异常：", err)
		}
	}()

	mCall := reflect.ValueOf(w)

	var values []reflect.Value

	for _, v := range params {
		values = append(values, reflect.ValueOf(v))
	}

	time.Sleep(time.Millisecond * 1000)
	_method := mCall.MethodByName(method)
	invalidMethod := _method == reflect.Value{}

	logs.Critical("回调次数：", w.Limit, method, invalidMethod, "企业ID:", w.Appid, w)
	if invalidMethod {
		w.Limit = 0
		return "", errors.New("找不到回调方法：" + method)
	}

	result := _method.Call(values)

	logs.Critical("回调结果：", result[0], result[1])
	var err error
	if result[1].Interface() == nil {
		err = nil
	} else {
		err = result[1].Interface().(error)
	}

	w.Limit = 0
	return result[0].String(), err
}

/**
 *  发起网络请求
 */
func (w *WorkWechat) sendForm(method string, uri string, data map[string]interface{}, params ...interface{}) (string, error) {
	if len(w.Appid) > 0 && len(w.Appid) < 18 {
		w.Limit = 0
		return "", errors.New("非法APPID")
	}

	var str string
	var err error
	result := make(map[string]interface{})

	if method == "GET" {
		str, err = request.Get(uri, data)
	} else {
		str, err = request.Post(uri, data)
	}

	inputData, _ := json.Marshal(data)

	if err != nil {
		logs.Critical("Http请求", "请求方法："+method, "Url:"+uri, "请求参数：", string(inputData), "\r\n", "请求结果：", str, err)
		w.Limit = 0
		return str, err
	} else {
		err = json.Unmarshal([]byte(str), &result)
		if err != nil {
			w.Limit = 0
			return str, err
		}

		errcode, ok := result["errcode"]
		if ok && inArray(int64(errcode.(float64)), []int64{41001, 42001, 40014, 45033, 45009}) && w.Limit < 3 {
			logs.Critical("Http请求", "请求方法："+method, "Url:"+uri, "请求参数：", string(inputData), "\r\n", "请求结果：", str, err)

			pc, _, _, ok1 := runtime.Caller(1)
			logs.Critical("W:", w)
			if ok1 {
				w.Limit = w.Limit + 1

				if !inArray(int64(errcode.(float64)), []int64{45033, 45009}) {
					Rds.Del("WECHAT_QY::AUTH_CORPID_" + w.Appid)
					if len(w.SuiteId) == 0 {
						w.AccessToken, _ = w.getAccessToken()
					} else {
						_, _ = w.GetCorpToken()
					}
				} else {
					time.Sleep(time.Millisecond * 500)
					w.Limit = 3
				}

				f := runtime.FuncForPC(pc)
				actionArr := strings.Split(f.Name(), ".")
				if len(actionArr) > 0 {
					_action := actionArr[len(actionArr)-1]
					return w.callBack(_action, params...)
				}
			}
			w.Limit = 0
		}

		if ok && inArray(int64(errcode.(float64)), []int64{45033, 45009}) {
			logs.Critical("\n请求方法："+method, "\nUrl:"+uri, "\n请求参数：", string(inputData), "\n请求结果：", str, err)
		}
	}

	w.Limit = 0
	return str, err
}

/**
 *  判定在数组中
 */
func inArray(pattem int64, array []int64) bool {
	for _, va := range array {
		if pattem == va {
			return true
		}
	}
	return false
}

/**
 *  转字符串
 */
func ToString(data interface{}) string {
	if data == nil {
		return ""
	}

	switch reflect.TypeOf(data).String() {
	case "string":
		return data.(string)
	case "int":
		return strconv.Itoa(data.(int))
	case "float64":
		return strconv.Itoa(int(int64(data.(float64))))
	case "float32":
		return strconv.Itoa(int(int32(data.(float32))))
	default:
		bytes, _ := json.Marshal(data)
		return string(bytes)
	}
}

/**
 *  转字符串
 */
func ToInt(data interface{}) int {
	if data == nil {
		return 0
	}

	switch reflect.TypeOf(data).String() {
	case "string":
		va, _ := strconv.Atoi(data.(string))
		return va
	case "int":
		return data.(int)
	case "float64":
		return int(int64(data.(float64)))
	case "float32":
		return int(int32(data.(float32)))
	case "int64":
		return int(data.(int64))
	default:
		return 0
	}
}
