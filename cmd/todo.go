package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/devekkx/go-task-tracker/internal/display"
	"github.com/devekkx/go-task-tracker/internal/models"
	"github.com/devekkx/go-task-tracker/internal/storage"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage todo lists",
}

var todoCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new todo list",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		list := models.NewTodoList(name)
		store, err := storage.New()
		if err != nil { return err }
		if err := store.AddTodoList(list); err != nil { return err }
		display.Success("Todo list created: %s (ID: %s)", list.Name, list.ID)
		return nil
	},
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"ls"},
	Short: "List all todo lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		display.PrintTodoLists(store.ListTodoLists())
		return nil
	},
}

var todoShowCmd = &cobra.Command{
	Use:   "show <list-id>",
	Short: "Show items in a todo list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		list, err := store.GetTodoList(args[0])
		if err != nil { return err }
		display.PrintTodoList(list)
		return nil
	},
}

var todoAddCmd = &cobra.Command{
	Use:   "add <list-id> <content>",
	Short: "Add an item to a todo list",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		list, err := store.GetTodoList(args[0])
		if err != nil { return err }
		item := list.AddItem(strings.Join(args[1:], " "))
		if err := store.UpdateTodoList(list); err != nil { return err }
		display.Success("Item added: %s (ID: %s)", item.Content, item.ID)
		return nil
	},
}

var todoCheckCmd = &cobra.Command{
	Use:   "check <list-id> <item-id>",
	Short: "Mark a todo item as done",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		list, err := store.GetTodoList(args[0])
		if err != nil { return err }
		if err := list.CheckItem(args[1]); err != nil { return err }
		if err := store.UpdateTodoList(list); err != nil { return err }
		display.Success("Item checked.")
		return nil
	},
}

var todoUncheckCmd = &cobra.Command{
	Use:   "uncheck <list-id> <item-id>",
	Short: "Unmark a todo item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		list, err := store.GetTodoList(args[0])
		if err != nil { return err }
		if err := list.UncheckItem(args[1]); err != nil { return err }
		if err := store.UpdateTodoList(list); err != nil { return err }
		display.Success("Item unchecked.")
		return nil
	},
}

var todoRemoveCmd = &cobra.Command{
	Use:   "remove <list-id> <item-id>",
	Short: "Remove an item from a todo list",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		list, err := store.GetTodoList(args[0])
		if err != nil { return err }
		item, err := list.GetItem(args[1])
		if err != nil { return err }
		content := item.Content
		if err := list.RemoveItem(args[1]); err != nil { return err }
		if err := store.UpdateTodoList(list); err != nil { return err }
		display.Success("Item removed: %s", content)
		return nil
	},
}

var todoDeleteCmd = &cobra.Command{
	Use:     "delete <list-id>",
	Aliases: []string{"rm"},
	Short:   "Delete a todo list",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil { return err }
		list, err := store.GetTodoList(args[0])
		if err != nil { return err }
		name := list.Name
		if err := store.DeleteTodoList(args[0]); err != nil { return err }
		display.Success("Todo list deleted: %s", name)
		return nil
	},
}

func init() {
	todoCmd.AddCommand(todoCreateCmd)
	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoShowCmd)
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoCheckCmd)
	todoCmd.AddCommand(todoUncheckCmd)
	todoCmd.AddCommand(todoRemoveCmd)
	todoCmd.AddCommand(todoDeleteCmd)
}
