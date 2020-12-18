CREATE TABLE courses (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    canvas_id integer NOT NULL,
    name text NOT NULL,
    code text NOT NULL,
    workflow_state text NOT NULL
);

CREATE TABLE pages (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    url text NOT NULL,
    title integer NOT NULL,
    created_at text NOT NULL,
    updated_at text NOT NULL,
    body text NOT NULL
    published integer NOT NULL,
    front_page integer NOT NULL
);

CREATE TABLE modules (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    canvas_id integer NOT NULL,
    position integer NOT NULL,
    name text NOT NULL,
    unlock_at text NOT NULL,
    require_sequential_progress boolean
    items_count integer NOT NULL,
    items_url text NOT NULL,
    published boolean not null
);

CREATE TABLE moduleitems (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    canvas_id integer not null,
    module_id integer not null,
    position integer not null,
    title text not null,
    indent integer not null,
    type text not null,
    content_id integer not null,
    html_url text not null,
    url text not null,
    page_url text not null,
    external_url text not null,
    new_tab boolean not null,
    completion_requirement_id integer not null,
    published boolean not null
);

CREATE TABLE moduleprerequisites (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    module_id INTEGER,
    prerequisite_module_id INTEGER
);
