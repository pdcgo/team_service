package team_v1

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	common "github.com/pdcgo/schema/services/common/v1"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/team_service/team_models"
	"gorm.io/gorm"
)

// TeamList implements [team_ifaceconnect.TeamServiceHandler]. It returns non-deleted teams
// (newest first), optionally filtered by a keyword (name / team_code) and team type, paged.
// Any authenticated caller.
func (s *teamServiceImpl) TeamList(
	ctx context.Context,
	req *connect.Request[team_iface.TeamListRequest],
) (*connect.Response[team_iface.TeamListResponse], error) {
	pay := req.Msg
	if pay.Page == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("page is required"))
	}
	db := s.db.WithContext(ctx)

	var rows []*team_models.Team
	paginated, pageInfo, err := db_connect.SetPaginationQuery(db, func() (*gorm.DB, error) {
		query := db.
			Model(&team_models.Team{}).
			Scopes(func(d *gorm.DB) *gorm.DB {
				d = d.Where("deleted = ?", false)
				if q := strings.TrimSpace(pay.Q); q != "" {
					like := "%" + q + "%"
					d = d.Where("name ILIKE ? OR team_code ILIKE ?", like, like)
				}
				if pay.TeamType != common.TeamType_TEAM_TYPE_UNSPECIFIED {
					d = d.Where("type = ?", string(teamTypeToModel(pay.TeamType)))
				}
				return d
			})
		return query, nil
	}, pay.Page)
	if err != nil {
		return nil, err
	}

	err = paginated.Order("id DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := &team_iface.TeamListResponse{
		Teams:    make([]*team_iface.Team, 0, len(rows)),
		PageInfo: pageInfo,
	}
	for _, row := range rows {
		result.Teams = append(result.Teams, toProtoTeam(row))
	}
	return connect.NewResponse(result), nil
}
