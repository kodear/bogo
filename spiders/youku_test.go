package spiders

import (
	"net/http"
	"strings"
	"testing"
)

func TestYOUKUClient_Request(t *testing.T) {
	cookies := `__ysuid=1562657576484g7V; __arlft=1565336440; cna=0h+eFVk7jUoCAXFaImmt/21r; juid=01di2dk1imtr7; ysestep=2; yseidcount=14; ystep=37; isg=BMLCuNjUKPHZVTeQIXI9k22CEMgkk8at6xiFkgzb3jXuX2PZ9CGDvA3VC9sjFD5F; UM_distinctid=16d48c4ffe3a-00adba159cc9948-4c312373-1fa400-16d48c4ffe4bd9; user_name=%E5%BC%A0%E9%91%AB98222; P_gck=NA%7Cmiev12TDYqw%2BMFG%2F0wMF7w%3D%3D%7CNA%7C1583416335712; P_pck_rm=4swDIg3Od16f4da0ebcd41ZBcrAoyMq9ceW6Tgi1qSkOv9yDiByryR9CcZoB4pYd7TRz9RowUPfsPdwmKmeTsanRap%2BmbvPsSdZWfMz3aoFHJS9zsYRKguRpRvKztFWm4r%2FRAJpNU2yqs0w9wGYolmOemoKDjpn8dNaXdQ%3D%3D%5FV2; __ayft=1584165175448; __aysid=1584165175448SLN; __arpvid=1584165203432DQz7Le-1584165203451; __ayscnt=1; __aypstp=2; __ayspstp=2; _m_h5_tk=53e53e8aad73467d8df4e20722a88076_1584170215798; _m_h5_tk_enc=931ba391ccf3ec06d63d735b2543ceed; P_ck_ctl=A9B8329412002663DD5FF63853DA362B; _m_h5_c=6fc499fed6f4306c85bf88fdf8c4bb2f_1584175256158%3B8c2b69e9583c062a6607b3e613463fe4; __arycid=dv-3-00; __arcms=dv-3-00; __ayvstp=1; __aysvstp=1`
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

	test := YOUKUClient{
		Client{
			Header:    http.Header{},
			URL:       "https://v.youku.com/v_show/id_XNDY1NDkyNTYxMg==.html",
			CookieJar: cookieJar,
		},
	}

	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
