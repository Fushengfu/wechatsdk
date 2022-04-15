package amulet

import "github.com/Fushengfu/wechatsdk/work"

/**
 *  获取企业微信实例
    weixin := &work.WorkWechat{
		CorpId:       option["corpid"],
		Secret:      option["secret"],
		SuiteId:     option["suite_id"],
		SuiteSecret: option["suite_secret"],
		SuiteTicket: option["suite_ticket"],
		PermanentCode: option["permanent_code"],
	}
*/
func NewWorkWechat(option map[string]string) (w *work.WorkWechat) {
	if _, ok := option["corpid"]; !ok {
		option["corpid"] = ""
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
		CorpId:        option["corpid"],
		Secret:        option["secret"],
		SuiteId:       option["suite_id"],
		SuiteSecret:   option["suite_secret"],
		SuiteTicket:   option["suite_ticket"],
		PermanentCode: option["permanent_code"],
	}

	return weixin.GetWorkWechat()
}
