package work

import "fmt"

/**
 *  获取企业的jsapi_ticket
 */
func (w *WorkWechat) GetJsapiTicket() (str string, err error) {
	key := "WECHAT_QY::GET_JSAPI_TICKET_" + w.Appid
	str = Rds.Get(key)

	if str != "" {
		return str, nil
	}

	url := API_URL_PREFIX + GET_JSAPI_TICKET + w.AccessToken
	str, err = w.sendForm("GET", url, nil)
	return str, err
}

/**
 *  获取应用的jsapi_ticket
 */
func (w *WorkWechat) GetAgentJsapiTicket() (str string, err error) {
	key := "WECHAT_QY::GET_AGENT_JSAPI_TICKET_" + w.Appid
	str = Rds.Get(key)

	if str != "" {
		return str, nil
	}

	url := fmt.Sprintf(API_URL_PREFIX+GET_AGENT_JSAPI_TICKET, w.AccessToken)
	str, err = w.sendForm("GET", url, nil)
	return str, err
}
