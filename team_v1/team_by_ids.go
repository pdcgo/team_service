package team_v1

import (
	"context"

	"connectrpc.com/connect"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
)

// TeamByIds implements [team_ifaceconnect.TeamServiceHandler]. It bulk-loads non-deleted
// teams by id, keyed by id, for preloading team names in other UIs (so raw ids are never
// shown). Missing and soft-deleted ids are omitted from the map. Any authenticated caller.
func (s *teamServiceImpl) TeamByIds(
	ctx context.Context,
	req *connect.Request[team_iface.TeamByIdsRequest],
) (*connect.Response[team_iface.TeamByIdsResponse], error) {
	pay := req.Msg
	result := &team_iface.TeamByIdsResponse{Data: map[uint64]*team_iface.Team{}}
	if len(pay.Ids) == 0 {
		return connect.NewResponse(result), nil
	}

	db := s.db.WithContext(ctx)

	var rows []*team_models.Team
	err := db.
		Where("deleted = ?", false).
		Where("id IN ?", pay.Ids).
		Find(&rows).
		Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result.Data[uint64(row.ID)] = toProtoTeam(row)
	}

	return connect.NewResponse(result), nil
}
