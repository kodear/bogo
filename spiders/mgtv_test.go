package spiders

import (
	"net/http"
	"strings"
	"testing"
)

func TestMGTVRequest_Request(t *testing.T) {
	cookies := `locale=CHN; WWW_LOCALE=CN; __random_seed=0.8574044248094317; mba_deviceid=708d7311-baf8-5750-5c0e-94178f757da7; mba_cxid_expiration=1584028800000; mba_cxid=8kpqp21m4cm; sessionid=1583986048561_8kpqp21m4cm; MQGUID=1237952993639518208; __MQGUID=1237952993639518208; pc_v6=v6; id=52531462; rnd=rnd; seqid=bpkrb5hlqhggi8ohjs20; uuid=e6c27eb1894245deaae5ae61e66f9958; vipStatus=3; wei=7e4799e8757e05a9b16bebe9da429d6c; wei2=439bXDGf69cfqDyTUfB3hwROpjABBhNctRRC4a7V2ehAG7lm63xQ1374hK74j9AvZTPEnfYSVsBe2hTBUdjVf9Y5Ue0dGThrwOsbFNrwNOOjpupcEjDaZju0g00MBJD%2F3wVFHHW3rVkKkGwMRKBWccA%2BoPpJ8CbNpncBnMuIfRmVFjIgKkwqcCoTqSM93VdtjVbb94W7RCkqo9A; HDCN=BPKRB5HPAHH7767E04BG-864042134; PM_CHKID=0b9a503289232a45; mba_sessionid=4c07928f-ee89-5745-dea4-5c05e56600d1; mba_last_action_time=1583993031036; beta_timer=1583993032019; lastActionTime=1583994066911`
	var cookieJar []*http.Cookie
	for _, cookie := range strings.Split(cookies, ";") {
		c := strings.Split(cookie, "=")
		name := c[0]
		value := c[1]
		cookieJar = append(cookieJar, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	test := MGTVRequest{
		SpiderRequest{
			Header:    http.Header{},
			URL:       "https://www.mgtv.com/b/328268/5144766.html",
			CookieJar: cookieJar,
		},
	}

	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
