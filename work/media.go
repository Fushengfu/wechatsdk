package work

type Media struct {
}

/**
 *  上传附件资源
 */
func (w *WorkWechat) UploadAttachment(data map[string]interface{}) (string, error) {
	url := API_URL_PREFIX + UPLOAD_ATTACHMENT + w.AccessToken + "&media_type=" + ToString(data["media_type"]) + "&attachment_type=" + ToString(data["attachment_type"])
	data["isfile"] = 1
	str, err := w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  上传临时素材
 */
func (w *WorkWechat) Upload(data map[string]interface{}) (string, error) {
	url := API_URL_PREFIX + MEDIA_UPLOAD + w.AccessToken + "&type=" + ToString(data["type"])
	data["isfile"] = 1
	str, err := w.sendForm("POST", url, data, data)
	return str, err
}

/**
 *  上传图片
 */
func (w *WorkWechat) Uploadimg(data map[string]interface{}) (string, error) {
	url := API_URL_PREFIX + UPLOADIMG + w.AccessToken
	data["isfile"] = 1
	str, err := w.sendForm("POST", url, data, data)
	return str, err
}
