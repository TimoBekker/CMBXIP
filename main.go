package main

import (
	"fmt"
	"log"
	"nikeron/cmbxip/config"
	"os"
	"runtime"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	config.DebugLogger().Print("подготовка приложения")
	var err error
	app, err = gtk.ApplicationNew("ru.digitalreg.cmbxip", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		panic(err)
	}
	app.Connect("activate", appActivateSignal)
}

func main() {
	config.DebugLogger().Print("запуск приложения")
	defer func() {
		if err := recover(); err != nil {
			var data string
			switch v := err.(type) {
			case error:
				data = v.Error()
			default:
				data = fmt.Sprintf("%v", v)
			}
			data += "\n\n"

			buf := make([]byte, 1024*1024*10)
			data += string(buf[:runtime.Stack(buf, true)])

			os.WriteFile(config.PanicDumpPath(), []byte(data), 0644)
		}
	}()
	app.Run(config.NonFlagArgs())
}

func appActivateSignal() {
	var err error
	if builder, err = gtk.BuilderNewFromString(builderResource); err != nil {
		log.Fatalln("Couldn't make builder:", err)
	}

	{ // settings window
		settingsWindowObj, err := builder.GetObject("settings-window")
		if err != nil {
			log.Fatalln("Couldn't get settings window")
		}
		settingsWindow = settingsWindowObj.(*gtk.ApplicationWindow)
		settingsWindow.Connect("delete-event", func(window *gtk.ApplicationWindow) bool {
			window.Hide()
			return true
		})
		app.AddWindow(settingsWindow)
	}

	{ // main window
		mainWindowObj, err := builder.GetObject("main-window")
		if err != nil {
			log.Fatalln("Couldn't get main window")
		}
		mainWindow = mainWindowObj.(*gtk.ApplicationWindow)
		mainWindow.Connect("destroy", app.Quit)
		curTitle, _ := mainWindow.GetTitle()
		mainWindow.SetTitle(curTitle + " v" + config.Version())
		app.AddWindow(mainWindow)
	}

	{ // main window cm url entry
		mainWindowDocumentURLEntryObj, err := builder.GetObject("main-url-entry")
		if err != nil {
			log.Fatalln("Couldn't get cm doc url entry")
		}
		mainWindowDocumentURLEntry = mainWindowDocumentURLEntryObj.(*gtk.Entry)
	}

	{ // main window task view
		mainWindowTaskViewObj, err := builder.GetObject("main-task-view")
		if err != nil {
			log.Fatalln("Couldn't get task view")
		}
		mainWindowTaskView = mainWindowTaskViewObj.(*gtk.Box)
	}

	{ // main window task title
		mainWindowTaskTitleObj, err := builder.GetObject("main-task-title")
		if err != nil {
			log.Fatalln("Couldn't get task title")
		}
		mainWindowTaskTitle = mainWindowTaskTitleObj.(*gtk.Entry)
	}

	{ // main window task description
		mainWindowTaskDescriptionObj, err := builder.GetObject("main-task-description")
		if err != nil {
			log.Fatalln("Couldn't get task description")
		}
		mainWindowTaskDescription = mainWindowTaskDescriptionObj.(*gtk.TextView)
	}

	{ // main window executors tree
		mainWindowExecutorsTreeObj, err := builder.GetObject("main-executors-tree")
		if err != nil {
			log.Fatalln("Couldn't get executors tree")
		}
		mainWindowExecutorsTree = mainWindowExecutorsTreeObj.(*gtk.TreeView)
		mainWindowExecutorsTree.Connect("row-collapsed", func(t *gtk.TreeView) bool {
			t.ExpandAll()
			return true
		})
		{ // name column
			cellRenderer, _ := gtk.CellRendererTextNew()
			column, _ := gtk.TreeViewColumnNewWithAttribute("Полное имя", cellRenderer, "text", mainWindowExecutorsTreeColumnName)
			column.AddAttribute(cellRenderer, "cell-background", mainWindowExecutorsTreeColumnColor)
			column.SetExpand(true)
			column.SetResizable(true)
			mainWindowExecutorsTree.AppendColumn(column)
		}
		{ // role column
			cellRenderer, _ := gtk.CellRendererComboNew()
			mainWindowExecutorsTreeRoleComboBoxStore, _ = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
			for k, v := range bxRoleName {
				mainWindowExecutorsTreeRoleComboBoxStore.Set(mainWindowExecutorsTreeRoleComboBoxStore.Append(),
					[]int{0, 1}, []interface{}{k, v})
			}
			cellRenderer.SetProperty("text-column", 1)
			cellRenderer.SetProperty("has-entry", false)
			cellRenderer.SetProperty("editable", true)
			setCellRendererComboModel(cellRenderer, mainWindowExecutorsTreeRoleComboBoxStore.ToTreeModel())
			cellRenderer.Connect("changed", mainWindowExecutorsTreeRoleCellChanged)
			column, _ := gtk.TreeViewColumnNewWithAttribute("Роль", cellRenderer, "text", mainWindowExecutorsTreeColumnRole)
			column.SetExpand(true)
			column.SetResizable(true)
			column.SetReorderable(true)
			column.AddAttribute(cellRenderer, "cell-background", mainWindowExecutorsTreeColumnColor)
			mainWindowExecutorsTree.AppendColumn(column)
		}
		mainWindowExecutorsTreeStore, _ = gtk.TreeStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
		mainWindowExecutorsTree.SetModel(mainWindowExecutorsTreeStore)
	}

	{ // main window notify label
		mainWindowNotifyLabelObj, err := builder.GetObject("main-notify-label")
		if err != nil {
			log.Fatalln("Couldn't get notify text")
		}
		mainWindowNotifyLabel = mainWindowNotifyLabelObj.(*gtk.Label)
		mainWindowNotifyLabel.SetText("Введите ссылку на документ Company Media!")
	}

	builder.ConnectSignals(map[string]interface{}{
		"main-window-settings-button-clicked": settingsWindow.Present,
		"settings-window-shown":               func() { settingsWindowAction(0) },
		"settings-window-save-button-clicked": func() {
			settingsWindowAction(1)
			defer config.Save()
			settingsWindow.Hide()
		},
		"settings-window-bx-root-autofind-toggled": func() { settingsWindowAction(2) },
		"main-url-entry-changed":                   mainWindowCMDocumentURLEntryChanged,
		"main-import-button-clicked":               mainWindowImportButtonClicked,
	})

	mainWindow.Show()
}
