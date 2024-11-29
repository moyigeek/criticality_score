/*
 * @Author: 7erry
 * @Date: 2024-11-29 17:34:47
 * @LastEditTime: 2024-11-29 17:38:02
 * @Description:
 */
package gitlab

import (
	"encoding/json"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

func Fetch(c *req.Client, url string) (*Response, error) {
	resp := &Response{}
	res, err := c.R().Get(url)

	if err != nil || res.GetStatusCode() != 200 {
		logrus.Errorf(
			"[Gitlab] failed: code=%d, msg=%s, err=%v",
			res.GetStatusCode(),
			res.String(),
			err,
		)
		return nil, err
	}

	json.Unmarshal(res.Bytes(), resp)
	return resp, nil
}
