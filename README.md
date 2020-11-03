# Easel

A Canvas course management tool.

## TODO

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
