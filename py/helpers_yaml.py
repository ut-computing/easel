import yaml

import external_tool

def read(filepath):
    with open(filepath) as f:
        return yaml.load(f, Loader=yaml.FullLoader)

def write(filepath, obj):
    with open(filepath, 'w') as f:
        f.write(yaml.dump(obj))

