package work

type Statistics struct {
}

/**
 *  获取联系客户统计数据
 */
func (self *WorkWechat) GetUserBehaviorData(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GET_USER_BEHAVIOR_DATA + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}
