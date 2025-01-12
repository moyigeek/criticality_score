/*
 * @Author: 7erry
 * @Date: 2024-12-02 16:34:35
 * @LastEditTime: 2024-12-02 17:17:32
 * @Description:
 */
package maven

// Request
type Response struct {
	Response       Resp           `json:"response"`
	ResponseHeader ResponseHeader `json:"responseHeader"`
}

type Resp struct {
	Docs     []Doc `json:"docs"`
	NumFound int64 `json:"numFound"`
	Start    int64 `json:"start"`
}

type Doc struct {
	A             string   `json:"a"`
	Ec            []string `json:"ec"`
	G             string   `json:"g"`
	ID            string   `json:"id"`
	LatestVersion string   `json:"latestVersion"`
	P             string   `json:"p"`
	RepositoryID  string   `json:"repositoryId"`
	Text          []string `json:"text"`
	Timestamp     int64    `json:"timestamp"`
	VersionCount  int64    `json:"versionCount"`
}

type ResponseHeader struct {
	Params Params `json:"params"`
	QTime  int64  `json:"QTime"`
	Status int64  `json:"status"`
}

type Params struct {
	Core            string `json:"core"`
	FL              string `json:"fl"`
	Indent          string `json:"indent"`
	Q               string `json:"q"`
	Rows            string `json:"rows"`
	Sort            string `json:"sort"`
	Spellcheck      string `json:"spellcheck"`
	SpellcheckCount string `json:"spellcheck.count"`
	Start           string `json:"start"`
	Version         string `json:"version"`
	Wt              string `json:"wt"`
}
