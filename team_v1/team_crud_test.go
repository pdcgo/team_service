package team_v1

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	common "github.com/pdcgo/schema/services/common/v1"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"
	"github.com/pdcgo/team_service/team_models"
	"github.com/pdcgo/user_service/user_models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTeamCrud(t *testing.T) {
	var scenario moretest_mock.DbScenario
	moretest.Suite(t, "team crud",
		moretest.SetupListFunc{moretest_mock.MockPostgresDatabase(&scenario)},
		func(t *testing.T) {
			scenario(t, func(tx *gorm.DB) {
				assert.NoError(t, tx.AutoMigrate(
					&team_models.Team{},
					&team_models.TeamInfo{},
					&user_models.UserTeamRole{},
				))

				seed := func(name, code string, tp team_models.TeamType, deleted bool) *team_models.Team {
					team := &team_models.Team{Name: name, TeamCode: team_models.TeamCode(code), Type: tp, Deleted: deleted}
					assert.NoError(t, tx.Create(team).Error)
					assert.NoError(t, tx.Create(&team_models.TeamInfo{TeamID: team.ID, ContactNumber: "0812"}).Error)
					return team
				}
				sel := seed("Alpha Selling", "ALPHA", team_models.SellingTeamType, false)
				beta := seed("Beta Warehouse", "BETA", team_models.WarehouseTeamType, false)
				gone := seed("Gone Team", "GONE", team_models.SellingTeamType, true) // soft-deleted

				svc := NewTeamService(tx)
				ctx := context.Background()
				page := func() *common.PageFilter { return &common.PageFilter{Page: 1, Limit: 20} }

				t.Run("list excludes deleted, filters by type and keyword", func(t *testing.T) {
					res, err := svc.TeamList(ctx, connect.NewRequest(&team_iface.TeamListRequest{Page: page()}))
					assert.NoError(t, err)
					assert.Len(t, res.Msg.Teams, 2) // GONE excluded

					res, err = svc.TeamList(ctx, connect.NewRequest(&team_iface.TeamListRequest{
						TeamType: common.TeamType_TEAM_TYPE_WAREHOUSE,
						Page:     page(),
					}))
					assert.NoError(t, err)
					assert.Len(t, res.Msg.Teams, 1)
					assert.Equal(t, "Beta Warehouse", res.Msg.Teams[0].Name)

					res, err = svc.TeamList(ctx, connect.NewRequest(&team_iface.TeamListRequest{Q: "alpha", Page: page()}))
					assert.NoError(t, err)
					assert.Len(t, res.Msg.Teams, 1)
					assert.Equal(t, "ALPHA", res.Msg.Teams[0].TeamCode)
				})

				t.Run("detail returns the team with info", func(t *testing.T) {
					res, err := svc.TeamDetail(ctx, connect.NewRequest(&team_iface.TeamDetailRequest{TeamId: uint64(sel.ID)}))
					assert.NoError(t, err)
					assert.Equal(t, "Alpha Selling", res.Msg.Team.Name)
					assert.NotNil(t, res.Msg.Team.Info)
					assert.Equal(t, "0812", res.Msg.Team.Info.ContactNumber)
				})

				t.Run("detail missing team → not found", func(t *testing.T) {
					_, err := svc.TeamDetail(ctx, connect.NewRequest(&team_iface.TeamDetailRequest{TeamId: 999999}))
					assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
				})

				t.Run("by ids returns non-deleted keyed by id, omits missing and deleted", func(t *testing.T) {
					res, err := svc.TeamByIds(ctx, connect.NewRequest(&team_iface.TeamByIdsRequest{
						Ids: []uint64{uint64(sel.ID), uint64(beta.ID), uint64(gone.ID), 999999},
					}))
					assert.NoError(t, err)
					assert.Len(t, res.Msg.Data, 2) // gone (soft-deleted) and 999999 (missing) omitted
					assert.Equal(t, "Alpha Selling", res.Msg.Data[uint64(sel.ID)].Name)
					assert.Equal(t, "Beta Warehouse", res.Msg.Data[uint64(beta.ID)].Name)
					assert.Nil(t, res.Msg.Data[uint64(gone.ID)])
					assert.Nil(t, res.Msg.Data[999999])
				})

				t.Run("update changes name/description only", func(t *testing.T) {
					res, err := svc.TeamUpdate(ctx, connect.NewRequest(&team_iface.TeamUpdateRequest{
						TeamId:      uint64(sel.ID),
						Name:        "Alpha Renamed",
						Description: "updated",
					}))
					assert.NoError(t, err)
					assert.Equal(t, "Alpha Renamed", res.Msg.Team.Name)
					assert.Equal(t, "updated", res.Msg.Team.Description)
					assert.Equal(t, "ALPHA", res.Msg.Team.TeamCode)                         // immutable
					assert.Equal(t, string(team_models.SellingTeamType), res.Msg.Team.Type) // immutable
				})

				t.Run("update missing team → not found", func(t *testing.T) {
					_, err := svc.TeamUpdate(ctx, connect.NewRequest(&team_iface.TeamUpdateRequest{
						TeamId: 999999, Name: "no such", Description: "",
					}))
					assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
				})

				// Soft-delete last: it drops sel from the list.
				t.Run("delete soft-deletes and drops from list", func(t *testing.T) {
					_, err := svc.TeamDelete(ctx, connect.NewRequest(&team_iface.TeamDeleteRequest{TeamId: uint64(sel.ID)}))
					assert.NoError(t, err)

					var team team_models.Team
					assert.NoError(t, tx.Where("id = ?", sel.ID).First(&team).Error)
					assert.True(t, team.Deleted)

					res, err := svc.TeamList(ctx, connect.NewRequest(&team_iface.TeamListRequest{Page: page()}))
					assert.NoError(t, err)
					assert.Len(t, res.Msg.Teams, 1) // only Beta remains

					// deleting an already-deleted team → not found.
					_, err = svc.TeamDelete(ctx, connect.NewRequest(&team_iface.TeamDeleteRequest{TeamId: uint64(sel.ID)}))
					assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
				})
			})
		},
	)
}
