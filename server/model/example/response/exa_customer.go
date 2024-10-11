package response

import "github.com/EscapeBearSecond/curescan/server/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
