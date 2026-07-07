package team_v1

import (
	"github.com/pdcgo/schema/services/team_iface/v1/team_ifaceconnect"
	"gorm.io/gorm"
)

type teamServiceImpl struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *teamServiceImpl {
	return &teamServiceImpl{db: db}
}

// Compile-time assertion that the skeleton satisfies the generated handler.
var _ team_ifaceconnect.TeamServiceHandler = (*teamServiceImpl)(nil)
