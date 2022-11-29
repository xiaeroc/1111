package models

import (
	"encoding/base64"
	"encoding/json"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/buger/jsonparser"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ua2 = `okhttp/3.12.1;jdmall;android;version/10.1.2;build/89743;screen/1440x3007;os/11;network/wifi;`

type AutoGenerated struct {
	ClientVersion string `json:"clientVersion"`
	Client        string `json:"client"`
	Sv            string `json:"sv"`
	St            string `json:"st"`
	UUID          string `json:"uuid"`
	Sign          string `json:"sign"`
	FunctionID    string `json:"functionId"`
}

func getSign() *AutoGenerated {
	data, _ := httplib.Get("https://pan.smxy.xyz/sign").SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36").Bytes()
	t := &AutoGenerated{}
	json.Unmarshal(data, t)
	//logs.Info(t.Sign)
	if t != nil {
		t.FunctionID = "genToken"
	}
	return t
}

func GetWsKey1(wsKey string) (string, error) {
	v := url.Values{}
	s := getSign()
	v.Add("functionId", s.FunctionID)
	v.Add("clientVersion", s.ClientVersion)
	v.Add("client", s.Client)
	v.Add("uuid", s.UUID)
	v.Add("st", s.St)
	v.Add("sign", s.Sign)
	v.Add("sv", s.Sv)
	req := httplib.Post(`https://api.m.jd.com/client.action?` + v.Encode())
	req.Header("cookie", wsKey)
	req.Header("User-Agent", ua2)
	req.Header("content-type", `application/x-www-form-urlencoded; charset=UTF-8`)
	req.Header("charset", `UTF-8`)
	req.Header("accept-encoding", `br,gzip,deflate`)
	req.Body(`body=%7B%22action%22%3A%22to%22%2C%22to%22%3A%22https%253A%252F%252Fplogin.m.jd.com%252Fcgi-bin%252Fm%252Fthirdapp_auth_page%253Ftoken%253DAAEAIEijIw6wxF2s3bNKF0bmGsI8xfw6hkQT6Ui2QVP7z1Xg%2526client_type%253Dandroid%2526appid%253D879%2526appup_type%253D1%22%7D&`)
	data, err := req.Bytes()
	if err != nil {
		return "", err
	}
	//logs.Info(string(data))
	//logs.Info("获取token正常")
	tokenKey, _ := jsonparser.GetString(data, "tokenKey")
	cookie, err := appjmp(tokenKey)
	//logs.Info(cookie)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func appjmp(tokenKey string) (string, error) {
	v := url.Values{}
	v.Add("tokenKey", tokenKey)
	v.Add("to", ``)
	v.Add("client_type", "android")
	v.Add("appid", "879")
	v.Add("appup_type", "1")
	req := httplib.Get(`https://un.m.jd.com/cgi-bin/app/appjmp?` + v.Encode())
	req.Header("User-Agent", ua2)
	req.Header("accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3`)
	req.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})
	rsp, err := req.Response()
	if err != nil {
		return "", err
	}
	cookies := strings.Join(rsp.Header.Values("Set-Cookie"), " ")
	//ptKey := FetchJdCookieValue("pt_key", cookies)
	return cookies, nil
}
func GetWsKey(wsKey string) (string, error) {
	url := checkCloud()
	return getToken(url, wsKey)
}

func checkCloud() string {
	urlList := []string{"aHR0cDovLzQzLjEzNS45MC4yMy8=", "aHR0cHM6Ly9zaGl6dWt1Lm1sLw==", "aHR0cHM6Ly9jZi5zaGl6dWt1Lm1sLw=="}
	for i := range urlList {
		decodeString, _ := base64.StdEncoding.DecodeString(urlList[i])
		req := httplib.Get(string(decodeString))
		req.Header("User-Agent", "python-requests/2.25.1")
		s, err := req.String()
		//logs.Info(s, err)
		if strings.Contains(s, "200") && err == nil {
			return string(decodeString)
		}
	}
	decodeString, _ := base64.StdEncoding.DecodeString(urlList[0])
	return string(decodeString)
}

type T struct {
	Code      int    `json:"code"`
	Update    int    `json:"update"`
	Jdurl     string `json:"jdurl"`
	UserAgent string `json:"User-Agent"`
}
type T2 struct {
	FunctionId    string `json:"functionId"`
	ClientVersion string `json:"clientVersion"`
	Build         string `json:"build"`
	Client        string `json:"client"`
	Partner       string `json:"partner"`
	Oaid          string `json:"oaid"`
	SdkVersion    string `json:"sdkVersion"`
	Lang          string `json:"lang"`
	HarmonyOs     string `json:"harmonyOs"`
	NetworkType   string `json:"networkType"`
	Uemps         string `json:"uemps"`
	Ext           string `json:"ext"`
	Ef            string `json:"ef"`
	Ep            string `json:"ep"`
	St            int64  `json:"st"`
	Sign          string `json:"sign"`
	Sv            string `json:"sv"`
}

