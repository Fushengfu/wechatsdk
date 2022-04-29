package work

type Media struct {
}

/**
 *  上传附件资源
 */
func (self *WorkWechat) UploadAttachment(data map[string]interface{}) (string, error) {
	url := API_URL_PREFIX + UPLOAD_ATTACHMENT + self.AccessToken + "&media_type=" + ToString(data["media_type"]) + "&attachment_type=" + ToString(data["attachment_type"])
	data["isfile"] = 1
	str, err := self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  上传临时素材
 */
func (self *WorkWechat) Upload(data map[string]interface{}) (string, error) {
	url := API_URL_PREFIX + MEDIA_UPLOAD + self.AccessToken + "&type=" + ToString(data["type"])
	data["isfile"] = 1
	str, err := self.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  上传图片
 */
func (self *WorkWechat) Uploadimg(data map[string]interface{}) (string, error) {
	url := API_URL_PREFIX + UPLOADIMG + self.AccessToken
	data["isfile"] = 1
	str, err := self.sendForm("POST", url, data, data)
	return str, err
}
