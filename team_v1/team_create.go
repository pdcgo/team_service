package team_v1

import (
	"context"
	"io"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	role_base "github.com/pdcgo/schema/services/role_base/v1"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
	"github.com/pdcgo/user_service/access_interceptors"
	"github.com/pdcgo/user_service/user_models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// teamCreateLogger streams every slog line to the client as a TeamCreateResponse.message,
// per the long-running-task RPC pattern (see docs/code-implementation-guideline.md).
type teamCreateLogger struct {
	stream *connect.ServerStream[team_iface.TeamCreateResponse]
}

// Write implements [io.Writer].
func (l *teamCreateLogger) Write(p []byte) (int, error) {
	err := l.stream.Send(&team_iface.TeamCreateResponse{Message: string(p)})
	return len(p), err
}

// TeamCreate implements [team_ifaceconnect.TeamServiceHandler]. It creates the team, its
// team info, and adds the caller as team owner — all in one transaction — streaming a
// progress line per step and a final message carrying the created team. Admin only.
func (s *teamServiceImpl) TeamCreate(
	ctx context.Context,
	req *connect.Request[team_iface.TeamCreateRequest],
	stream *connect.ServerStream[team_iface.TeamCreateResponse],
) error {
	caller, err := access_interceptors.GetIdentityFromCtx(ctx)
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	var logwriter io.Writer = &teamCreateLogger{stream: stream}
	logger := slog.New(slog.NewTextHandler(logwriter, nil))

	var created *team_models.Team
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var terr error
		created, terr = createTeam(tx, req.Msg, uint(caller.IdentityId), func(msg string) {
			logger.Info(msg)
		})
		return terr
	})
	if err != nil {
		logger.Error("team create failed", "err", err.Error())
		return err
	}

	return stream.Send(&team_iface.TeamCreateResponse{
		Team:    toProtoTeam(created),
		Message: "team create done",
	})
}

// createTeam runs the three create steps within the given transaction, logging a progress
// line after each. It is the testable core of TeamCreate (no stream / ctx identity needed).
func createTeam(
	tx *gorm.DB,
	in *team_iface.TeamCreateRequest,
	callerID uint,
	log func(string),
) (*team_models.Team, error) {
	// 1. create the team (unique team_code enforced by the index).
	team := &team_models.Team{
		Type:        teamTypeToModel(in.Type),
		Name:        in.Name,
		TeamCode:    team_models.TeamCode(strings.ToUpper(in.TeamCode)),
		Description: in.Description,
	}
	err := tx.Create(team).Error
	if err != nil {
		return nil, err
	}
	log("team created")

	// 2. create the team info.
	info := &team_models.TeamInfo{TeamID: team.ID}
	err = tx.Create(info).Error
	if err != nil {
		return nil, err
	}
	team.TeamInfo = info
	log("team info created")

	// 3. add the creator as team owner (idempotent upsert on (team_id, user_id)).
	owner := &user_models.UserTeamRole{
		TeamID: team.ID,
		UserID: callerID,
		Role:   role_base.Role_ROLE_TEAM_OWNER,
		Alias:  "own",
	}
	err = tx.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "team_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"role", "alias"}),
		}).
		Create(owner).
		Error
	if err != nil {
		return nil, err
	}
	log("team owner added")

	return team, nil
}
