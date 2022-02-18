package vpn

import (
	"path/filepath"
	"strings"

	"github.com/c3os-io/c3os/installer/config"
	"github.com/c3os-io/c3os/installer/systemd"
	"github.com/c3os-io/c3os/installer/utils"
)

func Setup(instance, apiAddress, rootDir string, start bool, c *config.Config) error {
	svc, err := systemd.EdgeVPN(instance, rootDir)
	if err != nil {
		return err
	}

	apiAddress = strings.ReplaceAll(apiAddress, "https://", "")
	apiAddress = strings.ReplaceAll(apiAddress, "http://", "")

	vpnOpts := map[string]string{
		"EDGEVPNTOKEN":         c.C3OS.NetworkToken,
		"API":                  "true",
		"APILISTEN":            apiAddress,
		"EDGEVPNLOWPROFILEVPN": "true",
		"DHCP":                 "true",
		"DHCPLEASEDIR":         "/usr/local/.c3os/lease",
	}
	// Override opts with user-supplied
	for k, v := range c.VPN {
		vpnOpts[k] = v
	}
	// Setup edgevpn instance
	err = utils.WriteEnv(filepath.Join(rootDir, "/etc/systemd/system.conf.d/edgevpn-c3os.env"), vpnOpts)
	if err != nil {
		return err
	}

	err = svc.WriteUnit()
	if err != nil {
		return err
	}

	if start {
		err = svc.Start()
		if err != nil {
			return err
		}

		return svc.Enable()
	}
	return nil
}
