'''
Author: 7erry
Date: 2024-10-25 15:25:17
LastEditTime: 2024-11-14 22:28:01
Description:
'''
import asyncio
import csv
import json
from time import ctime,sleep
import httpx

from config import (
    BITBUCKET_ENUMERATE_PAGE,
    BITBUCKET_ENUMERATE_API_URL,
    BITBUCKET_OUTPUT_FILEPATH,
    TIMEOUT,
    TIME_INTERVAL
)

API_URL = BITBUCKET_ENUMERATE_API_URL

def parse(repo_list):
    '''
    parse json data
    '''
    result = []
    if repo_list and len(repo_list) >= 1 and \
        isinstance(repo_list,(list,dict)) and \
        isinstance(repo_list[0],dict):
        for repo in repo_list:
            if not repo.get("is_private",None) and repo.get("type",None) == "repository":
                result.append({
                    "Repo Name":repo.get("full_name",None),
                    "API URL":repo.get("links",{}).get("self",{}).get("href",None),
                    "Git Link":repo.get("links",{}).get("html",{}).get("href",None),
                })

    return result

async def get_repo_list(page:int):
    '''get repo list of bitbucket'''
    repo_list = []
    global API_URL
    async with httpx.AsyncClient() as client:
        for _ in range(page):
            sleep(TIME_INTERVAL)
            try:
                resp = await client.get(
                    API_URL,
                    timeout=httpx.Timeout(timeout=TIMEOUT)
                )
            except Exception as e:
                print(f"[!] {e} {ctime()}")
                resp = None
            if not resp or not resp.content or not resp.text or len(resp.text) <= 2:
                continue
            try:
                contents = json.loads(resp.text)
            except json.decoder.JSONDecodeError as e:
                print(f"[!] {e} {resp.text} {ctime()}")
                continue
            repo_list += parse(contents.get("values",[]))
            next_page = contents.get("next",None)
            API_URL = next_page
            if not next_page:
                break

    return repo_list

async def enumerate_bitbucket():
    '''enumerate bitbucket'''
    index = 1
    with open(BITBUCKET_OUTPUT_FILEPATH, mode='a', newline='', encoding='utf-8') as file:
        repo_list = await get_repo_list(1)
        fieldnames = fieldnames = repo_list[0].keys()
        writer = csv.DictWriter(file, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(repo_list)

        while True:
            print(f"[*] Enumerating Page {index} - {index + BITBUCKET_ENUMERATE_PAGE}")
            repo_list = await get_repo_list(BITBUCKET_ENUMERATE_PAGE)
            writer.writerows(repo_list)
            index += BITBUCKET_ENUMERATE_PAGE
            if not API_URL:
                break

if __name__ == "__main__":
    print(f"[*] Enumerate at {ctime()}")
    try:
        asyncio.run(enumerate_bitbucket())
        print(f"[*] Bitbucket Repositories Enumerated at {ctime()}")
    except KeyboardInterrupt:
        print(f"[*] Enumerate Terminated {ctime()}")
