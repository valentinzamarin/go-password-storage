package request

import "password-storage/internal/app/query"

type IsMasterPasswordSetRequest struct{}

func (req *IsMasterPasswordSetRequest) ToQuery() *query.IsMasterPasswordSetQuery {
	return &query.IsMasterPasswordSetQuery{}
}
