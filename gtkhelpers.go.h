#pragma once

#include <stdint.h>
#include <stdlib.h>
#include <string.h>

static void set_cell_renderer_combo_model(void *cellRenderer, void *model) {
    g_object_set(GTK_CELL_RENDERER_COMBO(cellRenderer), "model", GTK_TREE_MODEL(model), NULL);
}
