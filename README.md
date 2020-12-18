# Easel

A Canvas course management tool.

## Operations

### Login

```
easel login <canvas_base_url> <api_token>
```

E.g.,

```
easel login https://dixie.instructure.com aASDFoSf23kD0aS9fAByuA0fyA0yf8e9ha
```

Only needs to be run once per client machine. Records the Canvas url and token
to be used for later.

### Init

```
easel init
```

Initializes the easel database in the current directory. Run this one time per
course directory.

### Course

```
easel course <canvas_course_url>
```

E.g.,

```
easel course https://dixie.instructure.com/courses/615446
```

Hooks up the database to a Canvas course. Run this one time per Canvas course
(once per section taught per semester).

```
easel course list
```

List all Canvas courses that are tracked in the database.

### Pull

```
easel pull
```

Pulls (most) everything from the configured courses. A pull is defined as
getting the information from Canvas and persisting it locally (db + file).

```
easel pull [component_type]
```

E.g.,

```
easel pull pages
```

Pulls all items of a single component type from the assigned courses. Works for
the following components:

- courses (mainly just grabs the syllabus)
- pages
- more to come!

```
easel pull [component_type] [component_id]
```

E.g.,

```
easel pull pages lesson1-variables
```

Pulls a single item of the given component type from the assigned courses. Works
for the same components as previously listed.

### Push

## TODO

- I've been assuming user pulls or pushes from the course's root directory. Need
  to search for the component dirs
- multiple courses (i.e., sections).
    - implicit iteration
        - push: pushes to all courses, unless specified (e.g., -c 02)
        - pull: pulls from all courses, checks for and reports any differences
            - need to add a prompt for overwrite, manually merge, or abort
- Figure out the workflow for editing page/assignment content. Canvas uses html,
  I'd prefer to express it in markdown.
  - First proposal: locally in markdown, convert to html when pushing. Don't
    edit content in Canvas (since we can't faithfully convert html to md).
    Pulling would not overwrite the component's contents.
- pull/push everything in transactions
    - use db as intermediate step, only go to Canvas if db transaction succeeded
    - workflow for pulling whether to overwrite, manually merge, or abort
- represent dates as time.Time
    - API requires strings (e.g., "2013-01-23T23:59:00-07:00")
- add a progress bar for pushing and pulling
- When pushing, update database with result (e.g., when pushing to a new course,
  the canvas id will be different)

### Thoughts

- Enforce directories? (e.g., pages, assignments, modules)
- We should enable expressing dates/times that are relative to the section
  meeting time (e.g., beginning of class, end of class)
- would it be worth adding in grading stuff eventually?
- Some fields would be useful to Easel but not necessary for instructor edits
  (e.g., record ids, component status).
  Do we keep those in the DB but not write them to file?
