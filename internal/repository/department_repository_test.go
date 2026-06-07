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

func TestDepartmentCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	dept := &domain.Department{CompanyID: uuid.New(), DepartmentName: "IT", IsActive: true}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	repo := NewDepartmentRepository(gormDB)
	err := repo.Create(context.Background(), dept)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentCreate_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	dept := &domain.Department{CompanyID: uuid.New(), DepartmentName: "IT"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewDepartmentRepository(gormDB)
	err := repo.Create(context.Background(), dept)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentCreateTx_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	dept := &domain.Department{CompanyID: uuid.New(), DepartmentName: "HR"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	tx := gormDB.Begin()
	repo := NewDepartmentRepository(gormDB)
	err := repo.CreateTx(tx, context.Background(), dept)
	tx.Commit()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentCreateTx_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	dept := &domain.Department{CompanyID: uuid.New(), DepartmentName: "HR"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	tx := gormDB.Begin()
	repo := NewDepartmentRepository(gormDB)
	err := repo.CreateTx(tx, context.Background(), dept)
	tx.Rollback()

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentCreateBulk_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	depts := []domain.Department{
		{CompanyID: uuid.New(), DepartmentName: "IT"},
		{CompanyID: uuid.New(), DepartmentName: "HR"},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(uuid.New()).AddRow(uuid.New()))
	mock.ExpectCommit()

	repo := NewDepartmentRepository(gormDB)
	err := repo.CreateBulk(context.Background(), depts)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentCreateBulk_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	depts := []domain.Department{{CompanyID: uuid.New(), DepartmentName: "IT"}}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "departments"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewDepartmentRepository(gormDB)
	err := repo.CreateBulk(context.Background(), depts)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentGetByCompanyID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "company_id", "department_name", "is_active"}).
		AddRow(uuid.New(), companyID, "IT", true).
		AddRow(uuid.New(), companyID, "HR", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "departments"`)).
		WithArgs(companyID).
		WillReturnRows(rows)

	repo := NewDepartmentRepository(gormDB)
	result, err := repo.GetByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentGetByCompanyID_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "departments"`)).
		WithArgs(companyID).
		WillReturnError(gorm.ErrInvalidDB)

	repo := NewDepartmentRepository(gormDB)
	result, err := repo.GetByCompanyID(context.Background(), companyID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentUpdateStatus_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "departments"`)).
		WithArgs(true, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewDepartmentRepository(gormDB)
	err := repo.UpdateStatus(context.Background(), id, true)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentUpdateStatus_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "departments"`)).
		WithArgs(true, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewDepartmentRepository(gormDB)
	err := repo.UpdateStatus(context.Background(), id, true)

	assert.Error(t, err)
	assert.Equal(t, "department not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentUpdateStatus_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "departments"`)).
		WithArgs(true, sqlmock.AnyArg(), id).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewDepartmentRepository(gormDB)
	err := repo.UpdateStatus(context.Background(), id, true)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentDelete_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "departments"`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewDepartmentRepository(gormDB)
	err := repo.Delete(context.Background(), id)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDepartmentDelete_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "departments"`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewDepartmentRepository(gormDB)
	err := repo.Delete(context.Background(), id)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
