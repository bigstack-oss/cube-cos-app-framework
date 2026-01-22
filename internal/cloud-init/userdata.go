package cloudinit

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	header = "#cloud-config\n"
)

type UserData struct {
	Bootcmd    []string    `yaml:"bootcmd,omitempty"`
	WriteFiles []WriteFile `yaml:"write_files"`
	RunCmd     []string    `yaml:"runcmd"`
}

type WriteFile struct {
	Owner       string `yaml:"owner"`
	Encoding    string `yaml:"encoding"`
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Append      bool   `yaml:"append"`
	Permissions string `yaml:"permissions"`
}

func (u *UserData) AddBootCmd(bootCmd string) {
	u.Bootcmd = append(u.Bootcmd, bootCmd)
}

func (u *UserData) AddRunCmd(runCmd string) {
	u.RunCmd = append(u.RunCmd, runCmd)
}

func (u *UserData) AddWriteFile(writeFile WriteFile) {
	u.WriteFiles = append(u.WriteFiles, writeFile)
}

func (u *UserData) GenConfig() string {
	b, err := yaml.Marshal(u)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s%s", header, string(b))
}

func (u *UserData) GenConfigBytes() []byte {
	return []byte(u.GenConfig())
}
