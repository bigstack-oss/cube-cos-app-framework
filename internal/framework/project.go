package framework

import (
	storageQuota "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/quotasets"
	computeQuota "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	networkQuotas "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
)

type memberAndRoleId struct {
	MemberId string
	RoleId   string
}

func (h *Helper) applyProjectAndResourceQuota() error {
	var err error
	h.Project, err = h.Openstack.CreateProject(h.Project.Name)
	if err != nil {
		return err
	}

	err = h.applyUnlimitedResourceQuotas()
	if err != nil {
		return err
	}

	h.Log.Infof(
		"project and resource quota set successfully (%s %s)",
		h.Project.Name,
		h.Project.ID,
	)
	return nil
}

func (h *Helper) addUserAndRolesToProject() error {
	memberRoles, err := h.getMemberAndRoleIdPairs()
	if err != nil {
		return err
	}

	err = h.applyMembersToProject(memberRoles)
	if err != nil {
		return err
	}

	h.Log.Info("users and roles added to project successfully")
	return nil
}

func (h *Helper) applyUnlimitedResourceQuotas() error {
	err := h.Openstack.UpdateComputeQuotas(
		h.Project.ID,
		h.genUnlimitedComputeQuota(),
	)
	if err != nil {
		return err
	}

	err = h.Openstack.UpdateNetworkQuotas(
		h.Project.ID,
		h.genUnlimitedNetworkQuota(),
	)
	if err != nil {
		return err
	}

	err = h.Openstack.UpdateStorageQuotas(
		h.Project.ID,
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

	for _, r := range h.Config.Openstack.Roles {
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
				ProjectID: h.Project.ID,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
