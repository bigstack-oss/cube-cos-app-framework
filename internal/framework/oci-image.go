package framework

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	log "go-micro.dev/v5/logger"
)

type ImageResp struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Errors []Error  `json:"errors"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *Helper) CheckOciImages() error {
	for _, image := range h.Spec.Framework.OciImages {
		u := url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", base.DataCenterVip, 5080),
			Path:   fmt.Sprintf("/v2/%s/%s/tags/list", image.Space, image.Name),
		}

		resp, err := h.http.R().SetResult(&ImageResp{}).Get(u.String())
		if err != nil {
			log.Errorf("framework: failed to send oci images request(%v)", err)
			return err
		}

		if resp.IsError() {
			err := fmt.Errorf("framework: has response error: %d(%s)", resp.StatusCode(), resp.String())
			return err
		}

		info := resp.Result().(*ImageResp)
		if !slices.Contains(info.Tags, image.Tag) {
			err := fmt.Errorf("framework: oci image %s/%s:%s not found", image.Space, image.Name, image.Tag)
			log.Errorf("framework: %v", err)
			return err
		}
	}

	return nil
}
