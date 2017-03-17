package digest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
	"net/http/cookiejar"
	"strings"
)

type Request struct {
	Body          string
	Method        string
	Password      string
	Uri           string
	Username      string
	Auth          *authorization
	Wa            *wwwAuthenticate
	ContentType   string
	SkipTLSVerify bool
	client        *http.Client
}

func NewRequest(username string, password string, method string, uri string, body string) Request {
	dr := Request{}
	dr.UpdateRequest(username, password, method, uri, body)
	return dr
}

func (dr *Request) Close() {
}

func (dr *Request) UpdateRequest(username string,
	password string, method string, uri string, body string) *Request {

	dr.Body = body
	dr.Method = method
	dr.Password = password
	dr.Uri = uri
	dr.Username = username
	return dr
}

func (dr *Request) generateClient(timeout time.Duration) *http.Client {
	if dr.client == nil {
		cookieJar, _ := cookiejar.New(nil)
		if dr.SkipTLSVerify {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			dr.client = &http.Client{
				Jar:       cookieJar,
				Timeout:   timeout,
				Transport: tr,
			}
		} else {
			dr.client = &http.Client{Timeout: timeout, Jar: cookieJar}
		}
	}
	return dr.client
}

func (dr *Request) Execute() (resp *http.Response, err error) {
	resp, err = dr.auth()
	resp, err = dr.executeExistingDigest()
	return resp, err
}

func (dr *Request) ExecuteQuery(query string) (resp *http.Response, err error) {
	resp, err = dr.auth()
	newReq := *dr

	if strings.Contains(newReq.Uri, "?") {
		newReq.Uri = newReq.Uri + query
	} else {
		newReq.Uri = newReq.Uri + "?" + query
	}
	resp, err = newReq.executeExistingDigest()

	return resp, err
}

func (dr *Request) auth() (resp *http.Response, err error) {
	if dr.Auth == nil {
		var req *http.Request
		if req, err = http.NewRequest(dr.Method, dr.Uri, bytes.NewReader([]byte(dr.Body))); err != nil {
			return nil, err
		}

		if dr.ContentType != "" {
			req.Header.Set("Content-Type", dr.ContentType)
		}

		client := dr.generateClient(30 * time.Second)
		resp, err = client.Do(req)

		if resp.StatusCode == 401 {
			return dr.executeNewDigest(resp)
		}
		return
	}
	return resp, err
}

func (dr *Request) executeNewDigest(resp *http.Response) (*http.Response, error) {
	var (
		auth *authorization
		err  error
		wa   *wwwAuthenticate
	)

	waString := resp.Header.Get("WWW-Authenticate")
	if waString == "" {
		return nil, fmt.Errorf("Failed to get WWW-Authenticate header, please check your server configuration.")
	}
	wa = newWwwAuthenticate(waString)
	dr.Wa = wa

	if auth, err = newAuthorization(dr); err != nil {
		return nil, err
	}
	authString := auth.toString()

	if resp, err := dr.executeRequest(authString); err != nil {
		return nil, err
	} else {
		dr.Auth = auth
		return resp, nil
	}
}

func (dr *Request) executeExistingDigest() (*http.Response, error) {
	var (
		auth *authorization
		err  error
	)

	if auth, err = dr.Auth.refreshAuthorization(dr); err != nil {
		return nil, err
	}
	dr.Auth = auth

	authString := dr.Auth.toString()
	return dr.executeRequest(authString)
}

func (dr *Request) executeRequest(authString string) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)

	if req, err = http.NewRequest(dr.Method, dr.Uri, bytes.NewReader([]byte(dr.Body))); err != nil {
		return nil, err
	}

	// fmt.Printf("AUTHSTRING: %s\n\n", authString)
	req.Header.Add("Authorization", authString)

	if dr.ContentType != "" {
		req.Header.Set("Content-Type", dr.ContentType)
	}

	client := dr.generateClient(30 * time.Second)
	return client.Do(req)
}
