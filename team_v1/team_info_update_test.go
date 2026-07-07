package team_v1

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"
	"github.com/pdcgo/team_service/team_models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTeamInfoUpdate(t *testing.T) {
	var scenario moretest_mock.DbScenario
	moretest.Suite(t, "team info update",
		moretest.SetupListFunc{moretest_mock.MockPostgresDatabase(&scenario)},
		func(t *testing.T) {
			scenario(t, func(tx *gorm.DB) {
				assert.NoError(t, tx.AutoMigrate(
					&team_models.Team{},
					&team_models.TeamInfo{},
				))

				team := &team_models.Team{Name: "Alpha", TeamCode: "ALPHA", Type: team_models.SellingTeamType}
				assert.NoError(t, tx.Create(team).Error)
				assert.NoError(t, tx.Create(&team_models.TeamInfo{TeamID: team.ID}).Error) // seeded by TeamCreate

				svc := NewTeamService(tx)
				ctx := context.Background()

				t.Run("sets bank / return / contact", func(t *testing.T) {
					res, err := svc.TeamInfoUpdate(ctx, connect.NewRequest(&team_iface.TeamInfoUpdateRequest{
						TeamId:            uint64(team.ID),
						ContactNumber:     "0812",
						BankType:          "BCA",
						BankOwnerName:     "Owner",
						BankAccountNumber: "123456",
						ReturnWarehouseId: 5,
						ReturnUserId:      9,
					}))
					assert.NoError(t, err)
					assert.Equal(t, "BCA", res.Msg.Info.BankType)
					assert.Equal(t, uint64(5), res.Msg.Info.ReturnWarehouseId)
					assert.Equal(t, uint64(9), res.Msg.Info.ReturnUserId)

					var info team_models.TeamInfo
					assert.NoError(t, tx.Where("team_id = ?", team.ID).First(&info).Error)
					assert.Equal(t, "0812", info.ContactNumber)
					assert.Equal(t, "123456", info.BankAccountNumber)
					assert.NotNil(t, info.ReturnWarehouseID)
					assert.Equal(t, uint(5), *info.ReturnWarehouseID)
				})

				t.Run("full replace; return ids 0 clears to NULL; no duplicate row", func(t *testing.T) {
					_, err := svc.TeamInfoUpdate(ctx, connect.NewRequest(&team_iface.TeamInfoUpdateRequest{
						TeamId:            uint64(team.ID),
						BankType:          "BNI",
						ReturnWarehouseId: 0,
						ReturnUserId:      0,
					}))
					assert.NoError(t, err)

					var infos []team_models.TeamInfo
					assert.NoError(t, tx.Where("team_id = ?", team.ID).Find(&infos).Error)
					assert.Len(t, infos, 1) // updated in place, not duplicated
					assert.Equal(t, "BNI", infos[0].BankType)
					assert.Empty(t, infos[0].ContactNumber) // omitted → cleared (full replace)
					assert.Nil(t, infos[0].ReturnWarehouseID)
					assert.Nil(t, infos[0].ReturnUserID)
				})

				t.Run("missing team → not found", func(t *testing.T) {
					_, err := svc.TeamInfoUpdate(ctx, connect.NewRequest(&team_iface.TeamInfoUpdateRequest{
						TeamId: 999999, BankType: "X",
					}))
					assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
				})
			})
		},
	)
}
