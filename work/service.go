package work

import (
	"encoding/json"
	"reflect"
)

/**
 *  获取服务商凭证
 */
//{
//"corpid":"xxxxx",
//"provider_secret":"xxx"
//}
func (self *WorkWechat) GetProviderToken(corpid, providerSecret string) (str string, err error) {
	key := "WECHAT_QY::PROVIDER_TOKEN_" + corpid
	str, err = self.GetCache(key)

	var providerAccToken providerAccessToken

	if str == "" {
		url := API_URL_PREFIX + GET_PROVIDER_TOKEN
		data := make(map[string]interface{})
		data["corpid"] = corpid
		data["provider_secret"] = providerSecret
		str, err = self.sendForm("POST", url, data, corpid, providerSecret)
	}

	if str != "" {
		er := json.Unmarshal([]byte(str), &providerAccToken)
		if er == nil && providerAccToken.Errcode == 0 {
			self.ProviderAccessToken = reflect.ValueOf(providerAccToken.ProviderAccessToken).Elem().String()
			self.SetCacheNx(key, str, providerAccToken.ExpiresIn-120)
		}
	}

	return str, err
}

/**
 *  corpid转换
 */
//{
//"corpid":"xxxxx"
//}
func (w *WorkWechat) CorpidToOpencorpid(corpid string) (str string, err error) {
	url := API_URL_PREFIX + CORPID_TO_OPENCORPID + w.SuiteAccessToken
	data := make(map[string]interface{})
	data["corpid"] = corpid
	str, err = w.sendForm("POST", url, data, corpid)
	return str, err
}
