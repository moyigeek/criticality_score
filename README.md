# Criticality Score Tool

## Quick Start

make sure `docker` is installed, and run the following.

```
./setup.sh
```

After the script finish, try to connect to database (the 
password is stored in `data/DB_PASSWD` and populate 
git_link fields in arch_packages and debian_packages 
manually and finally run following command.

```
docker compose exec app bash /gitlink.sh
```