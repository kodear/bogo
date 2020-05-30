package spiders

import (
	"net/http"
	"strings"
	"testing"
)

func TestQQClient_Request(t *testing.T) {
	cookies := `_ga=GA1.2.1762018758.1563275273; pgv_pvid=3791155328; pgv_pvi=2601128960; RK=sehd5mmMTr; ptcz=776fa43d26f5d5328f9020b1e225bd17aea8c369e729ca21f741932cf37ad92c; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2216c756b73502c2-0308a90f9b69d88-4c312272-2073600-16c756b7351a9a%22%2C%22%24device_id%22%3A%2216c756b73502c2-0308a90f9b69d88-4c312272-2073600-16c756b7351a9a%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%7D%7D; tvfe_boss_uuid=d236fd0662d67a42; video_guid=f95f19357edb465c; video_platform=2; ptui_loginuin=641015302; main_login=qq; vqq_access_token=D5AA7CEB9FDB0A73BF999EF063F6752C; vqq_appid=101483052; vqq_openid=ED87B84A15549ED3E60D6D2F927314D6; vqq_vuserid=119772433; vqq_vusession=-KVfOFOwyDcffaDKqgFceQ..; vqq_refresh_token=8EC653149DAF690299935F78CEBBFEFC; vqq_next_refresh_time=1399; vqq_login_time_init=1583483345; pgv_info=ssid=s2253475206; login_time_last=2020-3-6 16:29:7; uid=99954938`
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
			Header: http.Header{},
			URL:    "https://v.qq.com/x/cover/mtc2r2yum4j837e/s0034gfrspa.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
