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
func (self *WorkWechat) GetDepartmentList(id int) (str string, err error) {
	url := API_URL_PREFIX + GET_DEPARTMENT_LIST + self.AccessToken + "&id=" + strconv.Itoa(id)
	str, err = self.sendForm("GET", url, nil, id)
	logs.Info("get_token", "Url:"+url)

	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 *  添加部门
 */
func (self *WorkWechat) CreateDepartment(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + ADD_DEPARTMENT + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 *  更新部门
 */
func (self *WorkWechat) EditDepartment(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + UPDATE_DEPARTMENT + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 * 删除部门
 */
func (self *WorkWechat) DeleteDepartment(id int) (str string, err error) {
	url := API_URL_PREFIX + DELETE_DEPARTMENT + self.AccessToken + "&id=" + strconv.Itoa(id)
	str, err = self.sendForm("GET", url, nil, id)
	if err == nil {
		return str, nil
	}
	return str, err
}
