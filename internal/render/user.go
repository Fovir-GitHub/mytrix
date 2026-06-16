package render

import (
	"bytes"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/db"
)

func UserMarkdown(user *db.User) string {
	var buf bytes.Buffer
	if err := userTmpl.Execute(&buf, user); err != nil {
		return fmt.Sprintf("ID: %s, Role: %s", user.ID, user.Role)
	}
	return buf.String()
}
