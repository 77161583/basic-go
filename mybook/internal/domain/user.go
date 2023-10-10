package domain

import "time"

// User 领域对象， 是 DDD 中的entirely
// BO(business object)
type User struct {
	Id           int64
	Email        string
	Phone        string
	Password     string
	NickName     string
	Birthday     string
	Introduction string
	CreateTime   time.Time
	UpdateTime   time.Time
}
