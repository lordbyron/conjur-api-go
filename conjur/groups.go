package conjur

type Group struct {
	Id         string `json:"id"`
	UserId     string `json:"userid,omitempty"`
	OwnerId    string `json:"ownerid,omitempty"`
	Gid        int    `json:"gidnumber,omitempty"`
	RoleId     string `json:"roleid,omitempty"`
	ResourceId string `json:"resource_identifier,omitempty"`
}

func (c *ConjurClientImpl) Group() Group {
	return Group{}
}
