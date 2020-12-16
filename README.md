# Easel

A Canvas course management tool.

## Operations

### Login

```
easel login https://dixie.instructure.com apitoken
```

Only needs to be run once. Records the hostname and token to be used for later.

### Init

```
easel init
```

Initializes the easel database. Run this one time per course directory.

### Course

```
easel course https://dixie.instructure.com/courses/615446
```

Hooks up the database to a Canvas course. Run this one time per Canvas course
(typically once or twice per semester).

```
easel course list
```

List all Canvas courses that are tracked in the database.

### Pull

### Push

## TODO

- Still need to think about multiple courses (i.e., sections).
    - implicit iteration
        - push: pushes to all courses, unless specified
        - pull: pulls from all courses, checks for and reports any differences
    - A single file should potentially describe multiple courses, but then we
      can't put the course id in the file
- Enforce directories? (e.g., pages, assignments, modules)
- Canvas uses html for pages and assignment descriptions, we should use markdown
  and convert it. For now I'm just dumping the html though.
- How to represent dates?
    - API requires strings (e.g., "2013-01-23T23:59:00-07:00")
    - For relative dates, we should include times that are relative to the
      section meeting time (e.g., beginning of class, end of class)
- would it be worth adding in grading stuff eventually?
- Some fields would be useful to Easel but not necessary for instructor edits
  (e.g., record ids, component status).
  Do we keep those in the DB but not write them to file?
