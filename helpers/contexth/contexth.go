package contexth

import "context"

// GetAuthUserId is a helper function to get the current logged in user id from the given context
func GetAuthUserId(ctx context.Context) string {
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
