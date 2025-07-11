# Сборка

Для сборки окружения [MXE](https://github.com/mxe/mxe) (необхоимо для сборки приложения) использовать команду:
```
time make gtk3 MXE_TARGETS='x86_64-w64-mingw32.static' --jobs=4 JOBS=4
```

Последний коммит с успешной сборкой: `a31368b037221d56bcfcfd8c546b89adffe9ea04`.

Если проблемы сборки из-за недостающих символов C++ cairo - добавить в cairo.pc `-lstdc++`.

Для ручного обновления версии GTK3 в MXE необходимо править gtk3.mk, с соответствующей правкой хеша.

Соответственно не забываем прописывать окружение PATH для обнаружения программ MXE окружения, к примеру в `~/.profile`:
```
PATH=/opt/mxe/usr/bin:$PATH
```

Сборка самого приложения происходит через утилиту make: `make windows_amd64`

# Баги

В windows сборке GTK3 есть проблемы с контекстным меню при которых приложение вылетает. Фикс был в этом [коммите](https://gitlab.gnome.org/GNOME/gtk/-/merge_requests/5690/diffs?commit_id=22b091047f6a71670e0bcaad24c0ca5109a07279), соответственно нормально работает всё с версии `3.24.38`.
