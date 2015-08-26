package conjur

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/jmcvetta/napping"
)

type Conjur interface {
	Login(username, password string) string
	Authenticate() string

	User(login string) User
	CreateUser(user User) (*User, error)
}

type ErrMsg struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
}

type ConjurClientImpl struct {
	verbose  bool
	host     string
	username string
	api_key  string
	token    string
	session  *napping.Session
}

func NewConjurClient(host string, verbose bool) Conjur {
	c := ConjurClientImpl{
		verbose: verbose,
		host:    host,
		session: &napping.Session{
			Client: &http.Client{
				Transport: &http.Transport{
					// TODO don't ignore TLS...
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
			//Log: verbose,
		},
	}
	c.testConnection()
	return &c
}

func (c *ConjurClientImpl) testConnection() {
	const InfoPath = "/api/info"
	res := struct {
		Account string
	}{}
	var e ErrMsg
	resp, err := c.session.Get(c.host+InfoPath, nil, &res, &e)
	if err != nil {
		log.Fatal(err)
	}
	c.log("Connection test", resp.Status(), res.Account)
}

func (c *ConjurClientImpl) log(args ...interface{}) {
	if c.verbose {
		log.Println(args...)
	}
}

func (c *ConjurClientImpl) send(r *napping.Request) (response *napping.Response, err error) {
	if c.api_key == "" {
		log.Fatal("Must log in with either password or api_key before fetching resources")
	}

	for retry := 0; retry < 1; retry++ {
		if c.token == "" {
			c.Authenticate()
		}
		b64token := base64.StdEncoding.EncodeToString([]byte(c.token))
		c.session.Header.Set("Authorization", `Token token="`+b64token+`"`)

		response, err = c.session.Send(r)
		if err != nil {
			c.token = ""
		} else {
			return
		}
	}
	return
}

func (c *ConjurClientImpl) Get(path string, p *napping.Params, result interface{}) (response *napping.Response, err error) {
	r := napping.Request{
		Method: "GET",
		Url:    c.host + path,
		Params: p,
		Result: result,
		Error:  ErrMsg{},
	}
	return c.send(&r)
}

func (c *ConjurClientImpl) Post(path string, p *napping.Params, payload, result interface{}) (response *napping.Response, err error) {
	r := napping.Request{
		Method:  "POST",
		Url:     c.host + path,
		Payload: payload,
		Result:  result,
		Error:   ErrMsg{},
	}
	return c.send(&r)
}
