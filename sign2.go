package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"time"
)

// AWS signature v. 2
func Sign2(req *http.Request, keys *Keys) {
	if req.URL.Path == "" {
		req.URL.Path = "/"
	}

	q := req.URL.Query()
	q.Set("AWSAccessKeyId", keys.AccessKey)
	q.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))
	q.Set("SignatureVersion", "2")
	q.Set("SignatureMethod", "HmacSHA256")

	h := hmac.New(sha256.New, []byte(keys.SecretKey))
	h.Write([]byte(req.Method + "\n"))
	h.Write([]byte(req.URL.Host + "\n"))
	h.Write([]byte(req.URL.Path + "\n"))

	var pairs []string
	for k, _ := range q {
		for _, v := range q[k] {
			pairs = append(pairs, urlEncode(k)+"="+urlEncode(v))
		}
	}
	sort.Strings(pairs)
	h.Write([]byte(strings.Join(pairs, "&")))

	q.Set("Signature", base64.StdEncoding.EncodeToString(h.Sum(nil)))

	req.URL.RawQuery = q.Encode()
}

func urlEncode(str string) string {
	inLen, outLen := len(str), 0

	for i := 0; i < inLen; i++ {
		c := str[i]
		if isUnreserved(c) {
			outLen += 1
		} else {
			outLen += 3
		}
	}

	if inLen == outLen {
		return str
	}

	out := make([]byte, outLen)
	for i, off := 0, 0; i < inLen; i++ {
		c := str[i]
		if isUnreserved(c) {
			out[off] = c
			off++
		} else {
			out[off] = '%'
			out[off+1] = "0123456789ABCDEF"[c>>4]
			out[off+2] = "0123456789ABCDEF"[c&15]
			off += 3
		}
	}
	return string(out)
}

func isUnreserved(c byte) bool {
	return c >= 0x41 && c <= 0x5A || c >= 0x61 && c <= 0x7A || //ALPHA
		c >= 0x30 && c <= 0x39 || // DIGITS
		c == 0x2D || c == 0x2E || c == 0x5F || c == 0x7E // - . _ ~
}
