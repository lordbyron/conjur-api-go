package conjur

import (
	"errors"
	"fmt"
)

type User struct {
	Login      string `json:"login"`
	UserId     string `json:"userid,omitempty"`
	OwnerId    string `json:"ownerid,omitempty"`
	Uid        int    `json:"uidnumber,omitempty"`
	RoleId     string `json:"roleid,omitempty"`
	ResourceId string `json:"resource_identifier,omitempty"`
	ApiKey     string `json:"api_key,omitempty"`
}

func (c *ConjurClientImpl) User(login string) User {
	const UserPath = "/api/users/%s"
	path := fmt.Sprintf(UserPath, login)

	var res User

	//c.session.Log = true
	c.log("Fetching user:", login)
	c.Get(path, nil, &res)
	c.log("User:", res)

	return res
}

func (c *ConjurClientImpl) CreateUser(user User) (*User, error) {
	const UserPath = "/api/users"

	if user.Login == "" {
		return nil, errors.New("cannot create user without at least a Login")
	}

	var res User

	//c.session.Log = true
	c.log("Creating user:", user)
	_, err := c.Post(UserPath, nil, user, &res)
	c.log("Created:", res)

	return &res, err
}
