package cookie

//
//import (
//	"bytes"
//	"crypto/aes"
//	"crypto/cipher"
//	"database/sql"
//	"encoding/base64"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"os"
//	"syscall"
//	"unsafe"
//)
//
//var (
//	dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
//	dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")
//
//	pCryptUnprotectData = dllcrypt32.NewProc("CryptUnprotectData")
//	pLocalFree          = dllkernel32.NewProc("LocalFree")
//)
//
//var aesKey []byte
//
//type blob struct {
//	cbData uint32
//	pbData *byte
//}
//
//func newBlob(d []byte) *blob {
//	if len(d) == 0 {
//		return &blob{}
//	}
//	return &blob{
//		pbData: &d[0],
//		cbData: uint32(len(d)),
//	}
//}
//
//func (b *blob) toByteArray() []byte {
//	d := make([]byte, b.cbData)
//	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
//	return d
//}
//
//func getAesGCMKey() []byte {
//
//	var encryptedKey []byte
//	var path, _ = os.UserCacheDir()
//	var localStateFile = fmt.Sprintf("%s\\Google\\Chrome\\User Data\\Local State", path)
//
//	data, _ := ioutil.ReadFile(localStateFile)
//	var localState map[string]interface{}
//	_ = json.Unmarshal(data, &localState)
//
//	if localState["os_crypt"] != nil {
//
//		encryptedKey, _ = base64.StdEncoding.DecodeString(localState["os_crypt"].(map[string]interface{})["encrypted_key"].(string))
//
//		if bytes.Equal(encryptedKey[0:5], []byte{'D', 'P', 'A', 'P', 'I'}) {
//			encryptedKey, _ = decryptValue(encryptedKey[5:])
//		} else {
//			fmt.Print("encrypted_key does not look like DPAPI key\n")
//		}
//	}
//
//	return encryptedKey
//}
//
//func decryptValue(data []byte) ([]byte, error) {
//
//	if bytes.Equal(data[0:3], []byte{'v', '1', '0'}) {
//
//		aesBlock, _ := aes.NewCipher(aesKey)
//		aesgcm, _ := cipher.NewGCM(aesBlock)
//
//		nonce := data[3:15]
//		encryptedData := data[15:]
//
//		plaintext, _ := aesgcm.Open(nil, nonce, encryptedData, nil)
//
//		return plaintext, nil
//
//	} else {
//
//		var outblob blob
//		r, _, err := pCryptUnprotectData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
//		if r == 0 {
//			return nil, err
//		}
//		defer pLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
//		return outblob.toByteArray(), nil
//	}
//}
//
//type chromeDump struct {
//	Dump
//}
//
//func (d *chromeDump)Run()([]*http.Cookie, error)  {
//
//	db, err := sql.Open("sqlite3", d.file)
//	if err != nil{
//		return nil, err
//	}
//
//	rows, err := db.Query(`SELECT name, encrypted_value, host_key, is_httponly FROM cookies`)
//	if err != nil{
//		return nil, err
//	}
//
//	aesKey = getAesGCMKey()
//	var cookies []*http.Cookie
//	for rows.Next(){
//		var name, host string
//		var isHttpOnly int
//		var value []byte
//		err = rows.Scan(&name, &value, &host, &isHttpOnly)
//		if err != nil{
//			return nil, err
//		}
//
//		var isHttpOk bool
//		if isHttpOnly == 1 {
//			isHttpOk = true
//		}
//
//		value, err := decryptValue(value)
//		if err != nil{
//			return nil, err
//		}
//
//		cookies = append(cookies, &http.Cookie{
//			Name:       name,
//			Value:      string(value),
//			Domain:     host,
//			HttpOnly:   isHttpOk,
//		})
//	}
//
//	return cookies, nil
//}
