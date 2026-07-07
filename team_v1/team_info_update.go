package team_v1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	team_iface "github.com/pdcgo/schema/services/team_iface/v1"
	"github.com/pdcgo/team_service/team_models"
	"gorm.io/gorm"
)

// TeamInfoUpdate implements [team_ifaceconnect.TeamServiceHandler]. It edits the team's
// transfer/return metadata (bank details, return point, contact number) — a full replace
// of the editable info fields. Team-scoped: a team owner/admin may edit their own team
// (root/admin bypass), enforced by the interceptor via the use_scope team_id. A
// return_*_id of 0 clears the value (NULL).
func (s *teamServiceImpl) TeamInfoUpdate(
	ctx context.Context,
	req *connect.Request[team_iface.TeamInfoUpdateRequest],
) (*connect.Response[team_iface.TeamInfoUpdateResponse], error) {
	pay := req.Msg
	db := s.db.WithContext(ctx)

	// The team must exist and not be soft-deleted.
	var team team_models.Team
	err := db.
		Where("id = ? AND deleted = ?", pay.TeamId, false).
		First(&team).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("team not found"))
		}
		return nil, err
	}

	// Load the existing info row (TeamCreate seeds one) or start a fresh one; team_infos
	// has no unique index on team_id, so this load-then-Save is the safe upsert.
	var info team_models.TeamInfo
	err = db.
		Where("team_id = ?", pay.TeamId).
		First(&info).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		info = team_models.TeamInfo{TeamID: uint(pay.TeamId)}
	}

	info.ContactNumber = pay.ContactNumber
	info.BankType = pay.BankType
	info.BankOwnerName = pay.BankOwnerName
	info.BankAccountNumber = pay.BankAccountNumber
	info.ReturnWarehouseID = optionalID(pay.ReturnWarehouseId)
	info.ReturnUserID = optionalID(pay.ReturnUserId)

	err = db.Save(&info).Error
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&team_iface.TeamInfoUpdateResponse{Info: toProtoTeamInfo(&info)}), nil
}

// optionalID maps a proto id (0 = unset) to a nullable *uint column value.
func optionalID(id uint64) *uint {
	if id == 0 {
		return nil
	}
	v := uint(id)
	return &v
}
