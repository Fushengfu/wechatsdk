package work

type Group struct {
}

/**
 *  获取客户群列表
 */
func (self *WorkWechat) GetGroupList(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GROUP_CHAT_LIST + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  获取客户群详情
 */
func (self *WorkWechat) GetGroupInfo(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GROUP_CHAT_INFO + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}
