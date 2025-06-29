package configs

type Rancher struct {
	Url  string `json:"url"`
	Auth `json:"auth"`
}
