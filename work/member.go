package work

type Member struct {
}

/**
 *  获取客户列表
 */
func (w *WorkWechat) GetMemberList(userid string) (str string, err error) {
	url := API_URL_PREFIX + GET_MEMBER_LIST_OF_USER + w.AccessToken + "&userid=" + userid
	str, err = w.sendForm("GET", url, nil, userid)
	return str, err
}

/**
 *  获取客户详情
 */
func (w *WorkWechat) GetMemberInfo(userid string) (str string, err error) {
	url := API_URL_PREFIX + GET_MEMBER_DETAIL + w.AccessToken + "&external_userid=" + userid
	str, err = w.sendForm("GET", url, nil, userid)
	return str, err
}

/**
 *  设置客户备注
 */
func (w *WorkWechat) SetMemberRemark(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + SET_MEMBER_REMARK + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  批量获取获取客户详情
 */
func (w *WorkWechat) GetMemberBatchInfo(userid string, cursor string) (str string, err error) {
	data := map[string]interface{}{
		"userid": userid,
		"cursor": cursor,
		"limit":  100,
	}
	url := API_URL_PREFIX + BATCH_GET_BY_USER + w.AccessToken
	str, err = w.sendForm("POST", url, data, userid, cursor)
	return str, err
}

/**
 *  转换external_userid
 */
func (w *WorkWechat) GetNewExternalUserid(data []interface{}) (str string, err error) {
	url := API_URL_PREFIX + GET_NEW_EXTERNAL_USERID + w.AccessToken
	post := make(map[string]interface{})
	post["external_userid_list"] = data
	str, err = w.sendForm("POST", url, post, data)
	return str, err
}

/**
 *  unionid查询external_userid
 */
func (w *WorkWechat) GetUnionidToExternalUserid(post map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + UNIONID_TO_EXTERNAL_USERID_3RD + w.SuiteAccessToken
	str, err = w.sendForm("POST", url, post, post)
	return str, err
}
