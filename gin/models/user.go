package models

import (
	"context"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) Query(ctx context.Context, id int) error {
	return db.WithContext(ctx).First(u, id).Error
}

func (u *User) Insert(ctx context.Context) error {
	return db.WithContext(ctx).Create(u).Error
}
