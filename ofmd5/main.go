package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	response()
}

func strMd5(x string) string {
	h := md5.New()
	h.Write([]byte(x))                    // 需要加密的字符串为 sharejs.com
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果
}

func response() {
	User := "aaaaaa"
	A2Prefix := "AUTHENTICATE"
	Passwd := "111111"
	Realm := ""
	CNonce := "C6AC/KVrqlZPSKktfvbWQbom0RuKSavWTqnK8O72XKs="
	Nonce := "12916272820987184363"

	DigestURI := "xmpp/112.74.67.92"
	NC := "00000001"
	QOP := "auth"

	h := md5.New()
	h.Write([]byte(User + ":" + Realm + ":" + Passwd))
	MD5Hash := string(h.Sum(nil))

	A1 := MD5Hash + ":" + Nonce + ":" + CNonce
	A2 := A2Prefix + ":" + DigestURI
	T := strMd5(A1) + ":" + Nonce + ":" + NC + ":" + CNonce + ":" + QOP + ":" + strMd5(A2)
	fmt.Print(strMd5(T) + "\n")
}
