package downloader

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
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
