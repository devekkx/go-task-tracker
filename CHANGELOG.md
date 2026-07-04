# Changelog

## [1.1.0] - 2026-07-04

### Added
- Task archiving: `task archive` and `task unarchive` commands
- Task sorting: `--sort` flag supports title, priority, due, created
- Task copy: `task copy <id>` duplicates a task
- Task bulk-done: `task bulk-done` marks all filtered tasks as done
- Task clear-done: `task clear-done` removes completed tasks
- JSON export (`ExportJSON`) and import (`ImportJSON`) for data portability
- Stats: ArchivedTasks and PendingTodoItems fields
- `--archived`, `--limit`, `--json` flags on `task list`
- TodoList: Description field, SetDescription and Rename methods
- Task: HasTag, RemoveTag, ClearTags, TagsCount, Archive, Unarchive

### Fixed
- Duplicate godoc comments removed throughout the codebase
- Task title validation now trims whitespace before checking empty

### Changed
- Default task list now excludes archived tasks
- Stats output shows archived task and pending todo item counts

## [1.0.0] - 2025-12-31

- Initial release
