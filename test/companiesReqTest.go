package test

import "companies/service"

func CompaniesReqGetAllDataFromDatabaseExperimental() []byte {
	companies := service.CompaniesReqGetAllDataFromRequestService()
	return companies
}

// ///////////////////////////////////////// GET FROM PETITION ///////////////////////////////////////////
func CompaniesReqGetAllDataFromRequestExperimental() string {
	resData := service.CompaniesReqGetAllDataFromRequestService()
	return string([]byte(resData))
}

/*
func CompaniesReqLoadAllDataFromRequestExperimental() ([]byte, error) {
	resData, err := service.CompaniesReqLoadAllDataFromRequestService()
	return resData, err //{created:'created'}
}
*/
func FunnyCalc(a int, b int) int {
	resultado := a + b
	return resultado
}
