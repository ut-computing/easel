import yaml

import course
import helpers

EXT_TOOLS_PATH=course.COURSE_PATH+"/external_tools"
EXT_TOOLS_TABLE='external_tools'

class ExternalTool(yaml.YAMLObject):
    yaml_tag = "!ExternalTool"

    def __init__(self, name, consumer_key, shared_secret, config_type, config_url):
        self.name = name
        self.consumer_key = consumer_key
        self.shared_secret = shared_secret
        self.config_type = config_type
        self.config_url = config_url

    def __iter__(self):
        yield from vars(self).items()

    def __repr__(self):
        return f"ExternalTool(name={self.name})"

    def push(self, db, courses):
        # create only
        if not courses:
            courses = course.find_all(db)
        for course_ in courses:
            helpers.post(EXT_TOOLS_PATH.format(course_.canvas_id), self)
