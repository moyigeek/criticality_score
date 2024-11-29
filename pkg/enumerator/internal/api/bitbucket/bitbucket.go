package bitbucket

import (
	"encoding/json"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

const (
	API_URL = "https://api.bitbucket.org/2.0/repositories"
)

func Fetch(c *req.Client, url string) (*Response, error) {
	resp := &Response{}
	res, err := c.R().Get(url)

	if err != nil || res.GetStatusCode() != 200 {
		logrus.Errorf(
			"[Bitbucket] failed: code=%d, msg=%s, err=%v",
			res.GetStatusCode(),
			res.String(),
			err,
		)
		return nil, err
	}

	json.Unmarshal(res.Bytes(), resp)

	return resp, nil
}
