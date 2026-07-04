# Database Schema & Model

Schema and Model that have :
1. Team
2. TeamInfo



## Team Schema
1. Legacy compatibility
    Because Team Schema is exist in legacy system before. we must aware about the migration. in new migration this project, create team **If Only** that table is not exist.
2. This is legacy golang struct that reflected the schema. for now, the field and schema already accomodate this system. No need to change.
    ```
    type Team struct {
        ID                uint     `json:"id" gorm:"primarykey"`
        Type              TeamType `json:"type"`
        Name              string   `json:"name"`
        TeamCode          TeamCode `json:"team_code" gorm:"index:team_code_unique,unique" binding:"required"`
        Description       string   `json:"desc"`

        Deleted  bool         `json:"deleted"`
    }
    ```
3. if on `./team_models` doesn't have golang model for that definition, duplicate legacy and place at `./team_models/team.go`



## TeamInfo Schema
1. Legacy compatibility
    Because TeamInfo Schema is exist in legacy system before. we must aware about the migration. in new migration this project, create TeamInfo **If Only** that table is not exist.
2. This is legacy golang struct that reflected the schema. for now, the field and schema already accomodate this system. No need to change.
    ```
    type TeamInfo struct {
        ID     uint `json:"id" gorm:"primaryKey"`
        TeamID uint `json:"team_id"`

        ReturnWarehouseID *uint `json:"return_warehouse_id"` // untuk destinasi gudang return team
        ReturnUserID      *uint `json:"return_user_id"`      // untuk destinasi gudang return team

        // untuk info agar team lain bisa transfer
        ContactNumber     string `json:"contact_number"`
        BankType          string `json:"bank_type"`
        BankOwnerName     string `json:"bank_owner_name"`
        BankAccountNumber string `json:"bank_account_number"`
    }
    ```
3. if on `./team_models` doesn't have golang model for definition, duplicate legacy and place at `./team_models/team.go`