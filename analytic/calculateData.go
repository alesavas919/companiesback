package analytic

import (
	"companies/database"
	"companies/models"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AnalyticCaculatedResponse(c *gin.Context) {
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
	c.JSON(http.StatusAccepted, AnalyticCalculated(companies))
}

func AnalyticCalculated(companyInsertListReq []models.CompanyInsert) []models.CompaniesInsertCalculated {
	companyInsertList := []models.CompaniesInsertCalculated{}
	var ratingFromPoints float64 = 0
	var raitingToPoints float64 = 0
	for i := 0; i < len(companyInsertListReq); i += 1 {
		companyInsertList = append(companyInsertList, models.CompaniesInsertCalculated{
			Info:                companyInsertListReq[i],
			RatingPoints:        0,
			ActionPoints:        0,
			TargetPoints:        0,
			TotalPoints:         0,
			TargetToBetThanFrom: false,
			Summary:             "",
		})
		ratingFromPoints = 0
		raitingToPoints = 0
		/*
			- values 1 -
			1			2			3			4
			Strong Buy, Buy, Outperform, Overweight
			5		6			7			8
			Hold, Neutral, Equal Weight, Market Perform
			x		x			x			x
			Sell, Reduce, Underperform, Underweight
			from -> to | based on
			n -> buy | nice
			buy -> buy | good
			n -> n | depends (x)
			buy -> n | no recommend
		*/
		ratingFromPoints = DataRatingPointsCalc(companyInsertList[i].Info.RatingFrom)
		raitingToPoints = DataRatingPointsCalc(companyInsertList[i].Info.RatingTo)
		companyInsertList[i].RatingPoints = ratingFromPoints + raitingToPoints
		/*
			- values 2 -

			"initiated by", "resumed by"

			Estas acciones reflejan un cambio en la calificación de la acción:
			"upgraded by" ,"reiterated by" ,"downgraded by"

			Reflejan cambios en la estimación de valor de la acción:
			"target raised by",  "price target maintained by", "target lowered by" (X)
		*/
		companyInsertList[i].ActionPoints = DataActionPointsCalc(companyInsertList[i].Info.Action,
			raitingToPoints)
		/*
			- value 3 -
			((target_to - target_from) / target_from) * 100

			target_to > target_from -> Positive value

			target_to < target_from -> Negative value
		*/
		companyInsertList[i].TargetPoints = DataTargetPointsCalc(companyInsertList[i].Info.TargetTo,
			companyInsertList[i].Info.TargetFrom)

		companyInsertList[i].TotalPoints = companyInsertList[i].ActionPoints +
			companyInsertList[i].RatingPoints + companyInsertList[i].TargetPoints

		companyInsertList[i].TargetToBetThanFrom = companyInsertList[i].Info.TargetTo >
			companyInsertList[i].Info.TargetFrom

		companyInsertList[i].Summary = SummaryGenerator(companyInsertList[i])
	}
	return companyInsertList
}

func DataRatingPointsCalc(dataTo string) float64 {
	dataTo = strings.ToUpper(dataTo)
	var value float64 = 0
	//---------------------------------------
	if dataTo == strings.ToUpper("Hold") {
		value = 6
	}
	if dataTo == strings.ToUpper("Neutral") {
		value = 5
	}
	if dataTo == strings.ToUpper("Equal Weight") {
		value = 4
	}
	if dataTo == strings.ToUpper("Market Perform") {
		value = 3
	}
	//---------------------------
	if dataTo == strings.ToUpper("Overweight") {
		value = 7
	}
	if dataTo == strings.ToUpper("Outperform") {
		value = 8
	}
	if dataTo == strings.ToUpper("Buy") {
		value = 9
	}
	if dataTo == strings.ToUpper("Strong Buy") {
		value = 10
	}
	return value
}

func DataActionPointsCalc(action string, raitingToPoints float64) float64 {
	action = strings.ToUpper(action)
	var value float64 = 0
	if action == strings.ToUpper("initiated by") || action == strings.ToUpper("resumed by") {
		raitingToPoints -= 6
		if raitingToPoints <= 0 {
			raitingToPoints -= 1
			raitingToPoints = (raitingToPoints / 5 * 10) - 2
		} else {
			raitingToPoints = (raitingToPoints / 5 * 10) + 2
		}
		value = raitingToPoints
	}

	if action == strings.ToUpper("upgraded by") {
		//if raitingToPoints >= 4
		value += (((raitingToPoints - 3) / 7) * 100) + 1
	}
	if action == strings.ToUpper("reiterated by") {
		value = 0
	}
	if action == strings.ToUpper("downgraded by") {
		value = -1
		raitingToPoints -= 6
		if raitingToPoints <= 0 {
			raitingToPoints -= 1
			raitingToPoints = (raitingToPoints / 5 * 10) - 2
		} else {
			raitingToPoints = (raitingToPoints / 5 * 10) + 2
		}
		value += raitingToPoints
	}

	if action == strings.ToUpper("target raised by") {
		value = 1
	}
	if action == strings.ToUpper("price target maintained by") {
		value = 0
	}
	if action == strings.ToUpper("target lowered by") {
		value = -1
	}
	return value
}

func DataTargetPointsCalc(targetTo float64, targetFrom float64) float64 {
	var value float64 = ((targetTo - targetFrom) / targetFrom) * 100
	return value
}

// //////////////////////////////
func SummaryGenerator(companyInsert models.CompaniesInsertCalculated) string {
	var summary = ``
	var generalInfo = `De la entidad ` + companyInsert.Info.Company + ` se puede decir lo siguiente: 
	`
	var inversion = `Ya que la inversion inicial ` + strconv.FormatFloat(companyInsert.Info.TargetFrom, 'f', 2, 64) +
		` y la inversion final ` + strconv.FormatFloat(companyInsert.Info.TargetTo, 'f', 2, 64) +
		` significando que existe ` + TargetGenerator(companyInsert.Info.TargetTo, companyInsert.Info.TargetFrom) + `
		`
	var action = `Las acciones se encuentran en modo ` + companyInsert.Info.Action + ` lo cual sería ` +
		ActionGenerator(companyInsert.Info.Action) + `
		`
	var rating = `Las compras de acciones en su estado inicial es de ` + companyInsert.Info.RatingFrom +
		` y la inversion a futuro es de ` + companyInsert.Info.RatingTo + " " +
		RatingGenerator(companyInsert.Info.RatingTo, companyInsert.Info.RatingFrom)
	summary = generalInfo + inversion + action + rating
	return summary
}

func TargetGenerator(targetTo float64, targetFrom float64) string {
	var response = "Un comportamiento "
	if targetTo > targetFrom {
		response += "POSITIVO [+] "
	} else {
		if targetTo < targetFrom {
			response += "NEGATIVO [-] "
		} else {
			response += "NEUTRAL [=] "
		}
	}
	return response
}

func ActionGenerator(action string) string {
	action = strings.ToUpper(action)
	var response = "Un comportamiento "
	if action == strings.ToUpper("initiated by") || action == strings.ToUpper("resumed by") {
		response += "INICIAL, no se presenta NINGUNA ANOMALIA "
	}

	if action == strings.ToUpper("upgraded by") {
		response += "POSITIVO [+], ya que calificación de la accion a futuro MEJORÓ "
	}

	if action == strings.ToUpper("reiterated by") {
		response += "NEUTRAL [=], ya que calificación de la accion a futuro NO PRESENTA ANOMALIAS "
	}

	if action == strings.ToUpper("downgraded by") {
		response += "NEGATIVO [-], ya que calificación de la accion a futuro EMPEORÓ "
	}

	if action == strings.ToUpper("target raised by") {
		response += "POSITIVO [+], ya que AUMENTÓ el precio a objetivo "
	}

	if action == strings.ToUpper("price target maintained by") {
		response += "NEUTRAL [=], ya que SE MANTIENE el precio a objetivo "
	}

	if action == strings.ToUpper("target lowered by") {
		response += "NEGATIVO [-], ya que DISMINUYÓ el precio a objetivo "
	}

	return response
}

func RatingGenerator(ratingTo string, ratingFrom string) string {
	var response = "Un comportamiento "
	var dataToValue = DataRatingPointsCalc(ratingTo)
	var dataFromValue = DataRatingPointsCalc(ratingFrom)
	var calcDataToFromValue = (dataToValue - dataFromValue)
	if dataToValue < dataFromValue {
		if dataToValue <= 6 {
			calcDataToFromValue = -1
		}
	}
	//TODO UPGRADE THE METHOD rating & target
	if dataToValue > dataFromValue {
		response += "POSITIVO [+], " // + strconv.FormatFloat(calcDataToFromValue, 'f', 2, 64) + " "
		if calcDataToFromValue >= 0 && calcDataToFromValue <= 2 {
			response += "Se recomienda invertir a CORTO, MEDIANO y LARGO PLAZO "
		}
		if calcDataToFromValue >= 3 && calcDataToFromValue <= 5 {
			response += "Se recomienda invertir a MEDIANO y LARGO PLAZO "
		}
		if calcDataToFromValue >= 6 && calcDataToFromValue <= 10 {
			response += "Se recomienda invertir a LARGO PLAZO "
		}
	} else {
		if dataToValue < dataFromValue {
			response += "NEGATIVO [-], la compra tiende a haber perdidas futuras "
		} else {
			response = "Un comportamiento " + "NEUTRA [=], no es recomendable comprar ya que no se tiene nunguna ganancia "
		}
	}
	return response
}
