package role

//go:generate go tool enumer -type=Role -sql -json -transform=upper
type Role int

const (
	SUPERUSER Role = iota
	ADMIN
	USER
)
