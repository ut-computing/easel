import argparse

import commands
import helpers

parser = argparse.ArgumentParser(prog="easel", description="Easel - A Canvas "
        "course management tool.")
parser.add_argument('--api', action='store_true', help="report all API "
        "requests")
parser.add_argument('--api-dump', action='store_true', help="dump API request "
        "and response data")

# commands
subparsers = parser.add_subparsers(dest="command", help="the easel action to "
        "perform")
subparsers.required = True

## logging in
parser_login = subparsers.add_parser("login", help="login to Canvas")
parser_login.add_argument("hostname", help="the hostname of your Canvas "
        "instance")
parser_login.add_argument("token", help="your api token")
parser_login.set_defaults(func=commands.cmd_login)

## init
parser_init = subparsers.add_parser("init", help="initialize the db")
parser_init.set_defaults(func=commands.cmd_init)

## course commands
parser_course = subparsers.add_parser("course", help="course management "
        "commands")
parser_course.add_argument("subcommand", choices=["list", "add", "remove"])
parser_course.add_argument("subcommand_argument", nargs="?")
parser_course.set_defaults(func=commands.cmd_course)

args = parser.parse_args()

db = helpers.load_db()
args.func(db, args)
db.close()
