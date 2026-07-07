package team_v1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
)

// TeamDelete implements [team_ifaceconnect.TeamServiceHandler]. It soft-deletes a team
// (sets deleted = true). Admin only.
func (s *teamServiceImpl) TeamDelete(
	ctx context.Context,
	req *connect.Request[team_iface.TeamDeleteRequest],
) (*connect.Response[team_iface.TeamDeleteResponse], error) {
	pay := req.Msg
	db := s.db.WithContext(ctx)

	res := db.
		Model(&team_models.Team{}).
		Where("id = ? AND deleted = ?", pay.TeamId, false).
		Update("deleted", true)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("team not found"))
	}

	return connect.NewResponse(&team_iface.TeamDeleteResponse{}), nil
}
