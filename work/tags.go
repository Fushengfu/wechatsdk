package work

type Tags struct {
}

/**
 *  获取企业标签库
 */
func (w *WorkWechat) GetCorpTagList(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + GET_CORP_TAG_LIST + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  添加企业客户标签
 */
func (w *WorkWechat) AddCorpTag(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + ADD_CORP_TAG + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  编辑企业客户标签
 */
func (w *WorkWechat) EditCorpTag(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + Edit_CORP_TAG + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  删除企业客户标签
 */
func (w *WorkWechat) DeleteCorpTag(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + DELETE_CORP_TAG + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  编辑客户企业标签
 */
func (w *WorkWechat) MarkTag(data map[string]interface{}) (str string, err error) {
	url := API_URL_PREFIX + MARK_TAG + w.AccessToken
	str, err = w.sendForm("POST", url, data, data)
	return str, err
}
