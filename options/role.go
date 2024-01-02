package options

import (
	"github.com/imKota/permify_gorm/utils"
)

// RoleOption represents options when fetching roles.
type RoleOption struct {
	WithPermissions bool
	Pagination      *utils.Pagination
}
