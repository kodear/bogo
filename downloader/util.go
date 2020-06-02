package downloader

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func AESDecrypt(crypted, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

func getTestPath() string {
	home := os.Getenv("BOGO")
	return filepath.Join(home[1:len(home)-1], "test")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func urlJoin(BaseUrl, MediaUrl string) (fullUrl string) {
	if strings.HasPrefix(MediaUrl, "http") {
		fullUrl = MediaUrl
	} else if strings.HasPrefix(MediaUrl, "/") {
		urlParse, _ := url.Parse(BaseUrl)
		fullUrl = fmt.Sprintf("%s://%s/%s", urlParse.Scheme, urlParse.Host, MediaUrl)
	} else {
		urlSplit := strings.Split(BaseUrl, "/")
		urlSplit[len(urlSplit)-1] = MediaUrl
		fullUrl = strings.Join(urlSplit, "/")
	}

	return
}
