package framework

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/harbor"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/wait"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	log "go-micro.dev/v5/logger"
)

// harborReadyTimeout bounds how long we wait for the framework's Harbor to be
// reachable through its ingress Octavia LB, which provisions asynchronously
// (LB + amphora + floating IP) and can take several minutes on a fresh cluster.
const harborReadyTimeout = 15 * time.Minute

// waitForHarborReady polls Harbor's systeminfo endpoint until it answers 200 or
// the timeout elapses. A fixed short sleep is not enough on a fresh cluster: the
// ingress LB isn't routing yet, so the registry project/service-account calls
// would fail and the framework would come up with no usable registry.
func waitForHarborReady(httpUrl string) error {
	url := strings.TrimRight(httpUrl, "/") + "/api/v2.0/systeminfo"
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	deadline := time.Now().Add(harborReadyTimeout)
	for {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("harbor registry ingress %s not ready after %s", url, harborReadyTimeout)
		}
		log.Infof("harbor: waiting for registry ingress %s to become reachable...", url)
		wait.Seconds(15)
	}
}

func (h *Helper) createRegistryProject() error {
	h.setVipToPrimaryDnsServer()
	defer h.restoreOriginalDnsList()

	access := h.getCubeAppsAccess()
	if err := waitForHarborReady(access.HttpUrl); err != nil {
		return err
	}
	cli, err := harbor.NewHelper(
		harbor.Url(access.HttpUrl),
		harbor.Username(access.Username),
		harbor.Password(access.Password),
		harbor.InsecureSkipVerify(true),
	)
	if err != nil {
		log.Errorf("harbor: failed to create harbor client(%v)", err)
		return err
	}

	_, err = cli.CreateProject("extensions")
	if err != nil {
		if strings.Contains(err.Error(), "createProjectConflict") {
			return nil
		}

		log.Errorf("harbor: failed to create project for extensions(%s)", err.Error())
		return err
	}

	return nil
}

func (h *Helper) createRegistryServiceAccount() error {
	h.setVipToPrimaryDnsServer()
	defer h.restoreOriginalDnsList()

	access := h.getCubeAppsAccess()
	if err := waitForHarborReady(access.HttpUrl); err != nil {
		return err
	}
	cli, err := harbor.NewHelper(
		harbor.Url(access.HttpUrl),
		harbor.Username(access.Username),
		harbor.Password(access.Password),
		harbor.InsecureSkipVerify(true),
	)
	if err != nil {
		log.Errorf("harbor: failed to create harbor user(%v)", err)
		return err
	}

	_, err = cli.CreateUser("appctl", h.Spec.Openstack.Auth.Password, "appctl@registry.local")
	if err != nil {
		if strings.Contains(err.Error(), "createUserConflict") {
			return nil
		}

		log.Errorf("harbor: failed to create user for registry(%s)", err.Error())
		return err
	}

	users, err := cli.ListUsers(1, 10)
	if err != nil {
		log.Errorf("harbor: failed to list users(%v)", err)
		return err
	}

	userID := int64(0)
	for _, user := range users.Payload {
		if user.Username == "appctl" {
			userID = user.UserID
			break
		}
	}

	_, err = cli.SetUserSysAdmin(userID)
	if err != nil {
		log.Errorf("harbor: failed to set sysadmin for registry user(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) getCubeAppsAccess() configs.ExtensionRepo {
	for _, repo := range configs.DefaultSpec.Framework.ExtensionRepos {
		if repo.Name == "cube-apps" {
			repo.HttpUrl = h.findCubeAppsHttpUrl()
			return repo
		}
	}

	return configs.ExtensionRepo{}
}
