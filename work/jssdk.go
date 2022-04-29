package work

import (
	"encoding/json"
	"fmt"
)

/**
 *  获取企业的jsapi_ticket
 */
func (self *WorkWechat) GetJsapiTicket() (str string, err error) {
	key := "WECHAT_QY::GET_JSAPI_TICKET_" + self.CorpId + "_" + self.Agentid
	str, err = self.GetCache(key)
	var jsTicket jsapiTicket

	if str == "" {
		url := API_URL_PREFIX + GET_JSAPI_TICKET + self.AccessToken
		str, err = self.sendForm("GET", url, nil)
	}

	if err == nil {
		er := json.Unmarshal([]byte(str), &jsTicket)
		if er == nil && jsTicket.Errcode == 0 {
			self.SetCacheNx(key, str, jsTicket.ExpiresIn-120)
		}
	}

	return str, err
}

/**
 *  获取应用的jsapi_ticket
 */
func (self *WorkWechat) GetAgentJsapiTicket() (str string, err error) {
	key := "WECHAT_QY::GET_AGENT_JSAPI_TICKET_" + self.CorpId + "_" + self.Agentid
	str, err = self.GetCache(key)

	var jsTicket jsapiTicket

	if str == "" {
		url := fmt.Sprintf(API_URL_PREFIX+GET_AGENT_JSAPI_TICKET, self.AccessToken)
		str, err = self.sendForm("GET", url, nil)
	}

	if err == nil {
		er := json.Unmarshal([]byte(str), &jsTicket)
		if er == nil && jsTicket.Errcode == 0 {
			self.SetCacheNx(key, str, jsTicket.ExpiresIn-120)
		}
	}

	return str, err
}
