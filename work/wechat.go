package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Fushengfu/wechatsdk/tool"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

var (
	/**
	 *  缓存数据
	 */
	StoreCacheMap map[string]*StoreCache

	/**
	 *  永久授权码列表
	 */
	PermanentCodeMap map[string]string

	/**
	 *  授权企业token列表
	 */
	AuthAccessTokenMap map[string]authAccessToken

	/**
	 *  第三方应用凭证列表
	 */
	SuiteAccessTokenMap map[string]suiteAccessToken

	/**
	 *  jsapi_ticket 列表
	 */
	JsapiTicketMap map[string]jsapiTicket

	request          tools.Http
	API_URL_PREFIX   = "https://qyapi.weixin.qq.com"
	GET_ACCESS_TOKEN = "/cgi-bin/gettoken?"

	// 部门管理
	GET_DEPARTMENT_LIST = "/cgi-bin/department/list?access_token="
	ADD_DEPARTMENT      = "/cgi-bin/department/create?access_token="
	UPDATE_DEPARTMENT   = "/cgi-bin/department/update?access_token="
	DELETE_DEPARTMENT   = "/cgi-bin/department/delete?access_token="

	// 成员管理
	USER_GET                = "/cgi-bin/user/get?access_token="
	GET_USER_LIST           = "/cgi-bin/user/simplelist?access_token="
	GET_USER_LIST_BY_DETAIL = "/cgi-bin/user/list?access_token="
	GET_MOBILE_HASCODE      = "/cgi-bin/user/get_mobile_hashcode?access_token="
	ADD_CONTACT_WAY         = "/cgi-bin/externalcontact/add_contact_way?access_token="
	GET_CONTACT_WAY         = "/cgi-bin/externalcontact/get_contact_way?access_token="
	UPDATE_CONTACT_WAY      = "/cgi-bin/externalcontact/update_contact_way?access_token="
	DEL_CONTACT_WAY         = "/cgi-bin/externalcontact/del_contact_way?access_token="
	GET_FOLLOW_USER_LIST    = "/cgi-bin/externalcontact/get_follow_user_list?access_token="
	USERID_TO_OPENUSERID    = "/cgi-bin/batch/userid_to_openuserid?access_token="

	// 客户管理
	GET_MEMBER_LIST_OF_USER = "/cgi-bin/externalcontact/list?access_token="
	GET_MEMBER_DETAIL       = "/cgi-bin/externalcontact/get?access_token="
	SET_MEMBER_REMARK       = "/cgi-bin/externalcontact/remark?access_token="
	BATCH_GET_BY_USER       = "/cgi-bin/externalcontact/batch/get_by_user?access_token="
	GET_NEW_EXTERNAL_USERID = "/cgi-bin/externalcontact/get_new_external_userid?access_token="

	// 客户标签管理
	GET_CORP_TAG_LIST = "/cgi-bin/externalcontact/get_corp_tag_list?access_token="
	ADD_CORP_TAG      = "/cgi-bin/externalcontact/add_corp_tag?access_token="
	Edit_CORP_TAG     = "/cgi-bin/externalcontact/edit_corp_tag?access_token="
	DELETE_CORP_TAG   = "/cgi-bin/externalcontact/del_corp_tag?access_token="
	MARK_TAG          = "/cgi-bin/externalcontact/mark_tag?access_token="

	// 客户群管理
	GROUP_CHAT_LIST                = "/cgi-bin/externalcontact/groupchat/list?access_token="
	GROUP_CHAT_INFO                = "/cgi-bin/externalcontact/groupchat/get?access_token="
	UNIONID_TO_EXTERNAL_USERID_3RD = "/cgi-bin/service/externalcontact/unionid_to_external_userid_3rd?suite_access_token="

	// 企业群发
	ADD_GROUP_MSG_TEMPLATE = "/cgi-bin/externalcontact/add_msg_template?access_token="
	GET_GROUP_MSG_RESULT   = "/cgi-bin/externalcontact/get_group_msg_result?access_token="

	// 统计管理
	GET_USER_BEHAVIOR_DATA = "/cgi-bin/externalcontact/get_user_behavior_data?access_token="

	// 应用管理
	GET_AGENT_LIST = "/cgi-bin/agent/list?access_token="

	// 素材管理
	MEDIA_UPLOAD      = "/cgi-bin/media/upload?access_token="
	UPLOAD_ATTACHMENT = "/cgi-bin/media/upload_attachment?access_token="
	UPLOADIMG         = "/cgi-bin/media/uploadimg?access_token="

	// 消息推送
	SEND_WELCOME_MSG = "/cgi-bin/externalcontact/send_welcome_msg?access_token="

	// 离职管理
	GRT_UNASSIGNED_LIST = "/cgi-bin/externalcontact/get_unassigned_list?access_token="
	UNASSIGNED_TRANSFER = "/cgi-bin/externalcontact/transfer?access_token="
	GROUPCHAT_TRANSFER  = "/cgi-bin/externalcontact/groupchat/transfer?access_token="

	//js-sdk算法
	GET_JSAPI_TICKET       = "/cgi-bin/get_jsapi_ticket?access_token="
	GET_AGENT_JSAPI_TICKET = "/cgi-bin/ticket/get?access_token=%v&type=agent_config"

	/**
	 *  第三方应用
	 */
	//获取第三方应用凭证
	GET_SUITE_TOKEN = "/cgi-bin/service/get_suite_token"
	//获取预授权码
	GET_PRE_AUTH_CODE = "/cgi-bin/service/get_pre_auth_code?suite_access_token="
	//设置授权配置
	SET_SESSION_INFO = "/cgi-bin/service/set_session_info?suite_access_token="
	//获取企业永久授权码
	GET_PERMANMENT_CODE = "/cgi-bin/service/get_permanent_code?suite_access_token="
	//获取企业授权信息
	GET_AUTH_INFO = "/cgi-bin/service/get_auth_info?suite_access_token="
	//获取企业凭证
	GET_CORP_TOKEN = "/cgi-bin/service/get_corp_token?suite_access_token="
	//获取应用的管理员列表
	GET_ADMIN_LIST = "/cgi-bin/service/get_admin_list?suite_access_token="
	//服务商
	GET_PROVIDER_TOKEN   = "/cgi-bin/service/get_provider_token"
	CORPID_TO_OPENCORPID = "/cgi-bin/service/corpid_to_opencorpid?provider_access_token="

	//消息推送
	//发送应用消息
	MESSAGE_SEND = "/cgi-bin/message/send?access_token="
)

