package terraform

type State struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Type      string     `json:"type"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Attributes `json:"attributes"`
}

type Attributes struct {
	User  string `json:"user"`
	Token string `json:"token"`
	Url   string `json:"url"`
}
