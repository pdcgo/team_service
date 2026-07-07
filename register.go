package team_service

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/pdcgo/san_collection/san_caches"
	"github.com/pdcgo/schema/services/team_iface/v1/team_ifaceconnect"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/team_service/team_v1"
	"github.com/pdcgo/user_service/access_interceptors"
	"gorm.io/gorm"
)

type ServiceReflectNames []string
type RegisterHandler func() ServiceReflectNames

// NewRegister mounts the v2 Connect TeamService onto mux and returns its gRPC-reflection
// service name. The access interceptor enforces each request's (role_base.v1.request_policy)
// and injects the caller identity into context.
func NewRegister(
	mux *http.ServeMux,
	db *gorm.DB,
	cfg *configs.AppConfig,
	defaultInterceptor custom_connect.DefaultInterceptor,
	cacheMgr san_caches.CacheManager,
) RegisterHandler {
	return func() ServiceReflectNames {
		grpcReflects := ServiceReflectNames{}

		roleOpt := connect.WithInterceptors(access_interceptors.NewAccessInterceptor(db, cfg.JwtSecret, cacheMgr))
		path, handler := team_ifaceconnect.NewTeamServiceHandler(
			team_v1.NewTeamService(db),
			defaultInterceptor,
			roleOpt,
		)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, team_ifaceconnect.TeamServiceName)

		return grpcReflects
	}
}
