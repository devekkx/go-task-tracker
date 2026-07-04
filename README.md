# go-task-tracker

A fast, offline-first CLI tool for managing tasks and todo lists.

## Features

- Track tasks with priorities, due dates, and tags
- Manage named todo lists with checkable items
- Filter and search tasks by status, priority, or tag
- Colorized, table-formatted output
- Local JSON storage at `~/.task-tracker/data.json`

## Installation

```bash
go install github.com/devekkx/go-task-tracker@latest
```

Or build from source:

```bash
git clone https://github.com/devekkx/go-task-tracker
cd go-task-tracker
make install
```

## Commands

| Command | Description |
|---|---|
| `tracker task add <title>` | Add a task |
| `tracker task list` | List all tasks |
| `tracker task show <id>` | Show task details |
| `tracker task done <id>` | Mark task as done |
| `tracker task start <id>` | Mark task as in-progress |
| `tracker task update <id>` | Update task fields |
| `tracker task delete <id>` | Delete a task |
| `tracker todo create <name>` | Create a todo list |
| `tracker todo list` | List all todo lists |
| `tracker todo show <id>` | Show list items |
| `tracker todo add <list-id> <item>` | Add an item |
| `tracker todo check <list-id> <item-id>` | Check an item |
| `tracker todo uncheck <list-id> <item-id>` | Uncheck an item |
| `tracker todo remove <list-id> <item-id>` | Remove an item |
| `tracker todo delete <list-id>` | Delete a todo list |
| `tracker stats` | Show statistics |
| `tracker version` | Show version |

## License

MIT

## Archive

Archive tasks you want to keep but hide from the default view:

```bash
tracker task archive <id>    # archive a task
tracker task unarchive <id>  # restore an archived task
tracker task list --archived # list archived tasks
```

## Export & Import

Back up or migrate your data with JSON export and import (v1.1.0):

```bash
tracker export output.json
tracker import backup.json
```

## What's New in v1.1.0

- Archive tasks you want to keep but hide from the default view
- Sort task lists by title, priority, due date, or creation date
- Copy any task with `task copy`
- Mark filtered tasks done in bulk with `task bulk-done`
- Export and import data as JSON for backup or migration
- Machine-readable output with `--json` flag
