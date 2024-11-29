package api

const (
	PER_PAGE      = 100
	TIME_INTERVAL = 0.7
	TIMEOUT       = 1000

	BITBUCKET_ENUMERATE_API_URL = "https://api.bitbucket.org/2.0/repositories"
	GITLAB_ENUMERATE_API_URL    = "https://gitlab.com/api/v4/projects"
	GITEE_ENUMERATE_API_URL     = "https://api.indexea.com/v1/search/widget/wjawvtmm7r5t25ms1u3d"
	CRATES_IO_ENUMERATE_API_URL = "https://crates.io/api/v1/crates"
	PYPI_API_URL                = "https://pypi.org/simple/"

	GITLAB_TOTAL_PAGES = 100000

	BITBUCKET_ENUMERATE_PAGE = 40 //* repo_num = page * 10
	GITLAB_ENUMERATE_PAGE    = 20 //* repo_num = page * 100
	GITEE_ENUMERATE_PAGE     = 20 //* repo_num = page * 100
	CRATES_IO_ENUMERATE_PAGE = 20
)
