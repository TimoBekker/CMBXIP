package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

/*                    EMBEDED RESOURCES                       */

//go:embed _.glade
var builderResource string

/*                    GTK GUI COMPONENTS                      */

var app *gtk.Application
var builder *gtk.Builder
var mainWindow, settingsWindow *gtk.ApplicationWindow
var mainWindowDocumentURLEntry *gtk.Entry
var mainWindowTaskView *gtk.Box
var mainWindowTaskTitle *gtk.Entry
var mainWindowTaskDescription *gtk.TextView
var mainWindowExecutorsTree *gtk.TreeView
var mainWindowExecutorsTreeStore *gtk.TreeStore
var mainWindowExecutorsTreeRoleComboBoxStore *gtk.ListStore
var mainWindowExecutorsTreeRootEntry *ExecutorTreeEntry
var mainWindowNotifyLabel *gtk.Label

/*                         STRUCTS                            */

type ExecutorTreeEntry struct {
	FullName            string
	NameBackgroundColor string
	CompanyMediaID      string
	BitrixID            string
	BitrixRole          bitrixRole
	TreeIter            *gtk.TreeIter
	SubEntries          []*ExecutorTreeEntry
}

/*                       TYPES & ENUMS                        */

const (
	mainWindowExecutorsTreeColumnName = iota
	mainWindowExecutorsTreeColumnRole
	mainWindowExecutorsTreeColumnColor
)

type bitrixRole int

const (
	bxRoleExecutor bitrixRole = iota
	bxRoleAuditor
	bxRoleNone
)

var bxRoleName = map[bitrixRole]string{
	bxRoleExecutor: "Исполнитель",
	bxRoleAuditor:  "Наблюдатель",
	bxRoleNone:     "",
}
