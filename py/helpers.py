import json
import os
from pathlib import Path
import requests
import sqlite3
import sys
import tinydb

API="/api/v1"
SCHEMA="https://"

def write_config(hostname, token):
    home = Path.home()
    if home == "":
        print("home directory is not set")
        sys.exit(1)

    config_file = home / ".easelrc" # https://docs.python.org/3.7/library/pathlib.html#operators

    config = {"hostname": hostname, "token": token}
    try:
        with open(config_file, 'x') as f:
            f.write(json.dumps(config)) # TODO: 0644
    except FileExistsError:
        print(f"Config file {config_file} exists")

def load_db():
    return tinydb.TinyDB(".easeldb")

def setup_directories():
    dirs = ["assignment_groups", "assignments", "external_tools", "modules",
            "pages", "quizzes"]
    for d in dirs:
        os.mkdir(d)

def get(path, params={}, decode=True):
    conf = Config()
    if not path.startswith("/"):
        print("request path must start with /")
        sys.exit(1)
    r = requests.get(SCHEMA+conf.hostname+path,
            headers={'Authorization': 'Bearer '+conf.token},
            params=params)
    if r.status_code in [200]:
        if decode:
            r = r.json()
    else:
        print("=== REQUEST ERROR: {} ===".format(r.status_code))
        print(r.json())
        print("==========================")
    return r

class Config:

    def __init__(self):
        home = Path.home()
        if home == "":
            print("home directory is not set")
            sys.exit(1)

        config_file = home / ".easelrc" # https://docs.python.org/3.7/library/pathlib.html#operators

        f = open(config_file)
        c = json.loads(f.read())
        f.close()
        self.hostname = c["host"]
        self.token = c["token"]
