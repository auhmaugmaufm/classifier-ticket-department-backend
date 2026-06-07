package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestGetByEmail_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
		AddRow(companyID, "user@test.com", "hashed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "companies"`)).
		WithArgs("user@test.com", 1).
		WillReturnRows(rows)

	repo := NewCompanyRepository(gormDB)
	result, err := repo.GetByEmail(context.Background(), "user@test.com")

	assert.NoError(t, err)
	assert.Equal(t, "user@test.com", result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "companies"`)).
		WithArgs("ghost@test.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

	repo := NewCompanyRepository(gormDB)
	result, err := repo.GetByEmail(context.Background(), "ghost@test.com")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "companies"`)).
		WithArgs("user@test.com", 1).
		WillReturnError(gorm.ErrInvalidDB)

	repo := NewCompanyRepository(gormDB)
	result, err := repo.GetByEmail(context.Background(), "user@test.com")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(uuid.New(), "a@test.com").
		AddRow(uuid.New(), "b@test.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "companies"`)).
		WillReturnRows(rows)

	repo := NewCompanyRepository(gormDB)
	result, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_Empty(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "companies"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

	repo := NewCompanyRepository(gormDB)
	result, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "companies"`)).
		WillReturnError(gorm.ErrInvalidDB)

	repo := NewCompanyRepository(gormDB)
	result, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTx_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	company := &domain.Company{Email: "new@test.com", PasswordHash: "hashed"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "companies"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	tx := gormDB.Begin()
	repo := NewCompanyRepository(gormDB)
	err := repo.CreateTx(tx, context.Background(), company)

	tx.Commit()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTx_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	company := &domain.Company{Email: "new@test.com", PasswordHash: "hashed"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "companies"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	tx := gormDB.Begin()
	repo := NewCompanyRepository(gormDB)
	err := repo.CreateTx(tx, context.Background(), company)

	tx.Rollback()

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
