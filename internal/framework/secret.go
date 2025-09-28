package framework

import "encoding/base64"

func (h *Helper) base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
