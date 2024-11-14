'''
Author: 7erry
Date: 2024-10-28 23:39:04
LastEditTime: 2024-11-14 14:18:47
Description: 
'''
import asyncio
import csv
import json
import httpx
from time import ctime

from config import(
    CRATES_IO_ENUMERATE_API_URL,
    CRATES_IO_OUTPUT_FILEPATH,
    CRATES_IO_ENUMERATE_PAGE,
    PER_PAGE,
    TIMEOUT
)

API_URL = CRATES_IO_ENUMERATE_API_URL

def parse(crate_list):
    '''
    parse json data
    '''
    result = []
    for crate in crate_list:
        result.append({
            "ID":crate.get("id",None),
            "Git Link":crate.get("repository",None),
            "Home Page":crate.get("homepage",None),
            "Documentation":crate.get("documentation",None),
            "Download Count":crate.get("downloads",None),
            "Created At":crate.get("created_at",None),
            "Updated At":crate.get("updated_at",None),
        })

    return result

async def get_repo_list(begin,end):
    '''
    get gitlab repo list
    '''
    repo_list = []
    async with httpx.AsyncClient() as client:
        for page in range(begin,end):
            try:
                resp = await client.get(
                    API_URL,
                    timeout=httpx.Timeout(timeout=TIMEOUT),
                    params={
                    "sort":"downloads",
                    "per_page":PER_PAGE,
                    "page":page
                    }
                )
            except Exception as e:
                print(f"[!] {e} {page} {ctime()}")
                continue
            content = json.loads(resp.text)
            repo_list += parse(content.get("crates",[]))
            page += 1
            if not content.get("meta",{}).get("next_page",None):
                break

    return repo_list

async def enumerate_crates_io():
    '''enumerate gitlab'''
    with open(CRATES_IO_OUTPUT_FILEPATH, mode='a', newline='', encoding='utf-8') as file:
        print(f"[*] Enumerating Page 1 - {CRATES_IO_ENUMERATE_PAGE}")
        repo_list = await get_repo_list(1,1 + CRATES_IO_ENUMERATE_PAGE)
        fieldnames = fieldnames = repo_list[0].keys()
        writer = csv.DictWriter(file, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(repo_list)

        index = CRATES_IO_ENUMERATE_PAGE + 1
        while True:
            print(f"[*] Enumerating Page {index} - {index + CRATES_IO_ENUMERATE_PAGE - 1}")
            repo_list = await get_repo_list(index,index + CRATES_IO_ENUMERATE_PAGE)
            writer.writerows(repo_list)
            index += CRATES_IO_ENUMERATE_PAGE
            if not repo_list:
                break

if __name__ == "__main__":
    print(f"[*] Enumerate at {ctime()}")
    try:
        asyncio.run(enumerate_crates_io())
        print(f"[*] Crates.io Enumerated at {ctime()}")
    except KeyboardInterrupt:
        print(f"[*] Enumerate Terminated {ctime()}")
