package route53

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const APIRoot = "https://route53.amazonaws.com/2013-04-01"

func (c *Conn) Sign(req *http.Request) {
	dateString := time.Now().UTC().Format(time.RFC1123)

	h := hmac.New(sha256.New, []byte(c.keys.SecretKey))
	h.Write([]byte(dateString))
	sig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req.Header.Set("Date", dateString)
	req.Header.Set("X-Amzn-Authorization", "AWS3-HTTPS AWSAccessKeyId="+c.keys.AccessKey+",Algorithm=HmacSHA256,Signature="+sig)
}

// Common R53 error response
type ErrorResponse struct {
	HttpCode int `xml:"-"`
	Type     string
	Code     string
	Message  string
}

func (e *ErrorResponse) String() string {
	return fmt.Sprintf("%s error #%d: %s (%s)", e.Type, e.HttpCode, e.Code, e.Message)
}
func (e *ErrorResponse) Error() string { return e.String() }

////////////////

func checkError(resp *http.Response) error {
	if resp.StatusCode/100 == 2 { // 2xx
		return nil
	}
	// читаем ответ с ошибкой
	e := &struct{ Error ErrorResponse }{}
	e.Error.HttpCode = resp.StatusCode
	xml.NewDecoder(resp.Body).Decode(e)
	return &e.Error
}

func unescape(s string) string {
	re := regexp.MustCompile(`\\[0-3][0-7][0-7]`)
	return re.ReplaceAllStringFunc(s, func(a string) string {
		b, _ := strconv.ParseUint(a[1:], 8, 8)
		return string([]byte{uint8(b)})
	})
}

func setIf(q url.Values, name, value string) {
	if value != "" {
		q.Set(name, value)
	}
}
