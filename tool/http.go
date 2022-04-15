package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"time"
)

type Http struct {
	UserAgent string
	Request *http.Request
	Response *http.Response
	Jar *cookiejar.Jar
	client *http.Client
	cookies []*http.Cookie
	contentType string
}

/**
 *  GET 请求
 */
func (this *Http) Get(uri string, querys map[string]string) (res string, err error) {
	if nil != querys {
		uri = this.HttpBuildQuery(uri, querys)
	}

	err = this.InitRquest("GET", uri, nil, "json")
	if err != nil {
		return "", err
	}

	return this.send()
}

/**
 *  POST 请求
 */
func (this *Http) Post(url string, body interface{},) (res string, err error) {
	err = this.InitRquest("POST", url, body, "json")
	if err != nil {
		return "", err
	}
	return this.send()
}

/**
 *  POST 请求
 */
func (this *Http) PostForm(url string, body interface{},) (res string, err error) {
	err = this.InitRquest("POST", url, body, "form")
	if err != nil {
		return "", err
	}
	return this.send()
}

/**
 *  生成请求参数
 */
func (this *Http) HttpBuildQuery(uri string, querys map[string]string) string {
	params := url.Values{}
	Url, _ := url.Parse(uri)

	for k, v := range querys {
		params.Set(k, v)
	}

	Url.RawQuery = params.Encode()
	uri = Url.String()

	return uri
}

/**
 *  生成请求参数
 */
func (this *Http) HttpBuildForm(querys interface{}) *bytes.Buffer {
	params := url.Values{}

	data := querys.(map[string]string)

	for k, v := range data {
		params.Set(k, v)
	}

	return bytes.NewBuffer([]byte(params.Encode()))
}

/**
 *  发送Http请求
 */
func (this *Http) send() (respone string, err error) {
	/**
	 *  初始化cookie
	 */
	if nil == this.Jar {
		this.Jar,_ = cookiejar.New(nil)
		this.Jar.SetCookies(this.Request.URL, this.cookies)
	}

	/**
	 *  实例化客户端对象
	 */
	this.InitClient()

	this.Response, err = this.client.Do(this.Request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	this.cookies = this.Jar.Cookies(this.Request.URL)


	defer this.Response.Body.Close()

	body, err := ioutil.ReadAll(this.Response.Body)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	respone = string(body)

	return respone, nil
}

/**
 *  获取cookie
 */
func (this *Http) GetCookie(key string) string {
	for _,v := range this.cookies {
		fmt.Println(url.QueryUnescape(v.Value))
		fmt.Println(v.Name)
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

/**
 *  实例化请求对象
 */
func (this *Http) InitRquest(method, uri string, args interface{}, contentType string) (er error) {
	this.contentType = "application/json;charset=UTF-8"

	if method == "POST" {
		var buffer *bytes.Buffer

		jsonStr := []byte{}
		if reflect.TypeOf(args).String() == "string" {
			jsonStr = []byte(args.(string))
		} else {
			jsonStr, _ = json.Marshal(args)
		}

		buffer = bytes.NewBuffer(jsonStr)

		if contentType == "form" {
			buffer = this.HttpBuildForm(args)
			this.contentType = "application/x-www-form-urlencoded"
			fmt.Println("buffer", buffer.String())
		}

		this.Request, er = http.NewRequest(method, uri, buffer)
	} else {
		this.Request, er = http.NewRequest(method, uri, nil)
	}

	this.setHeaders()

	return er
}

/**
 *  实例化客户端对象
 */
func (this *Http) InitClient()  {
	if nil == this.client {
		this.client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: this.Jar,
			Timeout: 120 * time.Second,
		}
	}
}

/**
 *  设置请求头信息
 */
func (this *Http) setHeaders()  {
	this.Request.Header.Set("User-Agent", this.getUserAgent())
	this.Request.Header.Set("Content-Type", this.contentType)
	this.Request.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"100\", \"Google Chrome\";v=\"100\"")
	this.Request.Header.Set("Accept"," text/html,application/xhtml+xml,application/xml;application/json;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	this.Request.Header.Set("Host", this.Request.Host)
}

/**
 *  随机获取ua
 */
func (this *Http) getUserAgent() string {

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:38.0) Gecko/20100101 Firefox/38.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; InfoPath.3; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	}

	if this.UserAgent == "" {
		this.UserAgent = userAgents[rand.Intn(9)]
	}

	return this.UserAgent
}