package amulet

import (
	"github.com/astaxie/beego"
	"lmkweixin/amulet/weixin"
	"lmkweixin/amulet/work"
)

/**
 *  获取微信实例
 */
func NewWechat(option map[string]string) (w *weixin.Wechat) {
	wechat := &weixin.Wechat{
		Model:          option["model"],
		Appid:          option["appid"],
		ComponentAppid: beego.AppConfig.String("component::appid"),
		Secret:         beego.AppConfig.String("component::secret"),
	}

	return wechat.GetWechat()
}

/**
 *  获取企业微信实例
    weixin := &work.WorkWechat{
		Appid:       option["appid"],
		Secret:      option["secret"],
		SuiteId:     option["suite_id"],
		SuiteSecret: option["suite_secret"],
		SuiteTicket: option["suite_ticket"],
		PermanentCode: option["permanent_code"],
	}
*/
func NewWorkWechat(option map[string]string) (w *work.WorkWechat) {
	if _, ok := option["appid"]; !ok {
		option["appid"] = ""
	}

	if _, ok := option["secret"]; !ok {
		option["secret"] = ""
	}

	if _, ok := option["suite_id"]; !ok {
		option["suite_id"] = ""
	}

	if _, ok := option["suite_secret"]; !ok {
		option["suite_secret"] = ""
	}

	if _, ok := option["suite_ticket"]; !ok {
		option["suite_ticket"] = ""
	}

	if _, ok := option["permanent_code"]; !ok {
		option["permanent_code"] = ""
	}

	weixin := &work.WorkWechat{
		Appid:         option["appid"],
		Secret:        option["secret"],
		SuiteId:       option["suite_id"],
		SuiteSecret:   option["suite_secret"],
		SuiteTicket:   option["suite_ticket"],
		PermanentCode: option["permanent_code"],
	}

	return weixin.GetWorkWechat()
}
