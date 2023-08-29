package helper

import "echo-blog/models"

func WrapResponse(code int, status string, response interface{}) *models.WebResponse {
	newResponse := new(models.WebResponse)
	newResponse.Code = code
	newResponse.Status = status
	newResponse.Data = response
	return newResponse
}
