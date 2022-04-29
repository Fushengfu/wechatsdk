package work

import (
	"github.com/astaxie/beego/logs"
	"strconv"
)

type User struct {
}

/**
 *  将自建应用获取的userid转换为第三方应用获取的userid
 */
func (self *WorkWechat) UseridToOpenUserid(data []interface{}) (str string, err error) {
	post := make(map[string]interface{})
	post["userid_list"] = data
	url := API_URL_PREFIX + USERID_TO_OPENUSERID + self.AccessToken
	str, err = self.sendForm("POST", url, post, data)

	return str, err
}

/**
 * 获取部门信息列表
 */
func (self *WorkWechat) GetUserList(id, fetchChild int) (str string, err error) {
	var url string
	if fetchChild == 1 {
		url = API_URL_PREFIX + GET_USER_LIST_BY_DETAIL + self.AccessToken + "&department_id=" + strconv.Itoa(id) + "&fetch_child=" + strconv.Itoa(fetchChild)
	} else {
		url = API_URL_PREFIX + GET_USER_LIST + self.AccessToken + "&department_id=" + strconv.Itoa(id) + "&fetch_child=" + strconv.Itoa(fetchChild)
	}
	str, err = self.sendForm("GET", url, nil, id, fetchChild)
	logs.Info("get_token", "Url:"+url)
	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 * 获取读取成员
 */
func (self *WorkWechat) GetUserDetail(userid string) (str string, err error) {
	var url string
	url = API_URL_PREFIX + USER_GET + self.AccessToken + "&userid=" + userid
	str, err = self.sendForm("GET", url, nil, userid)
	logs.Info("get_token", "Url:"+url)
	if err == nil {
		return str, nil
	}
	return str, err
}

/**
 *  配置客户联系「联系我」方式
 */
func (self *WorkWechat) AddContactWay(data ...interface{}) (str string, err error) {
	post := make(map[string]interface{})
	post["type"] = 2
	post["scene"] = 2
	post["style"] = 1
	post["remark"] = "测试"
	post["skip_verify"] = true
	post["state"] = 1
	post["user"] = []string{"NanShanNan", "NanShanNan01"}

	url := API_URL_PREFIX + ADD_CONTACT_WAY + self.AccessToken
	str, err = self.sendForm("POST", url, post, data...)

	return str, err
}

/**
 *  获取企业已配置的「联系我」方式
 */
func (self *WorkWechat) GetContactWay(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GET_CONTACT_WAY + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  更新企业已配置的「联系我」方式
 */
func (self *WorkWechat) UpdateContactWay(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + UPDATE_CONTACT_WAY + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  删除企业已配置的「联系我」方式
 */
func (self *WorkWechat) DelContactWay(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + DEL_CONTACT_WAY + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  获取配置了客户联系功能的成员列表
 */
func (self *WorkWechat) GetFollowUserList() (str string, err error) {
	url := API_URL_PREFIX + GET_FOLLOW_USER_LIST + self.AccessToken
	str, err = self.sendForm("GET", url, nil)
	return str, err
}
