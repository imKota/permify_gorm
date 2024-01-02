package options

import (
	"github.com/imKota/permify_gorm/utils"
)

// PermissionOption represents options when fetching permissions.
type PermissionOption struct {
	Pagination *utils.Pagination
}
