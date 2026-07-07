package team_models

// TeamType is the stored team classification. It mirrors shared/db_models.TeamType
// (a plain string in the legacy schema).
type TeamType string

const (
	RootTeamType      TeamType = "root"
	AdminTeamType     TeamType = "admin"
	WarehouseTeamType TeamType = "warehouse"
	SellingTeamType   TeamType = "selling"
)

func (o TeamType) String() string { return string(o) }

// TeamCode is a caller-supplied, uppercased, unique short code (no generator exists;
// uniqueness is enforced by the team_code_unique index).
type TeamCode string

// Team mirrors the legacy shared/db_models.Team columns that team_service reads/writes.
// The teams table is legacy-owned; migrations create it only if absent.
type Team struct {
	ID          uint     `json:"id" gorm:"primarykey"`
	Type        TeamType `json:"type"`
	Name        string   `json:"name"`
	TeamCode    TeamCode `json:"team_code" gorm:"index:team_code_unique,unique"`
	Description string   `json:"desc"`
	Deleted     bool     `json:"deleted"`

	TeamInfo *TeamInfo `json:"team_info" gorm:"foreignKey:TeamID"`
}

// TeamInfo mirrors the legacy shared/db_models.TeamInfo transfer/return metadata.
type TeamInfo struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	TeamID uint `json:"team_id"`

	ReturnWarehouseID *uint `json:"return_warehouse_id"`
	ReturnUserID      *uint `json:"return_user_id"`

	ContactNumber     string `json:"contact_number"`
	BankType          string `json:"bank_type"`
	BankOwnerName     string `json:"bank_owner_name"`
	BankAccountNumber string `json:"bank_account_number"`
}
