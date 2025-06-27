package framework

import (
	"crypto/hmac"
	"crypto/sha1"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

func (h *Helper) genUserPassword(name string) string {
	secretKey := []byte("bigstackcoltd")
	message := []byte(name)

	hash := hmac.New(sha1.New, secretKey)
	hash.Write(message)

	return "123456"
}

func (h *Helper) createUserAndApplyToProject() error {
	user, err := h.Openstack.CreateUser(h.genUserCreationOpts())
	if err != nil {
		return err
	}

	h.User.ID = user.ID
	h.User.Password = h.genUserPassword(h.User.Name)
	err = h.addUserAndRolesToProject()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) genUserCreationOpts() users.CreateOpts {
	return users.CreateOpts{
		Name:             h.User.Name,
		Password:         h.genUserPassword(h.User.Name),
		DefaultProjectID: h.Project.ID,
	}
}
