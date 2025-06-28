package rancher

import "github.com/bigstack-oss/bigstack-dependency-go/pkg/terraform"

var (
	Token = ""
	Url   = ""
	User  = ""
)

func InitGlobalAuthIdentities() error {
	h := terraform.GetGlobalHelper()
	values, err := h.ShowResourceValues("rancher2_bootstrap")
	if err != nil {
		return err
	}

	for key, value := range values {
		switch key {
		case "url":
			Url = value.(string)
		case "user":
			User = value.(string)
		case "token":
			Token = value.(string)
		}
	}

	return nil
}
