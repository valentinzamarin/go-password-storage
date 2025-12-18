package query

import "password-storage/internal/app/common"

// Structure PasswordsResult with fields ID, URL, Login, Password, Description
// not a query, a result.
// query should only contain search criteria, not entity business data.

type PasswordsQueryResult struct {
	Result []*common.PasswordsResult
}

// type PasswordsQuery struct {
// 	ID          uint
// 	URL         string
// 	Login       string
// 	Password    string
// 	Description string
// }
