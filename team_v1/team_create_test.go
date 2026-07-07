package team_v1

import (
	"testing"

	common "github.com/pdcgo/schema/services/common/v1"
	role_base "github.com/pdcgo/schema/services/role_base/v1"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"
	"github.com/pdcgo/team_service/team_models"
	"github.com/pdcgo/user_service/user_models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTeam(t *testing.T) {
	var scenario moretest_mock.DbScenario
	moretest.Suite(t, "create team",
		moretest.SetupListFunc{moretest_mock.MockPostgresDatabase(&scenario)},
		func(t *testing.T) {
			scenario(t, func(tx *gorm.DB) {
				assert.NoError(t, tx.AutoMigrate(
					&team_models.Team{},
					&team_models.TeamInfo{},
					&user_models.UserTeamRole{},
				))

				noop := func(string) {}
				const callerID uint = 7

				t.Run("creates team, info, and owner", func(t *testing.T) {
					team, err := createTeam(tx, &team_iface.TeamCreateRequest{
						Type:        common.TeamType_TEAM_TYPE_SELLING,
						Name:        "Selling Alpha",
						TeamCode:    "alpha",
						Description: "first selling team",
					}, callerID, noop)
					assert.NoError(t, err)
					assert.NotZero(t, team.ID)
					assert.Equal(t, team_models.TeamCode("ALPHA"), team.TeamCode) // uppercased
					assert.Equal(t, team_models.SellingTeamType, team.Type)

					// team info row linked to the team.
					var info team_models.TeamInfo
					assert.NoError(t, tx.Where("team_id = ?", team.ID).First(&info).Error)

					// creator recorded as team owner.
					var role user_models.UserTeamRole
					assert.NoError(t, tx.Where("team_id = ? AND user_id = ?", team.ID, callerID).First(&role).Error)
					assert.Equal(t, role_base.Role_ROLE_TEAM_OWNER, role.Role)
					assert.Equal(t, "own", role.Alias)
				})

				// A duplicate team_code violates the unique index; keep this LAST — the
				// unique-violation aborts the shared scenario transaction.
				t.Run("duplicate team_code errors", func(t *testing.T) {
					_, err := createTeam(tx, &team_iface.TeamCreateRequest{
						Type:     common.TeamType_TEAM_TYPE_SELLING,
						Name:     "Selling Dup",
						TeamCode: "ALPHA",
					}, callerID, noop)
					assert.Error(t, err)
				})
			})
		},
	)
}
