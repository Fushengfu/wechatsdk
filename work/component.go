package work

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"net/url"
)

type Component struct {
}

/**
 *  获取第三方应用凭证
@return
{
    "errcode":0 ,
    "errmsg":"ok" ,
    "suite_access_token":"61W3mEpU66027wgNZ_MhGHNQDHnFATkDa9-2llMBjUwxRSNPbVsMmyD-yq8wZETSoE5NQgecigDrSHkPtIYA",
    "expires_in":7200
}
*/
func (w *WorkWechat) GetSuiteToken() *WorkWechat {
	w.mu.Lock()
	defer w.mu.Unlock()
	key := "WECHAT_QY::SUITE_ACCESS_TOKEN_" + w.SuiteId

	str := Rds.Get(key)
	data := make(map[string]interface{})
	var er error
	if str == "" {
		body := map[string]interface{}{
			"suite_id":     w.SuiteId,
			"suite_secret": w.SuiteSecret,
			"suite_ticket": w.SuiteTicket,
		}

		url := API_URL_PREFIX + GET_SUITE_TOKEN
		result, err := w.sendForm("POST", url, body)
		logs.Critical(url, result, err)
		if err == nil {
			er = json.Unmarshal([]byte(result), &data)

			_, ok := data["errcode"]
			if er == nil && (!ok || ToInt(data["errcode"]) == 0) {
				Rds.Set(key, result, ToInt(data["expires_in"])-120)
			}
		}

	} else {
		er = json.Unmarshal([]byte(str), &data)
	}

	errcode, ok := data["errcode"]
	if er != nil {
		logs.Critical("解析SUITE_ACCESS_TOKEN失败：", er.Error(), str)
	} else if !ok || ToInt(errcode) == 0 {
		w.SuiteAccessToken = ToString(data["suite_access_token"])
	}

	return w
}

/**
 *  获取预授权码
@return
{
    "errcode":0 ,
    "errmsg":"ok" ,
    "pre_auth_code":"Cx_Dk6qiBE0Dmx4EmlT3oRfArPvwSQ-oa3NL_fwHM7VI08r52wazoZX2Rhpz1dEw",
    "expires_in":1200
}
*/
func (w *WorkWechat) GetPreAuthCode() (str string, err error) {
	w.GetSuiteToken()
	url := API_URL_PREFIX + GET_PRE_AUTH_CODE + w.SuiteAccessToken
	str, err = w.sendForm("GET", url, nil)
	return str, err
}

/**
 *  设置授权配置
@return
{
    "errcode": 0,
    "errmsg": "ok"
}
*/
func (w *WorkWechat) SetSessionInfo(appids []string, authType int) (str string, err error) {
	w.GetPreAuthCode()
	data := map[string]interface{}{
		"pre_auth_code": w.PreAuthCode,
		"session_info": map[string]interface{}{
			"appid":     appids,
			"auth_type": authType, //授权类型：0 正式授权， 1 测试授权。 默认值为0。注意，请确保应用在正式发布后的授权类型为“正式授权”
		},
	}

	url := API_URL_PREFIX + SET_SESSION_INFO + w.SuiteAccessToken
	str, err = w.sendForm("POST", url, data, appids, authType)
	return str, err
}

