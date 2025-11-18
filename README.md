# go-task-tracker

A fast, offline-first CLI tool for managing tasks and todo lists.

## Features

- Track tasks with priorities, due dates, and tags
- Manage todo lists with items
- Filter and search tasks
- Colorized table output
- Local JSON storage at `~/.task-tracker/`

## Usage

```bash
tracker task add "Fix the bug" --priority high --due 2025-12-01
tracker task list --status pending
tracker todo create "Shopping"
tracker todo add <list-id> "Buy milk"
```

## Installation

```bash
go install github.com/devekkx/go-task-tracker@latest
```
