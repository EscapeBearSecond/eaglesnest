package response

import "github.com/EscapeBearSecond/eaglesnest/server/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
