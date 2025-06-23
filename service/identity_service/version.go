package impl

import (
	"context"
	"runtime"

	"github.com/chise0904/golang_template/service"
)

var (
	sha1ver   string
	buildTime string
	version   string
)

// Version implements service.IdentityService.
func (s *svc) Version(ctx context.Context) (result *service.Version) {

	result = &service.Version{
		Version:      version,
		BuildTime:    buildTime,
		Sha1ver:      sha1ver,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
	return
}
