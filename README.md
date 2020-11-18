# Easel

A Canvas course management tool.

## TODO

- Still need to think about multiple courses (i.e., sections).
    - A single file should potentially describe multiple courses, but then we
      can't put the course id in the file
    - Instead of one init command, I'm thinking there should be
      `easel course <course_url>` and `easel reset` commands.
- Enforce directories? (e.g., pages, assignments, modules)
- Canvas uses html for pages and assignment descriptions, we should use markdown
  and convert it.
- How to represent dates?
    - API requires strings (e.g., "2013-01-23T23:59:00-07:00")
    - For relative dates, we should include times that are relative to the
      section meeting time (e.g., beginning of class, end of class)
- would it be worth adding in grading stuff eventually?
- Some fields would be useful to Easel but not necessary for instructor edits
  (e.g., record ids, component status).
  Do we keep those in the DB but not write them to file?
