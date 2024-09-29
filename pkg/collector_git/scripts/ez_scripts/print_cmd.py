'''
Date: 2024-09-27 22:11:28
LastEditTime: 2024-09-29 17:03:09
Description: 
'''
from subprocess import Popen

from tqdm import tqdm
from . import config

def get_git_list(filepath: str) -> list[str]:
    "Read git link list from file"
    with open(filepath,"r",encoding="utf-8") as f:
        git_list: list[str] = f.readlines()
    return git_list

def run_task(git_list: list[str],begin: int,end: int):
    "run collector"
    cmd = config.CLONE_REPOS
    temp_path = f"{config.TEMP_DIR}/{begin}~{end}.csv"
    with open(temp_path,"w",encoding="utf-8") as f:
        f.writelines(git_list)
        
    return f"nohup {cmd} {temp_path} >  ./{begin}~{end}.txt & &&\n"

def main():
    "entrance of this script"
    git_list = get_git_list(config.GIT_LIST)
    list_len = len(git_list)
    script = "#! /bin/sh\n"

    for i in range(0,list_len,config.TASK_SIZE):
        if i + config.TASK_SIZE < list_len:
            new_line = run_task(git_list[i:i+config.TASK_SIZE],i,i+config.TASK_SIZE)
        else:
            new_line = run_task(git_list[i:list_len],i,list_len)
        script += new_line

    print(script)

if __name__ == "__main__":
    main()
