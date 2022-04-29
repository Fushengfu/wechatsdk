package work

type Message struct {
}

func (self *WorkWechat) MessageSend(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + MESSAGE_SEND + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}
