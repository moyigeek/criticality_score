#!/usr/bin/env python3

import logging
import os
from argparse import ArgumentParser

import psycopg2

# TODO: import from csv file
parser = ArgumentParser(description="Copy git link from existing database")

parser.add_argument("--user", help="Source database user", default="postgres")
parser.add_argument(
    "--database", help="Source database name", default="criticality_score"
)
parser.add_argument("--password", help="Source database password", required=True)
parser.add_argument("--host", help="Source database host", required=True)
parser.add_argument("--port", help="Source database port", default=5432)
parser.add_argument("--debug", help="Debug mode", action="store_true")


args = parser.parse_args()
logging.basicConfig(level=logging.INFO if not args.debug else logging.DEBUG)

logging.info("Connecting to database...")

conn = psycopg2.connect(
    database=args.database,
    user=args.user,
    password=args.password,
    host=args.host,
    port=args.port,
)

script_path = os.path.dirname(os.path.realpath(__file__))

env_file = open(script_path + "/../.env").read()

# get DB_HOST_PORT

envs = env_file.split("\n")

env_map = {}
for env in envs:
    if not env:
        continue
    strs = env.split("=", 2)
    env_map[strs[0]] = strs[1].strip("\"'")

logging.info("Connecting to dest database...")
conn_dest = psycopg2.connect(
    database="criticality_score",
    user="postgres",
    password=env_map["DB_PASSWD"],
    host="localhost",
    port=env_map["DB_HOST_PORT"],
)


def copy_table(table_name):
    curr = conn.cursor()
    logging.info("Fetching data from source database...")
    curr.execute("SELECT package, git_link FROM {}".format(table_name))
    while True:
        row = curr.fetchone()
        if row is None:
            break

        curr_dst = conn_dest.cursor()
        logging.debug("Updating package %s", row[0])
        curr_dst.execute(
            "UPDATE {} SET git_link = %s WHERE package = %s".format(table_name),
            (row[1], row[0]),
        )
        if curr_dst.rowcount == 0:
            logging.warning("No row updated for package %s", row[0])
        curr_dst.close()
    conn_dest.commit()
    curr.close()

logging.info("Copying debian_packages...")
copy_table("debian_packages")
logging.info("Copying arch_packages...")
copy_table("arch_packages")

conn.close()
conn_dest.close()
