package work

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"net/url"
	"reflect"
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
func (self *WorkWechat) GetSuiteToken() *WorkWechat {
	self.mu.Lock()
	defer self.mu.Unlock()
	key := "WECHAT_QY::SUITE_ACCESS_TOKEN_" + self.SuiteId
	var err error

	var suiteAccToken suiteAccessToken

	str, err := self.GetCache(key)
	if err != nil {
		body := map[string]interface{}{
			"suite_id":     self.SuiteId,
			"suite_secret": self.SuiteSecret,
			"suite_ticket": self.SuiteTicket,
		}

		url := API_URL_PREFIX + GET_SUITE_TOKEN
		str, err = self.sendForm("POST", url, body)
		logs.Critical(url, str, err)
	}

	if err == nil {
		err = json.Unmarshal([]byte(str), &suiteAccToken)
		if err == nil && suiteAccToken.Errcode == 0 {
			self.SetCacheNx(key, str, suiteAccToken.ExpiresIn-120)
		}
	}

	if suiteAccToken.Errcode != 0 || suiteAccToken.SuiteAccessToken == nil {
		logs.Critical("解析SUITE_ACCESS_TOKEN失败：", str)
	} else {
		self.SuiteAccessToken = reflect.ValueOf(suiteAccToken.SuiteAccessToken).Elem().String()
	}

	return self
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
func (self *WorkWechat) GetPreAuthCode() (str string, err error) {
	self.GetSuiteToken()
	url := API_URL_PREFIX + GET_PRE_AUTH_CODE + self.SuiteAccessToken
	str, err = self.sendForm("GET", url, nil)
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
func (self *WorkWechat) SetSessionInfo(appids []string, authType int) (str string, err error) {
	self.GetPreAuthCode()
	data := map[string]interface{}{
		"pre_auth_code": self.PreAuthCode,
		"session_info": map[string]interface{}{
			"appid":     appids,
			"auth_type": authType, //授权类型：0 正式授权， 1 测试授权。 默认值为0。注意，请确保应用在正式发布后的授权类型为“正式授权”
		},
	}

	url := API_URL_PREFIX + SET_SESSION_INFO + self.SuiteAccessToken
	str, err = self.sendForm("POST", url, data, appids, authType)
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
func (self *WorkWechat) GetPermanentCode(authCode string) (str string, err error) {
	self.GetSuiteToken()
	data := map[string]interface{}{
		"auth_code": authCode,
	}

	url := API_URL_PREFIX + GET_PERMANMENT_CODE + self.SuiteAccessToken
	str, err = self.sendForm("POST", url, data)
	if err == nil {
		info := make(map[string]interface{})
		er := json.Unmarshal([]byte(str), &info)
		_, ok := info["errcode"]

		if er == nil && (!ok || ToInt(info["errcode"]) == 0) {
			authCorpInfo := info["auth_corp_info"].(map[string]interface{})
			self.SetCacheNx("PERMANENT_CODE::"+"_"+self.SuiteId+"_"+authCorpInfo["corpid"].(string), info["permanent_code"].(string), 0)
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
func (self *WorkWechat) GetAuthInfo(authCorpid string) (str string, err error) {
	data := map[string]interface{}{
		"auth_corpid":    authCorpid,
		"permanent_code": self.GetCache("PERMANENT_CODE::" + "_" + self.SuiteId + "_" + authCorpid),
	}

	url := API_URL_PREFIX + GET_AUTH_INFO + self.SuiteAccessToken
	str, err = self.sendForm("POST", url, data)
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
func (self *WorkWechat) GetCorpToken() (str string, err error) {
	key := "WECHAT_QY::AUTH_CORPID_" + self.SuiteId + "_" + self.CorpId
	logs.Critical("WECHAT_QY_AUTH_CORPID:", key)
	str, err = self.GetCache(key)

	var authAccToken authAccessToken

	if err != nil {
		data := map[string]interface{}{
			"auth_corpid":    self.CorpId,
			"permanent_code": self.PermanentCode,
		}

		logs.Critical("strstrstr:", data, self.CorpId)

		url := API_URL_PREFIX + GET_CORP_TOKEN + self.SuiteAccessToken
		str, err = self.sendForm("POST", url, data)
	}

	if err != nil {
		return str, err
	}

	er := json.Unmarshal([]byte(str), &authAccToken)
	if er == nil && authAccToken.Errcode == 0 {
		self.AccessToken = reflect.ValueOf(authAccToken.AccessToken).Elem().String()
		self.SetCacheNx(key, str, authAccToken.ExpiresIn-120)
	}

	return str, err
}

/**
 * 获取accessToken
 */
func (self *WorkWechat) getAccessToken() (str string, er error) {
	key := "WECHAT::" + self.Secret + "qy_access_token"
	url := API_URL_PREFIX + GET_ACCESS_TOKEN + "corpid=" + self.CorpId + "&corpsecret=" + self.Secret

	var tokenInfo authAccessToken
	for i := 0; i < 3; i++ {
		res, err := self.sendForm("GET", url, nil)

		if err != nil {
			return str, err
		}

		if err = json.Unmarshal([]byte(res), &tokenInfo); err == nil && tokenInfo.Errcode == 0 {
			self.SetCacheNx(key, res, tokenInfo.ExpiresIn-120)

			return reflect.ValueOf(tokenInfo.AccessToken).Elem().String(), nil

		} else {
			return str, err
		}
	}

	return str, er
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
func (self *WorkWechat) GetAdminList(authCorpid, agentid interface{}) (str string, err error) {
	data := map[string]interface{}{
		"auth_corpid": authCorpid,
		"agentid":     agentid,
	}

	url := API_URL_PREFIX + GET_ADMIN_LIST + self.SuiteAccessToken
	str, err = self.sendForm("POST", url, data)
	return str, err
}

/**
 *  网页授权登录
 */
func (self *WorkWechat) GetAuthorizeUrl(uri string) string {
	uri, _ = url.QueryUnescape(uri)
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + self.SuiteId + "&redirect_uri=" + uri + "&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect"
	return url
}

/**
 *  获取访问用户身份
 */
func (self *WorkWechat) GetUserInfo(code string) (string, error) {
	uri := "https://qyapi.weixin.qq.com/cgi-bin/service/getuserinfo3rd?suite_access_token=" + self.SuiteAccessToken + "&code=" + code
	return request.Get(uri, nil)
}
