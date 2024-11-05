package repository

import (
	m "cqrsdemo/model"
	"log/slog"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) m.RoleOperations {
	return &RoleRepository{
		db: db,
	}
}

type Role struct {
	Id      int `gorm:"primarykey"`
	Service string
	Role    string
}

func (Role) TableName() string {
	return "role"
}

func (r *RoleRepository) GetRoleByMethod(service string) (string, error) {
	var role Role

	tx := r.db.Limit(1).Find(&role, "service = ?", service)
	if err := tx.Error; err != nil {
		slog.Error("eror getting role by method", "error", err)
		return "", err
	}

	return role.Role, nil
}

func (r *RoleRepository) SaveRole(role m.Role) (int, error) {
	slog.Info("init save role", "service", role.Service)
	roleEntity := Role{
		Role:    role.Role,
		Service: role.Service,
	}

	tx := r.db.Create(&roleEntity)
	if err := tx.Error; err != nil {
		slog.Error("eror saving role", "error", err)
		return -1, err
	}

	return roleEntity.Id, nil
}