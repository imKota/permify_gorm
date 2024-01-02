package repositories

import (
	"gorm.io/gorm"

	"github.com/imKota/permify_gorm/collections"
	"github.com/imKota/permify_gorm/models"
	"github.com/imKota/permify_gorm/models/pivot"
	"github.com/imKota/permify_gorm/repositories/scopes"
)

// IRoleRepository its data access layer abstraction of role.
type IRoleRepository interface {
	Migratable

	// single fetch options

	GetRoleByID(ID uint) (role models.Role, err error)
	GetRoleByIDWithPermissions(ID uint) (role models.Role, err error)

	GetRoleByGuardName(guardName string) (role models.Role, err error)
	GetRoleByGuardNameWithPermissions(guardName string) (role models.Role, err error)

	// Multiple fetch options

	GetRoles(roleIDs []uint) (roles collections.Role, err error)
	GetRolesWithPermissions(roleIDs []uint) (roles collections.Role, err error)

	GetRolesByGuardNames(guardNames []string) (roles collections.Role, err error)
	GetRolesByGuardNamesWithPermissions(guardNames []string) (roles collections.Role, err error)

	// ID fetch options

	GetRoleIDs(pagination scopes.GormPager) (roleIDs []uint, totalCount int64, err error)
	GetRoleIDsOfUser(userID uint, pagination scopes.GormPager) (roleIDs []uint, totalCount int64, err error)
	GetRoleIDsOfPermission(permissionID uint, pagination scopes.GormPager) (roleIDs []uint, totalCount int64, err error)

	// FirstOrCreate & Updates & Delete

	FirstOrCreate(role *models.Role) (err error)
	Updates(role *models.Role, updates map[string]interface{}) (err error)
	Delete(role *models.Role) (err error)

	// Actions

	AddPermissions(role *models.Role, permissions collections.Permission) (err error)
	ReplacePermissions(role *models.Role, permissions collections.Permission) (err error)
	RemovePermissions(role *models.Role, permissions collections.Permission) (err error)
	ClearPermissions(role *models.Role) (err error)

	// Controls

	HasPermission(roles collections.Role, permission models.Permission) (b bool, err error)
	HasAllPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error)
	HasAnyPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error)
}

// RoleRepository its data access layer of role.
type RoleRepository struct {
	Database *gorm.DB
}

// Migrate generate tables from the database.
// @return error
func (repository *RoleRepository) Migrate() (err error) {
	err = repository.Database.AutoMigrate(models.Role{})
	err = repository.Database.AutoMigrate(pivot.UserRoles{})
	return
}

// SINGLE FETCH OPTIONS

// GetRoleByID get role by id.
// @param uint
// @return models.Role, error
func (repository *RoleRepository) GetRoleByID(ID uint) (role models.Role, err error) {
	err = repository.Database.First(&role, "roles.id = ?", ID).Error
	return
}

// GetRoleByIDWithPermissions get role by id with its permissions.
// @param uint
// @return models.Role, error
func (repository *RoleRepository) GetRoleByIDWithPermissions(ID uint) (role models.Role, err error) {
	err = repository.Database.Preload("Permissions").First(&role, "roles.id = ?", ID).Error
	return
}

// GetRoleByGuardName get role by guard name.
// @param string
// @return models.Role, error
func (repository *RoleRepository) GetRoleByGuardName(guardName string) (role models.Role, err error) {
	err = repository.Database.Where("roles.guard_name = ?", guardName).First(&role).Error
	return
}

// GetRoleByGuardNameWithPermissions get role by guard name with its permissions.
// @param string
// @return models.Role, error
func (repository *RoleRepository) GetRoleByGuardNameWithPermissions(guardName string) (role models.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.guard_name = ?", guardName).First(&role).Error
	return
}

// MULTIPLE FETCH OPTIONS

// GetRoles get roles by ids.
// @param []uint
// @return collections.Role, error
func (repository *RoleRepository) GetRoles(IDs []uint) (roles collections.Role, err error) {
	err = repository.Database.Where("roles.id IN (?)", IDs).Find(&roles).Error
	return
}

