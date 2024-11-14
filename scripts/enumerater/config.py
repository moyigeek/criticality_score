'''
Author: 7erry
Date: 2024-10-25 15:23:12
LastEditTime: 2024-11-14 22:28:19
Description: 
'''
HEADERS = {
    "User-Agent"       : \
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) "+\
            "AppleWebKit/537.36 (KHTML, like Gecko) "+\
            "Chrome/84.0.4147.89 Safari/537.36" ,
}

PER_PAGE = 100
TIME_INTERVAL = 0.7
TIMEOUT = 1000

PYPI_SOURCE_FILE_PATH = "pypi.json"

BITBUCKET_OUTPUT_FILEPATH = "bitbucket.csv"
GITLAB_OUTPUT_FILEPATH = "gitlab.csv"
GITEE_OUTPUT_FILEPATH = "gitee.csv"
CRATES_IO_OUTPUT_FILEPATH = "crates_io.csv"

BITBUCKET_ENUMERATE_API_URL = "https://api.bitbucket.org/2.0/repositories"
GITLAB_ENUMERATE_API_URL = "https://gitlab.com/api/v4/projects"
GITEE_ENUMERATE_API_URL = "https://api.indexea.com/v1/search/widget/wjawvtmm7r5t25ms1u3d"
CRATES_IO_ENUMERATE_API_URL = "https://crates.io/api/v1/crates"
PYPI_API_URL = "https://pypi.org/simple/"
DEPS_DEV = {
    "PYPI": "https://deps.dev/pypi/"
}

BITBUCKET_ENUMERATE_PAGE = 40 #* repo_num = page * 10
GITLAB_ENUMERATE_PAGE = 20 #* repo_num = page * 100
GITEE_ENUMERATE_PAGE = 20 #* repo_num = page * 100
CRATES_IO_ENUMERATE_PAGE = 20
