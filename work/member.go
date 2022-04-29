package work

type Member struct {
}

/**
 *  获取客户列表
 */
func (self *WorkWechat) GetMemberList(userid string) (str string, err error) {
	url := API_URL_PREFIX + GET_MEMBER_LIST_OF_USER + self.AccessToken + "&userid=" + userid
	str, err = self.sendForm("GET", url, nil, userid)
	return str, err
}

/**
 *  获取客户详情
 */
func (self *WorkWechat) GetMemberInfo(userid string) (str string, err error) {
	url := API_URL_PREFIX + GET_MEMBER_DETAIL + self.AccessToken + "&external_userid=" + userid
	str, err = self.sendForm("GET", url, nil, userid)
	return str, err
}

/**
 *  设置客户备注
 */
func (self *WorkWechat) SetMemberRemark(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + SET_MEMBER_REMARK + self.AccessToken
	str, err = self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  批量获取获取客户详情
 */
func (self *WorkWechat) GetMemberBatchInfo(userid string, cursor string) (str string, err error) {
	data := map[string]interface{}{
		"userid": userid,
		"cursor": cursor,
		"limit":  100,
	}
	url := API_URL_PREFIX + BATCH_GET_BY_USER + self.AccessToken
	str, err = self.sendForm("POST", url, data, userid, cursor)
	return str, err
}

/**
 *  转换external_userid
 */
func (self *WorkWechat) GetNewExternalUserid(data []interface{}) (str string, err error) {
	url := API_URL_PREFIX + GET_NEW_EXTERNAL_USERID + self.AccessToken
	post := make(map[string]interface{})
	post["external_userid_list"] = data
	str, err = self.sendForm("POST", url, post, data)
	return str, err
}

/**
 *  unionid查询external_userid
 */
func (self *WorkWechat) GetUnionidToExternalUserid(post map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + UNIONID_TO_EXTERNAL_USERID_3RD + self.SuiteAccessToken
	str, err = self.sendForm("POST", url, post, post)
	return str, err
}
