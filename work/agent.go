package work

type Agent struct {
}

/**
 *  获取应用列表
 */
func (self *WorkWechat) GetAgentList() (str string, err error) {
	url := API_URL_PREFIX + GET_AGENT_LIST + self.AccessToken
	str, err = self.sendForm("GET", url, nil)
	return str, err
}
