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

func TestTicketCreate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	ticket := &domain.Ticket{FormID: uuid.New(), Title: "Bug", Message: "msg"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	repo := NewTicketRepository(gormDB)
	err := repo.Create(context.Background(), ticket)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketCreate_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	ticket := &domain.Ticket{FormID: uuid.New(), Title: "Bug", Message: "msg"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewTicketRepository(gormDB)
	err := repo.Create(context.Background(), ticket)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketCreateBulk_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	tickets := []domain.Ticket{
		{FormID: uuid.New(), Title: "T1"},
		{FormID: uuid.New(), Title: "T2"},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(uuid.New()).AddRow(uuid.New()))
	mock.ExpectCommit()

	repo := NewTicketRepository(gormDB)
	err := repo.CreateBulk(context.Background(), tickets)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketCreateBulk_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	tickets := []domain.Ticket{{FormID: uuid.New(), Title: "T1"}}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewTicketRepository(gormDB)
	err := repo.CreateBulk(context.Background(), tickets)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketGetByCompanyID_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "title", "message", "status"}).
		AddRow(uuid.New(), "T1", "msg1", "opened").
		AddRow(uuid.New(), "T2", "msg2", "closed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(companyID).
		WillReturnRows(rows)

	repo := NewTicketRepository(gormDB)
	result, err := repo.GetByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketGetByCompanyID_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	companyID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(companyID).
		WillReturnError(gorm.ErrInvalidDB)

	repo := NewTicketRepository(gormDB)
	result, err := repo.GetByCompanyID(context.Background(), companyID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketUpdateTicketStatus_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
		WithArgs(domain.TicketClosed, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewTicketRepository(gormDB)
	err := repo.UpdateTicketStatus(context.Background(), id, domain.TicketClosed)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketUpdateTicketStatus_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
		WithArgs(domain.TicketClosed, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewTicketRepository(gormDB)
	err := repo.UpdateTicketStatus(context.Background(), id, domain.TicketClosed)

	assert.Error(t, err)
	assert.Equal(t, "ticket not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketUpdateTicketStatus_Error(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
		WithArgs(domain.TicketClosed, sqlmock.AnyArg(), id).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	repo := NewTicketRepository(gormDB)
	err := repo.UpdateTicketStatus(context.Background(), id, domain.TicketClosed)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
