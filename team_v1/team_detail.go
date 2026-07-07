package team_v1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
	"gorm.io/gorm"
)

// TeamDetail implements [team_ifaceconnect.TeamServiceHandler]. It returns one team with
// its team info. Any authenticated caller.
func (s *teamServiceImpl) TeamDetail(
	ctx context.Context,
	req *connect.Request[team_iface.TeamDetailRequest],
) (*connect.Response[team_iface.TeamDetailResponse], error) {
	pay := req.Msg
	db := s.db.WithContext(ctx)

	var team team_models.Team
	err := db.
		Preload("TeamInfo").
		Where("id = ?", pay.TeamId).
		First(&team).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("team not found"))
		}
		return nil, err
	}

	return connect.NewResponse(&team_iface.TeamDetailResponse{Team: toProtoTeam(&team)}), nil
}