type WorkWechat struct {
	ProviderCorpid      string     `json:"provider_corpid"`
	ProviderSecret      string     `json:"provider_secret"`
	CorpId              string     `json:"corp_id"`
	Secret              string     `json:"secret"`
	AccessToken         string     `json:"access_token"`
	Agentid             string     `json:"agentid"`
	SuiteId             string     `json:"suite_id"`
	SuiteSecret         string     `json:"suite_secret"`
	SuiteTicket         string     `json:"suite_ticket"`
	SuiteAccessToken    string     `json:"suite_access_token"`
	ProviderAccessToken string     `json:"provider_access_token"`
	PreAuthCode         string     `json:"pre_auth_code"`
	PermanentCode       string     `json:"permanent_code"`
	Limit               int        `json:"limit"`
	Remark              int        `json:"remark"`
	mu                  sync.Mutex `json:"-"`
}

/**
 *  缓存数据
 */
type StoreCache struct {
	Type      string `json:"type"`
	Data      string `json:"data"`
	ExpiresIn int    `json:"expires_in"`
	StartAt   int    `json:"start_at"`
	EndAt     int    `json:"end_at"`
}

/**
 *  存储键值对数据
 */
func (self *WorkWechat) SetCache(key, val string, expires int) int {
	cache := new(StoreCache)
	cache.Type = "string"
	cache.Data = val
	cache.StartAt = int(time.Now().Unix())
	cache.ExpiresIn = expires
	cache.EndAt = cache.StartAt + expires
	StoreCacheMap[key] = cache

	return 1
}

/**
 *  存储键值对数据
 */
func (self *WorkWechat) SetCacheNx(key, val string, expires int) int {
	if _, ok := StoreCacheMap[key]; ok {
		return 0
	}

	cache := new(StoreCache)
	cache.Type = "string"
	cache.Data = val
	cache.StartAt = int(time.Now().Unix())
	cache.ExpiresIn = expires
	cache.EndAt = cache.StartAt + expires
	StoreCacheMap[key] = cache

	return 1
}

/**
 *  取出缓存数据
 */
func (self *WorkWechat) GetCache(key string) (string, error) {
	_, ok := StoreCacheMap[key]

	if ok &&
		StoreCacheMap[key].EndAt >= int(time.Now().Unix()) {
		return StoreCacheMap[key].Data, nil
	}

	return "", errors.New("empty key")
}

/**
 *  删除缓存数据
 */
func (self *WorkWechat) DelCache(key string) int {

	delete(StoreCacheMap, key)
	return 1
}

/**
 *  第三方应用凭证
 */
type suiteAccessToken struct {
	Errcode          int     `json:"errcode"`
	Errmsg           *string `json:"errmsg"`
	SuiteAccessToken *string `json:"suite_access_token"`
	ExpiresIn        int     `json:"expires_in"`
}

/**
 *  授权企业凭证
 */
type authAccessToken struct {
	Errcode     int     `json:"errcode"`
	Errmsg      *string `json:"errmsg"`
	AccessToken *string `json:"access_token"`
	ExpiresIn   int     `json:"expires_in"`
}

/**
 *  获取服务商凭证
 */
type providerAccessToken struct {
	Errcode             int     `json:"errcode"`
	Errmsg              *string `json:"errmsg"`
	ProviderAccessToken *string `json:"provider_access_token"`
	ExpiresIn           int     `json:"expires_in"`
}

/**
 *  jsapi_ticket
 */
type jsapiTicket struct {
	Errcode   int     `json:"errcode"`
	Errmsg    *string `json:"errmsg"`
	Ticket    string  `json:"ticket"`
	ExpiresIn int     `json:"expires_in"`
}

/**
 *  获取token
 */
func (self *WorkWechat) GetWorkWechat() *WorkWechat {
	if len(self.SuiteId) == 0 {
		key := "WECHAT::" + self.Secret + "qy_access_token"
		res, _ := self.GetCache(key)
		if res != "" {
			data := make(map[string]interface{})
			if err := json.Unmarshal([]byte(res), &data); err == nil {
				self.AccessToken = data["access_token"].(string)
			} else {
				fmt.Println("解析失败")
				fmt.Println(err.Error())
				return nil
			}
		} else {
			result, err := self.getAccessToken()
			if err != nil {
				return nil
			}
			self.AccessToken = result
		}
	} else {
		self.GetSuiteToken()
		_, err := self.GetCorpToken()
		if err != nil {
			logs.Critical("获取access_token异常：", err.Error())
		}

		_, er := self.GetProviderToken(self.ProviderCorpid, self.ProviderSecret)
		if er != nil {
			logs.Critical("获取provider_access_token异常：", er.Error())
		}
	}

	return self
}
