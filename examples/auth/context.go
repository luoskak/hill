package auth

import "context"

func GetUser(ctx context.Context) *User {
	i := ctx.Value(mwKey)
	return i.(*User)
}
