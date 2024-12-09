# Criticality Score

[简体中文](./README.zh_CN.md) | English


## Description

## Difference from [ossf/criticality_score](https://github.com/ossf/criticality_score)

1. Use distribution dependents 
   Fetch content by collecting distributions 
2. Support all git repos rather than github only
   Fetch content by cloning git repo
3. Approach to collect metrics  
4. Do not use google cloud and Big Query
5. Easy to deploy
6. Provide extra information, e.g. relationships

## Quick Start

make sure `docker` is installed, and run the following.

```sh
./setup.sh
```

After the script finish, try to connect to database (the 
password is stored in `data/DB_PASSWD` and populate 
git_link fields in arch_packages and debian_packages 
manually and finally run following command).

```sh
docker compose exec app bash /gitlink.sh
```

## Documentation

See `docs/` for details

## Reference

1. [ossf/criticality_score](https://github.com/ossf/criticality_score)
