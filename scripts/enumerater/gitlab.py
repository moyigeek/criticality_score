'''
Author: 7erry
Date: 2024-10-28 23:39:04
LastEditTime: 2024-11-14 22:28:58
Description: 
'''
import asyncio
import csv
import json
from time import ctime,sleep

import httpx

from config import(
    GITLAB_ENUMERATE_API_URL,
    GITLAB_OUTPUT_FILEPATH,
    GITLAB_ENUMERATE_PAGE,
    PER_PAGE,
    TIMEOUT,
    TIME_INTERVAL
)

API_URL = GITLAB_ENUMERATE_API_URL

def parse(repo_list:list[dict]):
    '''
    parse json data
    '''
    result = []
    if repo_list and \
        len(repo_list) > 0 and \
        isinstance(repo_list,(dict,list)):
        for repo in repo_list:
            try:
                result.append({
                    "Repo Name":repo.get("path_with_namespace",None),
                    "Git Link":repo.get("web_url",None),
                    "Forks Count":repo.get("forks_count",None),
                    "Stars Count":repo.get("star_count",None),       
                    "Created At":repo.get("created_at",None),
                    "Updated At":repo.get("last_activity_at",None),
                })
            except AttributeError as e:
                print(f"[!] {e} {repo} {ctime()}")

    return result

async def get_repo_list(begin,end):
    '''
    get gitlab repo list
    '''
    repo_list = []
    async with httpx.AsyncClient() as client:
        for page in range(begin,end):
            sleep(TIME_INTERVAL)
            try:
                resp = await client.get(
                    API_URL,
                    timeout=httpx.Timeout(timeout=TIMEOUT),
                    params={
                    "order_by":"star_count",
                    "sort":"desc",
                    "per_page":PER_PAGE,
                    "page":page
                    }
                )
            except Exception as e:
                print(f"[!] {e} {page} {ctime()}")
                continue
            content = json.loads(resp.text)
            if not content:
                break
            repo_list += parse(content)
            page += 1

    return repo_list

async def enumerate_gitlab():
    '''enumerate gitlab'''
    with open(GITLAB_OUTPUT_FILEPATH, mode='a', newline='', encoding='utf-8') as file:
        print(f"[*] Enumerating Page 1 - {GITLAB_ENUMERATE_PAGE}")
        repo_list = await get_repo_list(1,1 + GITLAB_ENUMERATE_PAGE)
        fieldnames = fieldnames = repo_list[0].keys()
        writer = csv.DictWriter(file, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(repo_list)

        index = GITLAB_ENUMERATE_PAGE + 1
        while True:
            print(f"[*] Enumerating Page {index} - {index + GITLAB_ENUMERATE_PAGE - 1}")
            repo_list = await get_repo_list(index,index + GITLAB_ENUMERATE_PAGE)
            writer.writerows(repo_list)
            index += GITLAB_ENUMERATE_PAGE
            if not repo_list:
                break

if __name__ == "__main__":
    print(f"[*] Enumerate at {ctime()}")
    try:
        asyncio.run(enumerate_gitlab())
        print(f"[*] Gitlab Repositories Enumerated at {ctime()}")
    except KeyboardInterrupt:
        print(f"[*] Enumerate Terminated {ctime()}")
