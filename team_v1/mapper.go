package team_v1

import (
	common "github.com/pdcgo/schema/services/common/v1"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
)

// toProtoTeam maps a team model to its proto message (info included when preloaded).
func toProtoTeam(m *team_models.Team) *team_iface.Team {
	if m == nil {
		return nil
	}
	return &team_iface.Team{
		Id:          uint64(m.ID),
		Type:        string(m.Type),
		Name:        m.Name,
		TeamCode:    string(m.TeamCode),
		Description: m.Description,
		Deleted:     m.Deleted,
		Info:        toProtoTeamInfo(m.TeamInfo),
	}
}

func toProtoTeamInfo(m *team_models.TeamInfo) *team_iface.TeamInfo {
	if m == nil {
		return nil
	}
	info := &team_iface.TeamInfo{
		TeamId:            uint64(m.TeamID),
		ContactNumber:     m.ContactNumber,
		BankType:          m.BankType,
		BankOwnerName:     m.BankOwnerName,
		BankAccountNumber: m.BankAccountNumber,
	}
	if m.ReturnWarehouseID != nil {
		info.ReturnWarehouseId = uint64(*m.ReturnWarehouseID)
	}
	if m.ReturnUserID != nil {
		info.ReturnUserId = uint64(*m.ReturnUserID)
	}
	return info
}

// teamTypeToModel maps the proto TeamType enum to the stored TeamType string
// ("" for unspecified). The proto enum has no root variant (root teams aren't
// creatable via the API).
func teamTypeToModel(t common.TeamType) team_models.TeamType {
	switch t {
	case common.TeamType_TEAM_TYPE_WAREHOUSE:
		return team_models.WarehouseTeamType
	case common.TeamType_TEAM_TYPE_SELLING:
		return team_models.SellingTeamType
	case common.TeamType_TEAM_TYPE_ADMIN:
		return team_models.AdminTeamType
	default:
		return ""
	}
}
