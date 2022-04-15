package work

type Agent struct {
}

/**
 *  获取应用列表
 */
func (w *WorkWechat) GetAgentList() (str string, err error) {
	url := API_URL_PREFIX + GET_AGENT_LIST + w.AccessToken
	str, err = w.sendForm("GET", url, nil)
	return str, err
}
