package user

import (
	"github.com/yunling101/ControllerManager/common"
	"time"
)

type User struct {
	Id          int       `json:"id"`
	Sid         int       `json:"sid"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Nickname    string    `json:"nickname"`
	IsActive    bool      `json:"is_active"`
	IsSuperuser bool      `json:"is_superuser"`
	RoleId      int       `json:"role_id"`
	LastLogin   time.Time `json:"last_Login"`
}

func (t *User) TableName() string {
	return common.Config().Global.TableNamePrefix() + "_user"
}
