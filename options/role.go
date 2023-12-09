package options

import (
	"github.com/ai-psyche/permify_gorm/utils"
)

// RoleOption represents options when fetching roles.
type RoleOption struct {
	WithPermissions bool
	Pagination      *utils.Pagination
}
