package main

import (
	"encoding/base64"
	"fmt"

	"github.com/lordbyron/conjur"
)

const Conjur_url = "https://conjur"

func main() {
	c := conjur.NewConjurClient(Conjur_url, true)
	c.Login("admin", "boxrocks")
	//c.User("admin")
	t := c.Authenticate()
	b64token := base64.StdEncoding.EncodeToString([]byte(t))
	fmt.Println(b64token)

	/*
		c.CreateUser(conjur.User{
			Login: "user3",
		})
	*/
}
