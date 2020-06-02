package spider

import (
	"testing"
)

func TestYOUKUClient_Request(t *testing.T) {
	var jar CookiesJar
	jar.SetValue("cna", "FOaSFXkoQE8CAb4CjttHruGf")
	jar.SetValue("P_pck_rm", "4swDIg3O795b6f4acd00e5ZBcrAoyMq9ceW6ThZE3ie%2FjN2ZhcJgd%2BRmXlbPDD8XTwvH3TkRS110kqGQCR36sdZQR0yK4NvelbRrtM1HrOVkMaKUuRTeXL4JcVH6qF9bUbRG9UlSmE2cyk5VGap0coutj%2BCtkI91EOI6Eg%3D%3D%5FV2")
	test := YOUKUClient{}
	test.Initialization("https://v.youku.com/v_show/id_XNDY1NDkyNTYxMg==.html", jar)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
