package spiders

import (
	"fmt"
	"net/http"
	"testing"
)


func TestIQIYIRequest_Request(t *testing.T) {
	test := IQIYIRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "https://www.iqiyi.com/v_19rxwminv0.html?vfrm=pcw_home&vfrmblk=G&vfrmrst=712211_zongyi_image4",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

	for _, x := range test.Response{
		fmt.Println(x.Quality, x.Size, x.StreamType)
	}

}

