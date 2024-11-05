package repository

import (
	"github.com/fabian99m/cqrsdemo/model"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db, mock = newMockDb()
	roleRows = []string{"id", "service", "role"}
)

func TestGetRolesByService(t *testing.T) {
	role := "testrole"
	rows := sqlmock.NewRows(roleRows).AddRow(1, "testservice", role)

	mock.ExpectQuery("SELECT *").WillReturnRows(rows)
	underTest := NewRoleRepository(db)

	roleRepository, err := underTest.GetRoleByMethod("testservice")

	assert.Nil(t, err)
	assert.Equal(t, role, roleRepository)
}

func TestGetRolesByServiceError(t *testing.T) {
	err := fmt.Errorf("table dont exists")

	mock.ExpectQuery("SELECT *").WillReturnError(err)
	underTest := NewRoleRepository(db)

	_, errRepository := underTest.GetRoleByMethod("testservice")

	assert.True(t, errors.Is(errRepository, err))
}

func TestSaveRole(t *testing.T) {
	idKey := 21
	rows := sqlmock.NewRows([]string{"id"}).AddRow(idKey)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO *").WillReturnRows(rows)
	mock.ExpectCommit()

	underTest := NewRoleRepository(db)

	id, err := underTest.SaveRole(model.Role{
		Service: "Test",
		Role:    "roleTest",
	})

	assert.Nil(t, err)
	assert.Equal(t, idKey, id)
}

func TestSaveRoleError(t *testing.T) {
	err := fmt.Errorf("table dont exists")

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO *").WillReturnError(err)
	mock.ExpectRollback()

	underTest := NewRoleRepository(db)

	_, errRepository := underTest.SaveRole(model.Role{
		Service: "Test",
		Role:    "roleTest",
	})

	assert.True(t, errors.Is(errRepository, err))
}

func newMockDb() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: true,
			},
		)},
	)

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}
