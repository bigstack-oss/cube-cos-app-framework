package framework

import (
	storageQuota "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/quotasets"
	computeQuota "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	networkQuotas "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	log "go-micro.dev/v5/logger"
)

type memberAndRoleId struct {
	MemberId string
	RoleId   string
}

func (h *Helper) createProject() error {
	var err error
	h.Spec.Openstack.Project, err = h.Openstack.CreateProject(h.Spec.Openstack.Project.Name)
	if err != nil {
		log.Errorf("framework: failed to create project(%v)", err)
		return err
	}

	log.Infof(
		"framework: project is applied successfully(%s %s)",
		h.Spec.Openstack.Project.Name,
		h.Spec.Openstack.Project.ID,
	)

	return nil
}

func (h *Helper) unlimitProjectQuotas() error {
	err := h.unlimitProjectResourceQuotas()
	if err != nil {
		log.Errorf("framework: failed to unlimit resource quotas(%v)", err)
		return err
	}

	log.Infof(
		"framework: resource quota us set successfully(%s %s)",
		h.Spec.Openstack.Project.Name,
		h.Spec.Openstack.Project.ID,
	)

	return nil
}

func (h *Helper) assignUserToProject() error {
	memberRoles, err := h.getMemberAndRoleIdPairs()
	if err != nil {
		log.Errorf("framework: failed to get member and role id pairs(%v)", err)
		return err
	}

	err = h.applyMembersToProject(memberRoles)
	if err != nil {
		log.Errorf("framework: failed to apply members to project(%v)", err)
		return err
	}

	log.Info("framework: users and roles added to project successfully")
	return nil
}

func (h *Helper) unlimitProjectResourceQuotas() error {
	err := h.Openstack.UpdateComputeQuotas(
		h.Spec.Openstack.Project.ID,
		h.genUnlimitedComputeQuota(),
	)
	if err != nil {
		return err
	}

	err = h.Openstack.UpdateNetworkQuotas(
		h.Spec.Openstack.Project.ID,
		h.genUnlimitedNetworkQuota(),
	)
	if err != nil {
		return err
	}

	err = h.Openstack.UpdateStorageQuotas(
		h.Spec.Openstack.Project.ID,
		h.genUnlimitedStorageQuota(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) genUnlimitedComputeQuota() computeQuota.UpdateOpts {
	unlimited := -1
	return computeQuota.UpdateOpts{
		Instances:          &unlimited,
		Cores:              &unlimited,
		RAM:                &unlimited,
		FixedIPs:           &unlimited,
		KeyPairs:           &unlimited,
		SecurityGroups:     &unlimited,
		SecurityGroupRules: &unlimited,
		FloatingIPs:        &unlimited,
	}
}

func (h *Helper) genUnlimitedNetworkQuota() networkQuotas.UpdateOpts {
	unlimited := -1
	return networkQuotas.UpdateOpts{
		FloatingIP:        &unlimited,
		Network:           &unlimited,
		Port:              &unlimited,
		Router:            &unlimited,
		SecurityGroup:     &unlimited,
		SecurityGroupRule: &unlimited,
		Subnet:            &unlimited,
		SubnetPool:        &unlimited,
	}
}

func (h *Helper) genUnlimitedStorageQuota() storageQuota.UpdateOpts {
	unlimited := -1
	return storageQuota.UpdateOpts{
		Gigabytes:          &unlimited,
		Volumes:            &unlimited,
		Snapshots:          &unlimited,
		Backups:            &unlimited,
		Groups:             &unlimited,
		PerVolumeGigabytes: &unlimited,
	}
}

func (h *Helper) getMemberAndRoleIdPairs() ([]memberAndRoleId, error) {
	memberRoles := []memberAndRoleId{}

	for _, r := range h.Spec.Openstack.Roles {
		role, err := h.Openstack.GetRoleByName(r.Name)
		if err != nil {
			return nil, err
		}

		member, err := h.Openstack.GetUserByName(r.User)
		if err != nil {
			return nil, err
		}

		memberRoles = append(
			memberRoles,
			memberAndRoleId{
				MemberId: member.ID,
				RoleId:   role.ID,
			},
		)
	}

	return memberRoles, nil
}

func (h *Helper) applyMembersToProject(memberRoles []memberAndRoleId) error {
	for _, memberRole := range memberRoles {
		err := h.Openstack.AddRole(
			memberRole.RoleId,
			roles.AssignOpts{
				UserID:    memberRole.MemberId,
				ProjectID: h.Spec.Openstack.Project.ID,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
