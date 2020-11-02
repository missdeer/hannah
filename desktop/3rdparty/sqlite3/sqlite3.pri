
DEFINES += SQLITE_ENABLE_FTS5=1

INCLUDEPATH += $$PWD
SOURCES += $$PWD/sqlite3.c
HEADERS += $$PWD/sqlite3.h \
           $$PWD/sqlite3ext.h
