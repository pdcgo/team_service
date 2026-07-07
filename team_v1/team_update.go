package team_v1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
)

// TeamUpdate implements [team_ifaceconnect.TeamServiceHandler]. It updates a team's name
// and description (type and team_code are immutable). Admin only.
func (s *teamServiceImpl) TeamUpdate(
	ctx context.Context,
	req *connect.Request[team_iface.TeamUpdateRequest],
) (*connect.Response[team_iface.TeamUpdateResponse], error) {
	pay := req.Msg
	db := s.db.WithContext(ctx)

	res := db.
		Model(&team_models.Team{}).
		Where("id = ? AND deleted = ?", pay.TeamId, false).
		Updates(map[string]interface{}{
			"name":        pay.Name,
			"description": pay.Description,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("team not found"))
	}

	var team team_models.Team
	err := db.
		Preload("TeamInfo").
		Where("id = ?", pay.TeamId).
		First(&team).
		Error
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&team_iface.TeamUpdateResponse{Team: toProtoTeam(&team)}), nil
}
