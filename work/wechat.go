package work

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"lmkweixin/amulet/tool"
	"lmkweixin/dbs"
	"lmkweixin/tools/redis"
	"sync"
)

var (
	Rds              *redis.RedisClient
	request          tool.Request
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
	Appid               string
	Secret              string
	AccessToken         string
	SuiteId             string
	SuiteSecret         string
	SuiteTicket         string
	SuiteAccessToken    string
	ProviderAccessToken string
	PreAuthCode         string
	PermanentCode       string
	Limit               int
	Remark              int
	mu                  sync.Mutex
}

func init() {
	Rds = redis.NewRedisClient(dbs.Rds)
}

/**
 *  获取token
 */
func (w *WorkWechat) GetWorkWechat() *WorkWechat {
	if len(w.SuiteId) == 0 {
		key := "WECHAT::" + w.Secret + "qy_access_token"
		res := Rds.Get(key)
		fmt.Println("获取TOKEN:", res)
		if res != "" {
			data := make(map[string]interface{})
			if err := json.Unmarshal([]byte(res), &data); err == nil {
				w.AccessToken = data["access_token"].(string)
			} else {
				fmt.Println("解析失败")
				fmt.Println(err.Error())
				return nil
			}
		} else {
			result, err := w.getAccessToken()
			if err != nil {
				return nil
			}
			w.AccessToken = result
		}
	} else {
		w.GetSuiteToken()
		_, err := w.GetCorpToken()
		if err != nil {
			logs.Critical("获取access_token异常：", err.Error())
		}

		_, er := w.GetProviderToken("ww98ddff0207beced9", "3N7iHhSJfEEx7DNYMWsEcZT-U6p5X19p01GYQQOZedlUzxraeQWNixyHf6QN9IFd")
		if er != nil {
			logs.Critical("获取provider_access_token异常：", er.Error())
		}
	}

	return w
}

/**
 * 获取token
 */
func (w *WorkWechat) getAccessToken() (str string, er error) {
	key := "WECHAT::" + w.Secret + "qy_access_token"
	url := API_URL_PREFIX + GET_ACCESS_TOKEN + "corpid=" + w.Appid + "&corpsecret=" + w.Secret

	data := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		res, err := w.sendForm("GET", url, nil)

		if err != nil {
			return str, err
		}

		if err = json.Unmarshal([]byte(res), &data); err == nil && ToInt(data["errcode"]) == 0 {
			Rds.Set(key, res, ToInt(data["expires_in"])-120)
			return data["access_token"].(string), nil
		} else {
			return str, err
		}
	}

	return str, er
}
