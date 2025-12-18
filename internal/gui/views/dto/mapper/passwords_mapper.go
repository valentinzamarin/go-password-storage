package mapper

import (
	"password-storage/internal/app/common"
	"password-storage/internal/app/query"
	"password-storage/internal/gui/views/dto/response"
)

func ToPasswordsResponse(password *common.PasswordsResult) *response.PasswordsResponse {
	return &response.PasswordsResponse{
		ID:          password.ID,
		URL:         password.URL,
		Login:       password.Login,
		Password:    password.Password,
		Description: password.Description,
	}
}

func ToPasswordsResponseList(queryResult *query.PasswordsQueryResult) []*response.PasswordsResponse {
	responses := make([]*response.PasswordsResponse, 0, len(queryResult.Result))

	for _, password := range queryResult.Result {
		responses = append(responses, ToPasswordsResponse(password))
	}

	return responses
}
