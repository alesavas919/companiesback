package models

import (
	"time"
)

type Company struct {
	ID         int8      `json:"-"`
	Ticker     string    `json:"ticker"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

type CompanyInsert struct {
	ID         int8      `json:"-"`
	Ticker     string    `json:"ticker"`
	TargetFrom float64   `json:"target_from"`
	TargetTo   float64   `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

type CompaniesReq struct {
	Items []Company `json:"items"`
}

type CompaniesReqInsert struct {
	Items []CompanyInsert `json:"items"`
}

type CompaniesInsertCalculated struct {
	Info                CompanyInsert
	RatingPoints        float64 `json:"rating_points"`
	ActionPoints        float64 `json:"action_points"`
	TargetPoints        float64 `json:"target_points"`
	TotalPoints         float64 `json:"total_points"`
	TargetToBetThanFrom bool    `json:"target_to_bet_than_from"`
	Summary             string  `json:"summary"`
}

func (CompaniesReq) TableName() string {
	return "companiesreq"
}
