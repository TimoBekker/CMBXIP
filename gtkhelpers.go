package main

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <gio/gio.h>
// #include <glib.h>
// #include <glib-object.h>
// #include <gtk/gtk.h>
// #include "gtkhelpers.go.h"
import "C"

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

// fix GtkCellRendererCombo bug from gotk3 issue #688 https://github.com/gotk3/gotk3/issues/688

func setCellRendererComboModel(cellRenderer *gtk.CellRendererCombo, model *gtk.TreeModel) {
	C.set_cell_renderer_combo_model(unsafe.Pointer(cellRenderer.GObject), unsafe.Pointer(model.GObject))
}
