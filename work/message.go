package work

type Message struct {
}

func (w *WorkWechat) MessageSend(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + MESSAGE_SEND + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}
