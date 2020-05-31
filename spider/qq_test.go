package spider

import (
	"net/http"
	"strings"
	"testing"
)

func TestQQClient_Request(t *testing.T) {
	cookies := `eas_sid=51c5m6q1e3E0N7v3y9U041x8L6;ied_qq=o1371235590;LW_sid=O1m5A9F0d4r9r6y5A347T011P8;LW_uid=r1F5P6O1P3F097i3j8I971a4b2;o_cookie=641015302;pac_uid=1_1371235590;pgv_info=ssid=s8192065488;pgv_pvi=546016256;pgv_pvid=4011816612;ptcz=5ca9b91710ffecdd008620930c8e9412302cfbf4cd388d775aa0e851b217b501;ptui_loginuin=641015302;RK=tWgV8kmNRp;tvfe_boss_uuid=37d48a26dd911ac5;uid=99954938;bucket_id=9231006;login_remember=qq;lw_nick=zxin|0|http://thirdqq.qlogo.cn/g?b=oidb&k=ZeOVQpibicOAVbItZianN8mLQ&s=640&t=1562659615|0;main_login=qq;ptag=www_baidu_com|videolist:click;qq_head=http://thirdqq.qlogo.cn/g?b=oidb&k=ZeOVQpibicOAVbItZianN8mLQ&s=640&t=1562659615;qq_nick=zxin;qv_als=YmLo8soUp2Yj+BJoA11590881397RMFUwg==;ts_last=v.qq.com/x/cover/mzc00200rxjna4e/t00340a0l12.html;ts_refer=www.baidu.com/link;ts_uid=7310358453;tvfe_search_uid=7df9c0ad-96d9-4d23-abf8-30d3cb8886bb;video_guid=b916ba14c8a33f3e;video_platform=2;vqq_access_token=D5AA7CEB9FDB0A73BF999EF063F6752C;vqq_appid=101483052;vqq_openid=ED87B84A15549ED3E60D6D2F927314D6;vqq_vuserid=119772433;vqq_vusession=Kq9BVwbcmfScIQ4-kERz5Q..;ad_play_index=35;QQLivePCVer=50201512`
	var cookieJar []*http.Cookie
	for _, cookie := range strings.Split(cookies, ";") {
		c := strings.Split(cookie, "=")
		name := strings.TrimSpace(c[0])
		value := strings.TrimSpace(c[1])
		cookieJar = append(cookieJar, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}
	test := QQClient{
		Client{
			CookieJar: cookieJar,
			Header:    http.Header{},
			URL:       "https://v.qq.com/x/cover/mzc00200uig074r/u0034ilke49.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
