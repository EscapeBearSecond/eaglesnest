package response

import "47.103.136.241/goprojects/curescan/server/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
