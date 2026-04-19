package cron

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/service"
)

type FormCron struct {
	FormService    *service.FormService
	CompanyService *service.CompanyService
	aiService      *service.AIService
	location       *time.Location
	stopCh         chan struct{}
	doneCh         chan struct{}
	onceStart      sync.Once
	onceStop       sync.Once
}

func NewFormCron(formService *service.FormService, companyService *service.CompanyService, aiService *service.AIService) *FormCron {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return &FormCron{
		FormService:    formService,
		CompanyService: companyService,
		aiService:      aiService,
		location:       loc,
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
	}
}

func (f *FormCron) Start() {
	f.onceStart.Do(func() {
		go f.runDailyJobLoop()
		log.Println("[FormCron] Cron scheduled at 12 AM (Asia/Bangkok)")
	})
}

func (f *FormCron) Stop() {
	f.onceStop.Do(func() {
		close(f.stopCh)
		<-f.doneCh
	})
}

func (f *FormCron) runDailyJobLoop() {
	defer close(f.doneCh)

	for {
		now := time.Now().In(f.location)
		// next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, f.location)
		// if !now.Before(next) {
		// 	next = next.Add(24 * time.Hour)
		// }
		next := now.Truncate(2 * time.Minute).Add(2 * time.Minute)

		timer := time.NewTimer(next.Sub(now))
		select {
		case <-timer.C:
			f.runJob()
		case <-f.stopCh:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
	}
}

func (f *FormCron) runJob() {
	log.Println("[FormCron] Starting job ....")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	companies, err := f.CompanyService.GetAllCompanies(ctx)
	if err != nil {
		log.Printf("[FormCron] fetch companies error: %v", err)
		return
	}

	var data []dto.CompanyFormItems
	lastDate := time.Now().In(f.location).AddDate(0, 0, -1).Format("2006-01-02")
	for _, company := range companies {
		forms, err := f.FormService.GetSubmitFormPerDayByCompanyID(ctx, company.ID, lastDate)
		if err != nil {
			log.Printf("[FormCron] fetch forms error: %v", err)
			return
		}

		companyData := dto.CompanyFormItems{
			CompanyID: company.ID,
			Forms:     forms,
		}
		data = append(data, companyData)
	}

	result, err := f.aiService.SendFormsToAI(ctx, data)
	if err != nil {
		log.Printf("[FormCron] send to AI error: %v", err)
		return
	}

	log.Printf("[FormCron] AI response: %s", result.Message)
	log.Println("[FormCron] Done.")
}
