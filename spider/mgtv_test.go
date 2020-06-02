package spider

import (
	"testing"
)

func TestMGTVClient_Request(t *testing.T) {
	var jar CookiesJar
	jar.SetValue("PM_CHKID", "73949fc6260ace04")
	test := MGTVClient{}
	test.Initialization("https://www.mgtv.com/b/328268/5144766.html", jar)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
