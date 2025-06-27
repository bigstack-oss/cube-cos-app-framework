package configs

type Rancher struct {
	Url  string `yaml:"url"`
	Auth `yaml:"auth"`
}
