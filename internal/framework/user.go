package framework

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) genUserPassword(name string) string {
	hash := sha1.New()
	hash.Write([]byte(name + base.SystemSeed))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func (h *Helper) createUser() error {
	opts := h.genUserCreationOpts()
	user, err := h.Openstack.CreateUser(opts)
	if err != nil {
		log.Errorf("framework: failed to create user(%v)", err)
		return err
	}

	h.Spec.Openstack.User.ID = user.ID
	h.Spec.Openstack.User.Password = opts.Password
	return nil
}

func (h *Helper) genUserCreationOpts() users.CreateOpts {
	return users.CreateOpts{
		Name:             h.Spec.Openstack.User.Name,
		Password:         h.genUserPassword(h.Spec.Openstack.User.Name),
		DefaultProjectID: h.Spec.Openstack.Project.ID,
	}
}
