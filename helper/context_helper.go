package helper

import (
	"context"
)

// CtxGetAuthUserId is a helper function to get the current logged in user id from the given context
func CtxGetAuthUserId(ctx context.Context) string {
	user, ok := ctx.Value("user").(map[string]any)
	if !ok {
		return ""
	}
	id, ok := user["id"].(string)
	if !ok {
		return ""
	}
	return id
}