/**
 *  获取企业永久授权码
@return
{
    "errcode":0 ,
    "errmsg":"ok" ,
    "access_token": "xxxxxx",
    "expires_in": 7200,
    "permanent_code": "xxxx",
    "dealer_corp_info":
    {
        "corpid": "xxxx",
        "corp_name": "name"
    },
    "auth_corp_info":
    {
        "corpid": "xxxx",
        "corp_name": "name",
        "corp_type": "verified",
        "corp_square_logo_url": "yyyyy",
        "corp_user_max": 50,
        "corp_agent_max": 30,
        "corp_full_name":"full_name",
        "verified_end_time":1431775834,
        "subject_type": 1,
        "corp_wxqrcode": "zzzzz",
        "corp_scale": "1-50人",
        "corp_industry": "IT服务",
        "corp_sub_industry": "计算机软件/硬件/信息服务",
        "location":"广东省广州市",
        "auth_type":1
    },
    "auth_info":
    {
        "agent" :
        [
            {
                "agentid":1,
                "name":"NAME",
                "round_logo_url":"xxxxxx",
                "square_logo_url":"yyyyyy",
                "appid":1,
                "privilege":
                {
                    "level":1,
                    "allow_party":[1,2,3],
                    "allow_user":["zhansan","lisi"],
                    "allow_tag":[1,2,3],
                    "extra_party":[4,5,6],
                    "extra_user":["wangwu"],
                    "extra_tag":[4,5,6]
                },
                "shared_from":
                {
                    "corpid":"wwyyyyy"
                }
            },
            {
                "agentid":2,
                "name":"NAME2",
                "round_logo_url":"xxxxxx",
                "square_logo_url":"yyyyyy",
                "appid":5,
                "shared_from":
                {
                    "corpid":"wwyyyyy"
                }
            }
        ]
    },
    "auth_user_info":
    {
        "userid":"aa",
        "open_userid":"xxxxxx",
        "name":"xxx",
        "avatar":"http://xxx"
    },
    "register_code_info":
    {
        "register_code":"1111",
        "template_id":"tpl111",
        "state":"state001"
    }
}
*/
func (w *WorkWechat) GetPermanentCode(auth_code string) (str string, err error) {
	w.GetSuiteToken()
	data := map[string]interface{}{
		"auth_code": auth_code,
	}

	url := API_URL_PREFIX + GET_PERMANMENT_CODE + w.SuiteAccessToken
	str, err = w.sendForm("POST", url, data)
	if err == nil {
		info := make(map[string]interface{})
		er := json.Unmarshal([]byte(str), &info)
		_, ok := info["errcode"]
		logs.Critical("PERMANENT_CODE:", er == nil && (!ok || ToInt(info["errcode"]) == 0), "permanent_code:", info["permanent_code"].(string))
		if er == nil && (!ok || ToInt(info["errcode"]) == 0) {
			auth_corp_info := info["auth_corp_info"].(map[string]interface{})
			logs.Critical("PERMANENT_CODE::"+auth_corp_info["corpid"].(string), "permanent_code:", info["permanent_code"].(string))
			Rds.Set("PERMANENT_CODE::"+auth_corp_info["corpid"].(string), info["permanent_code"].(string), 0)
		}
	}

	return str, err
}

/**
 *  获取企业授权信息
@return
{
    "errcode":0 ,
    "errmsg":"ok" ,
    "dealer_corp_info":
    {
        "corpid": "xxxx",
        "corp_name": "name"
    },
    "auth_corp_info":
    {
        "corpid": "xxxx",
        "corp_name": "name",
        "corp_type": "verified",
        "corp_square_logo_url": "yyyyy",
        "corp_user_max": 50,
        "corp_agent_max": 30,
        "corp_full_name":"full_name",
        "verified_end_time":1431775834,
        "subject_type": 1,
        "corp_wxqrcode": "zzzzz",
        "corp_scale": "1-50人",
        "corp_industry": "IT服务",
        "corp_sub_industry": "计算机软件/硬件/信息服务",
        "location":"广东省广州市"
    },
    "auth_info":
    {
        "agent" :
        [
            {
                "agentid":1,
                "name":"NAME",
                "round_logo_url":"xxxxxx",
                "square_logo_url":"yyyyyy",
                "appid":1,
                "privilege":
                {
                    "level":1,
                    "allow_party":[1,2,3],
                    "allow_user":["zhansan","lisi"],
                    "allow_tag":[1,2,3],
                    "extra_party":[4,5,6],
                    "extra_user":["wangwu"],
                    "extra_tag":[4,5,6]
                },
                "shared_from":
                {
                    "corpid":"wwyyyyy"
                }
            },
            {
                "agentid":2,
                "name":"NAME2",
                "round_logo_url":"xxxxxx",
                "square_logo_url":"yyyyyy",
                "appid":5,
                "shared_from":
                {
                    "corpid":"wwyyyyy"
                }
            }
        ]
    }
}
*/
func (w *WorkWechat) GetAuthInfo(auth_corpid string) (str string, err error) {
	data := map[string]interface{}{
		"auth_corpid":    auth_corpid,
		"permanent_code": Rds.Get("PERMANENT_CODE::" + w.Appid),
	}

	url := API_URL_PREFIX + GET_AUTH_INFO + w.SuiteAccessToken
	str, err = w.sendForm("POST", url, data)
	return str, err
}

