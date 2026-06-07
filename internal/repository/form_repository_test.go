package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFormCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	form := &domain.Form{CompanyID: uuid.New(), Title: "Bug", Description: "desc"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "forms"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	repo := NewFormRepository(gormDB)
	err := repo.Create(context.Background(), form)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFormCreate_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	form := &domain.Form{CompanyID: uuid.New(), Title: "Bug", Description: "desc"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "forms"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewFormRepository(gormDB)
	err := repo.Create(context.Background(), form)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFormGetByCompanyID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "company_id", "title", "description"}).
		AddRow(uuid.New(), companyID, "Form A", "desc A").
		AddRow(uuid.New(), companyID, "Form B", "desc B")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms"`)).
		WithArgs(companyID).
		WillReturnRows(rows)

	repo := NewFormRepository(gormDB)
	result, err := repo.GetByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFormGetByCompanyID_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms"`)).
		WithArgs(companyID).
		WillReturnError(gorm.ErrInvalidDB)

	repo := NewFormRepository(gormDB)
	result, err := repo.GetByCompanyID(context.Background(), companyID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFormGetFormCompanyID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()
	dateStr := "2025-01-15"

	rows := sqlmock.NewRows([]string{"id", "company_id", "title", "description"}).
		AddRow(uuid.New(), companyID, "Daily Form", "desc")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(dateStr, companyID).
		WillReturnRows(rows)

	repo := NewFormRepository(gormDB)
	result, err := repo.GetFormCompanyID(context.Background(), companyID, dateStr)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFormGetFormCompanyID_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()
	dateStr := "2025-01-15"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(dateStr, companyID).
		WillReturnError(gorm.ErrInvalidDB)

	repo := NewFormRepository(gormDB)
	result, err := repo.GetFormCompanyID(context.Background(), companyID, dateStr)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
