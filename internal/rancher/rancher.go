package rancher

import (
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	log "go-micro.dev/v5/logger"
)

type Helper struct {
	Http    *http.Helper
	Options Options
}

func NewHelper(opts ...Option) *Helper {
	syncedOpts := syncOptions(opts)
	c, err := http.NewHelper(
		http.TlsInsecureSkipVerify(true),
	)
	if err != nil {
		log.Errorf("failed to init http helper: %s \n", err.Error())
		return nil
	}

	return &Helper{
		Http:    c,
		Options: *syncedOpts,
	}
}

func syncOptions(opts []Option) *Options {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	return options
}

func genAuthHeaders(token string) map[string]string {
	return map[string]string{
		"content-type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
}
