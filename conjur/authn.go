package conjur

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/jmcvetta/napping"
)

/**
 * Takes a username and password and returns an api_key to be used for authentication
 */
func (c *ConjurClientImpl) Login(username, password string) string {
	const LoginPath = "/api/authn/users/login"

	s := *c.session // copy the session to modify it
	s.Userinfo = url.UserPassword(username, password)
	s.Header = &http.Header{}
	s.Header.Set("Accept", "text/plain")
	//s.Log = true

	var res string
	var e ErrMsg

	c.log("Logging in: ", username)
	resp, err := s.Get(c.host+LoginPath, nil, &res, &e)
	// err is expected because it isn't JSON
	if resp == nil || resp.Status() != 200 {
		log.Fatal(err)
	}

	api_key := resp.RawText()
	c.log("Login success: ", api_key, " ", len(api_key))

	c.username = username
	c.api_key = api_key

	return api_key
}

func (c *ConjurClientImpl) Authenticate() string {
	const AuthenticatePath = "/api/authn/users/%s/authenticate"

	path := fmt.Sprintf(AuthenticatePath, c.username)
	var res string
	var e ErrMsg

	r := napping.Request{
		Method:     "POST",
		Url:        c.host + path,
		Payload:    bytes.NewBufferString(c.api_key),
		Result:     res,
		Error:      e,
		RawPayload: true,
	}

	log.Print("Authenticating: ", c.username)
	c.session.Header = &http.Header{}
	c.session.Header.Set("Content-Type", "text/plain")
	c.session.Header.Set("Accept", "text/plain")
	//c.session.Log = true
	resp, err := c.session.Send(&r)

	// err is expected beause it isn't JSON
	if resp == nil || resp.Status() != 200 {
		log.Fatal(err)
	}

	token := resp.RawText()
	c.log("Authenticate success: ", token)

	c.token = token
	c.session.Header.Del("Content-Type")
	c.session.Header.Del("Accept")

	return token
}
