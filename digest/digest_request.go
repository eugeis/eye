package digest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Request struct {
	Body        string
	Method      string
	Password    string
	Uri         string
	Username    string
	Auth        *authorization
	Wa          *wwwAuthenticate
	ContentType string
}

func NewClient(skipTLSVerify bool, timeout time.Duration) (ret *http.Client) {
	cookieJar, _ := cookiejar.New(nil)
	if skipTLSVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		ret = &http.Client{
			Jar:       cookieJar,
			Timeout:   timeout,
			Transport: tr,
		}
	} else {
		ret = &http.Client{Timeout: timeout, Jar: cookieJar}
	}
	return
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

func (dr *Request) Execute(client *http.Client) (resp *http.Response, err error) {
	if dr.Auth == nil {
		var req *http.Request
		req, err = http.NewRequest(dr.Method, dr.Uri, bytes.NewReader([]byte(dr.Body)))
		if err == nil {
			if dr.ContentType != "" {
				req.Header.Set("Content-Type", dr.ContentType)
			}

			resp, err = client.Do(req)
			if err == nil && resp.StatusCode == 401 {
				resp, err = dr.executeNewDigest(resp, client)
			}
		}
	} else {
		resp, err = dr.executeExistingDigest(client)
		if err == nil && resp.StatusCode == 401 {
			//reset and start new request with auth
			resp.Body.Close()
			dr.Auth = nil
			resp, err = dr.Execute(client)
		}
	}
	return
}

func (dr *Request) executeNewDigest(resp *http.Response, client *http.Client) (*http.Response, error) {
	var (
		auth *authorization
		err  error
		wa   *wwwAuthenticate
	)

	waString := resp.Header.Get("WWW-Authenticate")

	resp.Body.Close()

	if waString == "" {
		return nil, fmt.Errorf("Failed to get WWW-Authenticate header, please check your server configuration.")
	}
	wa = newWwwAuthenticate(waString)
	dr.Wa = wa

	if auth, err = newAuthorization(dr); err != nil {
		return nil, err
	}
	authString := auth.toString()

	if resp, err := dr.executeRequest(authString, client); err != nil {
		return nil, err
	} else {
		dr.Auth = auth
		return resp, nil
	}
}

func (dr *Request) executeExistingDigest(client *http.Client) (*http.Response, error) {
	var (
		auth *authorization
		err  error
	)

	if auth, err = dr.Auth.refreshAuthorization(dr, true); err != nil {
		return nil, err
	}
	dr.Auth = auth

	authString := dr.Auth.toString()
	return dr.executeRequest(authString, client)
}

func (dr *Request) executeRequest(authString string, client *http.Client) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)

	if req, err = http.NewRequest(dr.Method, dr.Uri, bytes.NewReader([]byte(dr.Body))); err != nil {
		return nil, err
	}

	// fmt.Printf("AUTHSTRING: %v\n\n", authString)
	req.Header.Add("Authorization", authString)

	if dr.ContentType != "" {
		req.Header.Set("Content-Type", dr.ContentType)
	}
	return client.Do(req)
}
