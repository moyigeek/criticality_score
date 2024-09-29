'''
Date: 2024-09-08 03:18:46
LastEditTime: 2024-09-29 17:03:16
Description: Run Collector in Multi Process Way
'''
import asyncio
from subprocess import Popen

from tqdm import tqdm
from . import config

async def get_git_list(filepath: str) -> list[str]:
    "Read git link list from file"
    with open(filepath,"r",encoding="utf-8") as f:
        git_list: list[str] = f.readlines()
    return git_list

async def run_task(git_list: list[str],begin: int,end: int):
    "run collector"
    cmd = config.COLLECT_LOCAL
    temp_path = f"{config.TEMP_DIR}/{begin}~{end}.csv"
    with open(temp_path,"w",encoding="utf-8") as f:
        f.writelines(git_list)
    with open(config.OUTPUT_PATH,"ab+") as f:
        with open(config.ERR_PATH,"ab+") as e:
            process = Popen(f"{cmd} {temp_path}",shell=True,stdout=f,stderr=e)
    output, err = process.communicate()
    if err:
        print(err.decode("utf-8"))
    if output:
        with open(config.OUTPUT_PATH,"wb") as f:
            f.write(output)

async def main():
    "entrance of this script"
    git_list = await get_git_list(config.GIT_LIST)
    list_len = len(git_list)
    task_list = []

    for i in tqdm(range(0,list_len,config.TASK_SIZE),desc="Launching Collectors"):
        if i + config.TASK_SIZE < list_len:
            task = asyncio.create_task(run_task(git_list[i:i+config.TASK_SIZE],i,i+config.TASK_SIZE))
        else:
            task = asyncio.create_task(run_task(git_list[i:list_len],i,list_len))
        task_list.append(task)

    for t in tqdm(task_list,desc="Waiting for collectors"):
        await t

if __name__ == "__main__":
    asyncio.run(main())
