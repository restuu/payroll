package dto

import "github.com/shopspring/decimal"

type RegisterNewEmployeeParams struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Salary   decimal.Decimal `json:"salary"`
	AdminID  int64           `json:"-"`
}
