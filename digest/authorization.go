package digest

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type authorization struct {
	Algorithm string
	Cnonce    string
	Nc        int
	Nonce     string
	Opaque    string
	Qop       string
	Realm     string
	Response  string
	Uri       string
	Userhash  bool
	Username  string
	Username_ string
}

func newAuthorization(dr *Request) (*authorization, error) {

	ah := authorization{
		Algorithm: dr.Wa.Algorithm,
		Cnonce:    "",
		Nc:        0,
		Nonce:     dr.Wa.Nonce,
		Opaque:    dr.Wa.Opaque,
		Qop:       "",
		Realm:     dr.Wa.Realm,
		Response:  "",
		Uri:       "",
		Userhash:  dr.Wa.Userhash,
		Username:  "",
		Username_: "", // TODO
	}

	return ah.refreshAuthorization(dr, true)
}

func (ah *authorization) refreshAuthorization(dr *Request, updateUri bool) (*authorization, error) {

	ah.Username = dr.Username

	if ah.Userhash {
		ah.Username = ah.hash(fmt.Sprintf("%v:%v", ah.Username, ah.Realm))
	}

	ah.Nc++

	ah.Cnonce = ah.hash(fmt.Sprintf("%d:%v:my_value", time.Now().UnixNano(), dr.Username))

	if updateUri {
		url, err := url.Parse(dr.Uri)
		if err != nil {
			return nil, err
		}
		ah.Uri = url.RequestURI()
	}

	ah.Response = ah.computeResponse(dr)

	return ah, nil
}

func (ah *authorization) computeResponse(dr *Request) (s string) {

	kdSecret := ah.hash(ah.computeA1(dr))
	kdData := fmt.Sprintf("%v:%08x:%v:%v:%v", ah.Nonce, ah.Nc, ah.Cnonce, ah.Qop, ah.hash(ah.computeA2(dr)))

	return ah.hash(fmt.Sprintf("%v:%v", kdSecret, kdData))
}

func (ah *authorization) computeA1(dr *Request) string {

	if ah.Algorithm == "" || ah.Algorithm == "MD5" || ah.Algorithm == "SHA-256" {
		return fmt.Sprintf("%v:%v:%s", ah.Username, ah.Realm, dr.Password)
	}

	if ah.Algorithm == "MD5-sess" || ah.Algorithm == "SHA-256-sess" {
		upHash := ah.hash(fmt.Sprintf("%v:%v:%s", ah.Username, ah.Realm, dr.Password))
		return fmt.Sprintf("%v:%v:%v", upHash, ah.Nonce, ah.Cnonce)
	}

	return ""
}

func (ah *authorization) computeA2(dr *Request) string {

	if matched, _ := regexp.MatchString("auth-int", dr.Wa.Qop); matched {
		ah.Qop = "auth-int"
		return fmt.Sprintf("%v:%v:%v", dr.Method, ah.Uri, ah.hash(dr.Body))
	}

	if dr.Wa.Qop == "auth" || dr.Wa.Qop == "" {
		ah.Qop = "auth"
		return fmt.Sprintf("%v:%v", dr.Method, ah.Uri)
	}

	return ""
}

func (ah *authorization) hash(a string) (s string) {

	var h hash.Hash

	if ah.Algorithm == "" || ah.Algorithm == "MD5" || ah.Algorithm == "MD5-sess" {
		h = md5.New()
	} else if ah.Algorithm == "SHA-256" || ah.Algorithm == "SHA-256-sess" {
		h = sha256.New()
	}

	io.WriteString(h, a)
	s = hex.EncodeToString(h.Sum(nil))

	return
}

func (ah *authorization) toString() string {
	var buffer bytes.Buffer

	buffer.WriteString("Digest ")

	if ah.Algorithm != "" {
		buffer.WriteString(fmt.Sprintf("algorithm=%v, ", ah.Algorithm))
	}

	if ah.Cnonce != "" {
		buffer.WriteString(fmt.Sprintf("cnonce=\"%v\", ", ah.Cnonce))
	}

	if ah.Nc != 0 {
		buffer.WriteString(fmt.Sprintf("nc=%08x, ", ah.Nc))
	}

	if ah.Opaque != "" {
		buffer.WriteString(fmt.Sprintf("opaque=\"%v\", ", ah.Opaque))
	}

	if ah.Nonce != "" {
		buffer.WriteString(fmt.Sprintf("nonce=\"%v\", ", ah.Nonce))
	}

	if ah.Qop != "" {
		buffer.WriteString(fmt.Sprintf("qop=%v, ", ah.Qop))
	}

	if ah.Realm != "" {
		buffer.WriteString(fmt.Sprintf("realm=\"%v\", ", ah.Realm))
	}

	if ah.Response != "" {
		buffer.WriteString(fmt.Sprintf("response=\"%v\", ", ah.Response))
	}

	if ah.Uri != "" {
		buffer.WriteString(fmt.Sprintf("uri=\"%v\", ", ah.Uri))
	}

	if ah.Userhash {
		buffer.WriteString("userhash=true, ")
	}

	if ah.Username != "" {
		buffer.WriteString(fmt.Sprintf("username=\"%v\", ", ah.Username))
	}

	s := buffer.String()

	return strings.TrimSuffix(s, ", ")
}
