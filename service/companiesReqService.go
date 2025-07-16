package service

import (
	"companies/commonData"
	"companies/database"
	"companies/models"
	"companies/security"

	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func CompaniesReqGetAllDataFromDatabaseService() []models.CompanyInsert {
	rows, err := database.DB.Query(context.Background(), `SELECT ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, "time" 
	FROM public.companiesreq order by "time" desc`)
	if err != nil {
		log.Fatal("Error response -", err)
	}
	defer rows.Close()

	var companies []models.CompanyInsert
	for rows.Next() {
		var company models.CompanyInsert
		rows.Scan(
			&company.Ticker,
			&company.TargetFrom,
			&company.TargetTo,
			&company.Company,
			&company.Action,
			&company.Brokerage,
			&company.RatingFrom,
			&company.RatingTo,
			&company.Time,
		)
		companies = append(companies, company)
	}
	return companies
}

func CompaniesReqGetAllDataFromRequestService() ([]byte, error) {
	//HTTP GET DATA FROM TR
	var bearer, err1 = security.SecretString(1) //security.ResourceSecurityData("REGIS_INFO_A")
	var url, err2 = security.SecretString(2)    //security.ResourceSecurityData("REGIS_LIST")
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	if bearer == "" || url == "" {
		return nil, nil
	}
	bearer = "Bearer " + bearer

	res, err := http.NewRequest("GET", url, nil)
	res.Header.Add("Authorization", bearer)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Fatal("Error response -", err)
	}
	defer resp.Body.Close()
	resData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Rading data Error: ", err)
	}
	return resData, nil
}

func CompaniesReqLoadAllDataFromRequestService() ([]byte, error) {
	//var companyModel models.CompaniesReqLst
	//HTTP GET DATA FROM TR
	var bearer, err1 = security.SecretString(1) //security.ResourceSecurityData("REGIS_INFO_A")
	var url, err2 = security.SecretString(2)    //security.ResourceSecurityData("REGIS_LIST")
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	if bearer == "" || url == "" {
		return nil, nil
	}

	res, err := http.NewRequest("GET", url, nil)
	res.Header.Add("Authorization", bearer)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Fatal("Error response -", err)
	}
	defer resp.Body.Close()
	resData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Reading data Error: ", err)
	}

	var companiesReq models.CompaniesReq
	err = json.Unmarshal([]byte(resData), &companiesReq)
	if err != nil {
		fmt.Println("JSON data rror:", err)
		return nil, err
	}
	for i := 0; i < len(companiesReq.Items); i++ {

		rows, err := getAllRegisters(companiesReq, i)
		if rows == 0 {
			TargetFrom, _ := strconv.ParseFloat(strings.Split(companiesReq.Items[i].TargetFrom, "$")[1], 64)
			TargetTo, _ := strconv.ParseFloat(strings.Split(companiesReq.Items[i].TargetTo, "$")[1], 64)
			companyInsert := models.CompanyInsert{
				Ticker:     companiesReq.Items[i].Ticker,
				TargetFrom: TargetFrom,
				TargetTo:   TargetTo,
				Company:    companiesReq.Items[i].Company,
				Action:     companiesReq.Items[i].Action,
				Brokerage:  companiesReq.Items[i].Brokerage,
				RatingFrom: companiesReq.Items[i].RatingFrom,
				RatingTo:   companiesReq.Items[i].RatingTo,
				Time:       companiesReq.Items[i].Time,
			}
			createData(companyInsert)
		}
		if err != nil {
			return nil, err
		}

	}
	return resData, nil
}

func getAllRegisters(companiesReq models.CompaniesReq, i int) (int, error) {
	company, companyTime := strings.TrimSpace(strings.ToUpper(string(companiesReq.Items[i].Company))), companiesReq.Items[i].Time.Format(time.RFC3339)
	args := []interface{}{company}
	args = append(args, companyTime)
	rows, _ := database.DB.Query(context.Background(), `SELECT * FROM public.companiesreq WHERE UPPER(TRIM(company)) = $1 and "time" = $2`, args...)
	rows.Close()
	var valReturned = 1
	if rows.CommandTag().RowsAffected() == 0 {
		valReturned = 0
	}
	return valReturned, nil
}

func createData(companiesReq models.CompanyInsert) string {
	var createStatus int8 = 1
	createStatus *= commonData.ValueExists(companiesReq.Ticker, "Ticker ("+companiesReq.Ticker+")")
	if createStatus == 1 {
		createStatus *= commonData.ValueNumberExits(companiesReq.TargetFrom, "TargetFrom")
	}
	if createStatus == 1 {
		createStatus *= commonData.ValueNumberExits(companiesReq.TargetTo, "TargetFrom")
	}
	if createStatus == 1 {
		createStatus *= commonData.ValueExists(companiesReq.Company, "Company")
	}
	if createStatus == 1 {
		createStatus *= commonData.ValueExists(companiesReq.Action, "Action")
	}
	if createStatus == 1 {
		createStatus *= commonData.ValueExists(companiesReq.Brokerage, "Brokerage")
	}
	if createStatus == 1 {
		createStatus *= commonData.ValueExists(companiesReq.RatingFrom, "RatingFrom")
	}
	if createStatus == 1 {
		createStatus *= commonData.ValueExists(companiesReq.RatingTo, "RatingTo")
	}
	var TargetFrom = companiesReq.TargetFrom
	var TargetTo = companiesReq.TargetTo
	if createStatus == 1 {
		_, err := database.DB.Exec(context.Background(),
			`INSERT INTO public.companiesreq
				(ticker, target_from, target_to, company, "action", brokerage, rating_from, rating_to, "time") VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			string(companiesReq.Ticker),
			TargetFrom,
			TargetTo,
			string(companiesReq.Company),
			string(companiesReq.Action),
			string(companiesReq.Brokerage),
			string(companiesReq.RatingFrom),
			string(companiesReq.RatingTo),
			companiesReq.Time.Format(time.RFC3339)) //companiesReq.Time.Format(time.RFC3339))
		if err != nil {
			return "creating error"
		}
		return "created data"
	}
	return "completed"
}
