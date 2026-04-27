package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AssessmentResultModel struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionID        uuid.UUID `gorm:"type:uuid;not null;index"`
	RegulationItemID uuid.UUID `gorm:"type:uuid;not null;index"`
	ComplianceStatus string    `gorm:"not null"`
}

func (AssessmentResultModel) TableName() string { return "assessment_results" }

type RegulationAssessmentModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RegulationID uuid.UUID `gorm:"type:uuid;not null;index"`
	SessionID    uuid.UUID `gorm:"type:uuid;not null;index"`
	AmountPass   int
	AmountFail   int
	AmountNA     int
}

func (RegulationAssessmentModel) TableName() string { return "regulation_assesments" }

func main() {
	dsn := "host=localhost user=postgres password=Tambunan140705 dbname=grc_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- Assessment Results ---")
	var results []AssessmentResultModel
	db.Find(&results)
	for _, r := range results {
		fmt.Printf("Session: %s, Item: %s, Status: %s\n", r.SessionID, r.RegulationItemID, r.ComplianceStatus)
	}

	fmt.Println("\n--- Regulation Assessments (Summary) ---")
	var summaries []RegulationAssessmentModel
	db.Find(&summaries)
	for _, s := range summaries {
		fmt.Printf("Session: %s, Reg: %s, Pass: %d, Fail: %d, NA: %d\n", s.SessionID, s.RegulationID, s.AmountPass, s.AmountFail, s.AmountNA)
	}
}