func cloudInfo(url string) string {
	// 重试10次
	for i := 0; i <= 10; i++ {
		req := httplib.Get(url + "check_api")
		req.Header("User-Agent", "python-requests/2.25.1")
		req.Header("authorization", "Bearer Shizuku")
		s, _ := req.Bytes()
		t := T{}
		json.Unmarshal(s, &t)
		logs.Info(t.UserAgent)
		if t.UserAgent != "" {
			return t.UserAgent
		}
		time.Sleep(time.Second * 2)
	}
	// 重试10次 还是失败 搞个默认值
	return "jdapp;android;10.3.5;;;appBuild/92468;ef/1;ep/{\"hdid\":\"JM9F1ywUPwflvMIpYPok0tt5k9kW4ArJEU3lfLhxBqw=\",\"ts\":1647918020020,\"ridx\":-1,\"cipher\":{\"sv\":\"CJS=\",\"ad\":\"EWOzY2ZuCNHsEWHvEJc3EG==\",\"od\":\"ENY0CJS3ZtvvCtK5ZJC5Yq==\",\"ov\":\"CzO=\",\"ud\":\"EWOzY2ZuCNHsEWHvEJc3EG==\"},\"ciphertype\":5,\"version\":\"1.2.0\",\"appname\":\"com.jingdong.app.mall\"};Mozilla/5.0 (Linux; Android 12; M2102K1C Build/SKQ1.211006.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/97.0.4692.98 Mobile Safari/537.36"
}
func getToken1(urls string, ua string) T2 {
	t := T2{}
	for i := 0; i <= 10; i++ {
		req1 := httplib.Get(urls + "genToken")
		req1.Header("User-Agent", ua)
		s, _ := req1.Bytes()
		json.Unmarshal(s, &t)
		if t.Sign != "" {
			return t
		}
		time.Sleep(time.Second * 2)
	}
	bytes := []byte("{\"functionId\":\"genToken\",\"clientVersion\":\"10.3.5\",\"build\":\"92468\",\"client\":\"android\",\"partner\":\"google\",\"oaid\":\"06lr6e4tz0n0hjsq\",\"sdkVersion\":\"31\",\"lang\":\"zh_CN\",\"harmonyOs\":\"0\",\"networkType\":\"UNKNOWN\",\"uemps\":\"0-2\",\"ext\":\"{\\\"prstate\\\": \\\"0\\\", \\\"pvcStu\\\": \\\"1\\\"}\",\"ef\":\"1\",\"ep\":\"{\\\"hdid\\\":\\\"JM9F1ywUPwflvMIpYPok0tt5k9kW4ArJEU3lfLhxBqw=\\\",\\\"ts\\\":1647918629476,\\\"ridx\\\":-1,\\\"cipher\\\":{\\\"d_model\\\":\\\"JWunCVVidRTr\\\",\\\"wifiBssid\\\":\\\"dW5hbw93bq==\\\",\\\"osVersion\\\":\\\"CJS=\\\",\\\"d_brand\\\":\\\"WQvrb21f\\\",\\\"screen\\\":\\\"CJuyCMenCNq=\\\",\\\"uuid\\\":\\\"oXVmbQ1uD3P6azYzczLfbG==\\\",\\\"aid\\\":\\\"oXVmbQ1uD3P6azYzczLfbG==\\\",\\\"openudid\\\":\\\"oXVmbQ1uD3P6azYzczLfbG==\\\"},\\\"ciphertype\\\":5,\\\"version\\\":\\\"1.2.0\\\",\\\"appname\\\":\\\"com.jingdong.app.mall\\\"}\",\"st\":1647918629476,\"sign\":\"48869070986ea43e1b9dd6c4aad60ca7\",\"sv\":\"120\"}")
	json.Unmarshal(bytes, &t)
	return T2{}
}

func getToken(urls string, wskey string) (string, error) {
	ua := cloudInfo(urls)
	t := getToken1(urls, ua)
	v := url.Values{}
	v.Add("functionId", t.FunctionId)
	v.Add("clientVersion", t.ClientVersion)
	v.Add("build", t.Build)
	v.Add("client", t.Client)
	v.Add("partner", t.Partner)
	v.Add("oaid", t.Oaid)
	v.Add("sdkVersion", t.SdkVersion)
	v.Add("lang", t.Lang)
	v.Add("harmonyOs", t.HarmonyOs)
	v.Add("networkType", t.NetworkType)
	v.Add("uemps", t.Uemps)
	v.Add("ext", t.Ext)
	v.Add("ef", t.Ef)
	v.Add("ep", t.Ep)
	v.Add("st", strconv.FormatInt(t.St, 10))
	v.Add("sign", t.Sign)
	v.Add("sv", t.Sv)
	req := httplib.Post(`https://api.m.jd.com/client.action?` + v.Encode())
	req.Header("cookie", wskey)
	req.Header("User-Agent", ua)
	req.Header("content-type", `application/x-www-form-urlencoded; charset=UTF-8`)
	req.Header("charset", `UTF-8`)
	req.Header("accept-encoding", `br,gzip,deflate`)
	req.Body(`body=%7B%22to%22%3A%22https%253a%252f%252fplogin.m.jd.com%252fjd-mlogin%252fstatic%252fhtml%252fappjmp_blank.html%22%7D&`)
	data, err := req.Bytes()
	if err != nil {
		return "", err
	}
	tokenKey, _ := jsonparser.GetString(data, "tokenKey")
	cookie, err := appjmp(tokenKey)
	//logs.Info(cookie)
	if err != nil {
		return "", err
	}
	return cookie, nil
}
