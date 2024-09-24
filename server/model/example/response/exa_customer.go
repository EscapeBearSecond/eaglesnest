package response

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