// GetRolesWithPermissions get roles by ids with its permissions.
// @param []uint
// @return collections.Role, error
func (repository *RoleRepository) GetRolesWithPermissions(IDs []uint) (roles collections.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.id IN (?)", IDs).Find(&roles).Error
	return
}

// GetRolesByGuardNames get roles by guard names.
// @param []string
// @return collections.Role, error
func (repository *RoleRepository) GetRolesByGuardNames(guardNames []string) (roles collections.Role, err error) {
	err = repository.Database.Where("roles.guard_name IN (?)", guardNames).Find(&roles).Error
	return
}

// GetRolesByGuardNamesWithPermissions get roles by guard names.
// @param []string
// @return collections.Role, error
func (repository *RoleRepository) GetRolesByGuardNamesWithPermissions(guardNames []string) (roles collections.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.guard_name IN (?)", guardNames).Find(&roles).Error
	return
}

// ID FETCH OPTIONS

// GetRoleIDs get role ids. (with pagination)
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *RoleRepository) GetRoleIDs(pagination scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Model(&models.Role{}).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("roles.id", &roleIDs).Error
	return
}

// GetRoleIDsOfUser get role ids of user. (with pagination)
// @param uint
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *RoleRepository) GetRoleIDsOfUser(userID uint, pagination scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("user_roles.role_id", &roleIDs).Error
	return
}

// GetRoleIDsOfPermission get role ids of permission. (with pagination)
// @param uint
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *RoleRepository) GetRoleIDsOfPermission(permissionID uint, pagination scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("role_permissions").Where("role_permissions.permission_id = ?", permissionID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("role_permissions.role_id", &roleIDs).Error
	return
}

// FirstOrCreate & Updates & Delete

// FirstOrCreate create new role if name not exist.
// @param *models.Role
// @return error
func (repository *RoleRepository) FirstOrCreate(role *models.Role) error {
	return repository.Database.Where(models.Role{GuardName: role.GuardName}).FirstOrCreate(role).Error
}

// Updates update role.
// @param *models.Role
// @param map[string]interface{}
// @return error
func (repository *RoleRepository) Updates(role *models.Role, updates map[string]interface{}) (err error) {
	return repository.Database.Model(role).Updates(updates).Error
}

// Delete delete role.
// @param *models.Role
// @return error
func (repository *RoleRepository) Delete(role *models.Role) (err error) {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_roles.role_id = ?", role.ID).Delete(&pivot.UserRoles{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Delete(role).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

// ACTIONS

// AddPermissions add permissions to role.
// @param *models.Role
// @param collections.Permission
// @return error
func (repository *RoleRepository) AddPermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Append(permissions.Origin())
}

// ReplacePermissions replace permissions of role.
// @param *models.Role
// @param collections.Permission
// @return error
func (repository *RoleRepository) ReplacePermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Replace(permissions.Origin())
}

// RemovePermissions remove permissions of role.
// @param *models.Role
// @param collections.Permission
// @return error
func (repository *RoleRepository) RemovePermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Delete(permissions.Origin())
}

// ClearPermissions remove all permissions of role.
// @param *models.Role
// @return error
func (repository *RoleRepository) ClearPermissions(role *models.Role) (err error) {
	return repository.Database.Model(role).Association("Permissions").Clear()
}

// Controls

// HasPermission does the role or any of the roles have given permission?
// @param collections.Role
// @param models.Permission
// @return bool, error
func (repository *RoleRepository) HasPermission(roles collections.Role, permission models.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id = ?", permission.ID).Count(&count).Error
	return count > 0, err
}

// HasAllPermissions does the role or roles have all the given permissions?
// @param collections.Role
// @param collections.Permission
// @return bool, error
func (repository *RoleRepository) HasAllPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return roles.Len()*permissions.Len() == count, err
}

// HasAnyPermissions does the role or roles have any of the given permissions?
// @param collections.Role
// @param collections.Permission
// @return bool, error
func (repository *RoleRepository) HasAnyPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return count > 0, err
}

// paginate pagging if pagination option is true.
// @param repositories_scopes.GormPager
// @return func(db *gorm.DB) *gorm.DB
func (repository *RoleRepository) paginate(pagination scopes.GormPager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination != nil {
			db.Scopes(pagination.ToPaginate())
		}

		return db
	}
}
