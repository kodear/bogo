package spider

import (
	"testing"
)

func TestMGTVClient_Request(t *testing.T) {
	jar := CookiesJar{}
	jar.SetValue("PM_CHKID", "73949fc6260ace04")
	test := MGTVClient{}
	test.SetURL("https://www.mgtv.com/b/328268/5144766.html")
	test.Initialization(jar, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
