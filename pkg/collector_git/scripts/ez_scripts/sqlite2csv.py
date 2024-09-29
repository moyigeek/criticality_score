'''
Date: 2024-09-11 17:06:25
LastEditTime: 2024-09-29 17:03:02
Description: 
'''
import csv
import sqlite3

from . import config

def sqlite2csv():
    '''save data from sqlite into csv'''
    conn = sqlite3.connect(config.SQLITE_DB_PATH)

    cursor = conn.execute(
        "select * from read_metrics"
    )

    with open(config.CSV_PATH, 'w',encoding="utf-8",newline='') as file:
        writer = csv.writer(file)
        writer.writerow([i[0] for i in cursor.description])
        writer.writerows(cursor)

    conn.close()

if __name__ == "__main__":
    sqlite2csv()
