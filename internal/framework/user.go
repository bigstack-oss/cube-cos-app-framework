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
		log.Errorf("openstack: failed to create user(%v)", err)
		return err
	}

	h.Spec.Openstack.User.ID = user.ID
	h.Spec.Openstack.User.Password = opts.Password
	return nil
}

func (h *Helper) deleteUsers() error {
	opts := users.ListOpts{DomainID: h.Spec.Openstack.Domain.ID}
	list, err := h.Openstack.ListUsers(&opts)
	if err != nil {
		log.Errorf("openstack: failed to list users(%v)", err)
		return err
	}

	for _, user := range list {
		if user.Name == h.Spec.Openstack.User.Name {
			err := h.Openstack.DeleteUser(user.ID)
			if err != nil {
				log.Errorf("openstack: failed to delete user %s(%v)", user.Name, err)
				continue
			}

			log.Infof("openstack: user %s is deleted successfully", user.Name)
		}
	}

	return nil
}

func (h *Helper) genUserCreationOpts() users.CreateOpts {
	return users.CreateOpts{
		Name:             h.Spec.Openstack.User.Name,
		Password:         h.genUserPassword(h.Spec.Openstack.User.Name),
		DefaultProjectID: h.Spec.Openstack.Project.ID,
	}
}
