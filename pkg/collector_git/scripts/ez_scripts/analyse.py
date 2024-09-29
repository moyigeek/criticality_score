'''
Date: 2024-09-11 17:33:55
LastEditTime: 2024-09-29 17:02:56
Description: 
'''
import pandas as pd
import sweetviz as sv
from . import config


grm_data = pd.read_csv(config.CSV_PATH)

my_report = sv.analyze(grm_data)
my_report.show_html() # Default arguments will generate to "SWEETVIZ_REPORT.html"
