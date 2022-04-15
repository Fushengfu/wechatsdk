package work

import (
	"github.com/astaxie/beego/logs"
	"strconv"
)

type Department struct {
}

/**
 * 获取部门信息列表
 */
func (w *WorkWechat) GetDepartmentList(id int) (str string, err error) {
	url := API_URL_PREFIX + GET_DEPARTMENT_LIST + w.AccessToken + "&id=" + strconv.Itoa(id)
	str, err = w.sendForm("GET", url, nil, id)
	logs.Info("get_token", "Url:"+url)

	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 *  添加部门
 */
func (w *WorkWechat) CreateDepartment(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + ADD_DEPARTMENT + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 *  更新部门
 */
func (w *WorkWechat) EditDepartment(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + UPDATE_DEPARTMENT + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 * 删除部门
 */
func (w *WorkWechat) DeleteDepartment(id int) (str string, err error) {
	url := API_URL_PREFIX + DELETE_DEPARTMENT + w.AccessToken + "&id=" + strconv.Itoa(id)
	str, err = w.sendForm("GET", url, nil, id)
	if err == nil {
		return str, nil
	}
	return str, err
}
