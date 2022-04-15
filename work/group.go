package work

type Group struct {
}

/**
 *  获取客户群列表
 */
func (w *WorkWechat) GetGroupList(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GROUP_CHAT_LIST + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  获取客户群详情
 */
func (w *WorkWechat) GetGroupInfo(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GROUP_CHAT_INFO + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}