/**
 *  获取企业凭证
@return
{
    "errcode":0 ,
    "errmsg":"ok" ,
    "access_token": "xxxxxx",
    "expires_in": 7200
}
*/
func (w *WorkWechat) GetCorpToken() (str string, err error) {
	key := "WECHAT_QY::AUTH_CORPID_" + w.Appid
	logs.Critical("WECHAT_QY_AUTH_CORPID:", key)
	str = Rds.Get(key)

	if str != "" {
		tmp := make(map[string]interface{})
		er := json.Unmarshal([]byte(str), &tmp)
		if er == nil && w.AccessToken == ToString(tmp["access_token"]) {
			Rds.Del(key)
		} else if er == nil {
			w.AccessToken = ToString(tmp["access_token"])
		} else {
			logs.Critical("解析TOLKEN信息异常：", er, "AccessToken:", w.AccessToken, "tmp:", tmp)
		}
	}

	if str == "" {
		permanentCode := Rds.Get("PERMANENT_CODE::" + w.Appid)
		if permanentCode == "" {
			return str, errors.New("PERMANENT_CODE 为空")
		}

		data := map[string]interface{}{
			"auth_corpid":    w.Appid,
			"permanent_code": permanentCode,
		}
		logs.Critical("strstrstr:", data, w.Appid)

		url := API_URL_PREFIX + GET_CORP_TOKEN + w.SuiteAccessToken
		str, err = w.sendForm("POST", url, data)
		if err != nil {
			return str, err
		}
	}

	info := make(map[string]interface{})
	er := json.Unmarshal([]byte(str), &info)
	errcode, ok := info["errcode"]

	if er == nil && (!ok || ToInt(errcode) == 0) {
		w.AccessToken = ToString(info["access_token"])
		Rds.Set(key, str, ToInt(info["expires_in"])-120)
	}

	return str, err
}

/**
 *  获取应用的管理员列表
@return
{
    "errcode": 0,
    "errmsg": "ok",
    "admin":[
        {"userid":"zhangsan","open_userid":"xxxxx","auth_type":1},
        {"userid":"lisi","open_userid":"yyyyy","auth_type":0}
    ]
}
*/
func (w *WorkWechat) GetAdminList(auth_corpid string) (str string, err error) {
	data := map[string]interface{}{
		"auth_corpid": "auth_corpid_value",
		"agentid":     1000046,
	}

	url := API_URL_PREFIX + GET_ADMIN_LIST + w.SuiteAccessToken
	str, err = w.sendForm("POST", url, data)
	return str, err
}

/**
 *  网页授权登录
 */
func (w *WorkWechat) GetAuthorizeUrl(uri string) string {
	uri, _ = url.QueryUnescape(uri)
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + w.SuiteId + "&redirect_uri=" + uri + "&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect"
	return url
}

/**
 *  获取访问用户身份
 */
func (w *WorkWechat) GetUserInfo(code string) (string, error) {
	uri := "https://qyapi.weixin.qq.com/cgi-bin/service/getuserinfo3rd?suite_access_token=" + w.SuiteAccessToken + "&code=" + code
	return request.Get(uri, nil)
}
