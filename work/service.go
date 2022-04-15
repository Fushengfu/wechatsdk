package work

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

/**
 *  获取服务商凭证
 */
//{
//"corpid":"xxxxx",
//"provider_secret":"xxx"
//}
func (w *WorkWechat) GetProviderToken(corpid, providerSecret string) (str string, err error) {
	key := "WECHAT_QY::PROVIDER_TOKEN_" + corpid
	str = Rds.Get(key)

	if str != "" {
		tmp := make(map[string]interface{})
		er := json.Unmarshal([]byte(str), &tmp)
		if er == nil && w.ProviderAccessToken == ToString(tmp["provider_access_token"]) {
			Rds.Del(key)
		} else if er == nil {
			w.ProviderAccessToken = ToString(tmp["provider_access_token"])
		} else {
			logs.Critical("解析TOLKEN信息异常：", er, "ProviderAccessToken:", w.ProviderAccessToken, "tmp:", tmp)
		}
	}

	if str == "" {
		url := API_URL_PREFIX + GET_PROVIDER_TOKEN
		data := make(map[string]interface{})
		data["corpid"] = corpid
		data["provider_secret"] = providerSecret
		str, err = w.sendForm("POST", url, data, corpid, providerSecret)

		info := make(map[string]interface{})
		er := json.Unmarshal([]byte(str), &info)
		errcode, ok := info["errcode"]

		if er == nil && (!ok || ToInt(errcode) == 0) {
			w.ProviderAccessToken = ToString(info["provider_access_token"])
			Rds.Set(key, str, ToInt(info["expires_in"])-120)
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
