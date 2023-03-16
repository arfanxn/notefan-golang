package contexth

import "context"

// GetAuthUser is a helper function to get the current logged in user from the given context
func GetAuthUser(ctx context.Context) (userMap map[string]any) {
	userMap, ok := ctx.Value("user").(map[string]any)
	if !ok {
		return
	}
	return
}

// GetAuthUserId is a helper function to get the current logged in user id from the given context
func GetAuthUserId(ctx context.Context) string {
	userMap := GetAuthUser(ctx)
	id, ok := userMap["id"].(string)
	if !ok {
		return ""
	}
	return id
}
