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

- assignments
- assignment groups
- courses (mainly just grabs the syllabus)
- modules
- pages
- quizzes

```
easel pull [component_type] [component_id]
```

E.g.,

```
easel pull pages lesson-1-variables
```

Pulls a single item of the given component type from the assigned courses. Works
for the same components as previously listed.

### Push

```
easel push
```

Pushes everything to each of the configured courses (individual courses coming).
A push reads the information of each component stored locally and for each one,
makes a PUT request to Canvas.

```
easel push [component_type]
```

E.g.,

```
easel push pages
```

Reads each page in the component type's directory and pushes them to each
configured course in Canvas. Currently works for the following components:

- pages
- more to come!

```
easel push [component_type] [component_id]
```

E.g.,

```
easel push pages lesson-1-variables
```

Reads and pushes a single item of the given component type to the configured
courses. Works for the same components as previously listed.

## File Structure

Component files are stored in separate directories, named for their component
type (e.g., pages are stored in a directory called `pages`). This is required
for now, but may be more flexible in the future.

Each individual component is defined by a single file. Quizzes are an exception
as the questions of the quiz are store in a `<quiz_name>_questions.md` file that
is separate from the main `<quiz_name>.md` file.

Most components are defined using yaml with some associated body/description
content. The yaml should be defined in a fenced code block (used for
preformatted text) at the very beginning of the file. The component's body
content should immediately follow. Here is an example:

~~~
```
url: lesson-1-introduction
title: Lesson 1- Introduction
created_at: "2020-07-16T22:45:56Z"
updated_at: "2020-08-10T18:14:42Z"
published: true
front_page: false
editing_roles: teachers
todo_date: "2020-08-28T11:58:00-06:00"
```
In this lesson, you'll learn the basics of programs and programming. You'll get
a chance to write some short programs to escape a maze.

## <i class="icon-check-plus" aria-hidden="true"></i> Outcomes

- Describe major course concepts and policies
- Use Canvas to access course materials
- Define "programming"
- Explain why programming is useful
- Identify limitations of programming
- Enumerate the core components of programs
- Write simple programs

## <i class="icon-quiz" aria-hidden="true"></i> Activities

- [Slides 1- Introduction](https://dixie.instructure.com/courses/615446/pages/slides-1-introduction)
- [Quiz 1- Introduction](https://dixie.instructure.com/courses/615446/assignments/7701035)
- [Programming 01 - Maze Game](https://dixie.instructure.com/courses/615446/assignments/7782062)

## <i class="icon-educators" aria-hidden="true"></i> Summary</h2>

Programs are created by programmers to give computers instructions. Data comes
in, the computer processes it, and then returns the result to the user.
~~~

Note: Canvas prefers the body content to be in html, even though we prefer to
edit in markdown. For now, we are having an issue with this when pulling a
component from Canvas as it wants to overwrite your nicely written markdown with
the raw html it uses. Converting from html to markdown is inconsistent at best.
So I recommend not pulling a component if you've defined its content in markdown
locally. Hopefully we'll have a solution for this in the future.

## TODO

- I've been assuming user pulls or pushes from the course's root directory. Need
  to search for the component dirs
- multiple courses (i.e., sections).
    - implicit iteration
        - push: pushes to all courses, unless specified (e.g., -c 02)
        - pull: pulls from all courses, checks for and reports any differences
            - need to add a prompt for overwrite, manually merge, or abort
            - need to track multiple canvas ids per component in the db. I'm
              saving the canvas id on each component as if it would be the same
              across all courses, but this is not the case.
- Figure out the workflow for editing page/assignment content. Canvas uses html,
  I'd prefer to express it in markdown.
  - First proposal: locally in markdown, convert to html when pushing. Don't
    edit content in Canvas (since we can't faithfully convert html to md).
    Pulling would not overwrite the component's contents.
- pull/push everything in transactions
    - use db as intermediate step, only go to Canvas if db transaction succeeded
    - workflow for pulling whether to overwrite, manually merge, or abort
    - When pushing, update database with result (e.g., when pushing to a new
      course, the canvas id will be different)
- represent dates as time.Time
    - API requires strings in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ (e.g., "2013-01-23T23:59:00-07:00")
- add a progress bar for pushing and pulling

### Thoughts

- Enforce directories? (e.g., pages, assignments, modules)
    - Or when pushing a component, save its filepath in the db
- Component files that only have yaml (no md or html), should the extension be
  yaml or stay consistent with md?
- We should enable expressing dates/times that are relative to the section
  meeting time (e.g., beginning of class, end of class, Fridays)
- would it be worth adding in grading stuff eventually?
- Some fields would be useful to Easel but not necessary for instructor edits
  (e.g., record ids, component status).
  Do we keep those in the DB but not write them to file?
- should quiz questions be in their own file? Options:
    - a single quiz's questions in one file. easier to implement but it would be
      harder to reuse them
    - one file per question, easy to move around, but how to uniquely identify
      each question? (for the name of the file)
    - one file per question category (e.g., all requrements engineering
      questions) this is probably the best user-focused approach, but harder to
      implement?
