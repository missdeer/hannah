#pragma once

#include <cstdint>

#include "Sqlite3Helper.h"

struct sqlite3;

class Sqlite3DBManager
{
public:
    Sqlite3DBManager() : m_sqlite(m_db) {}
    ~Sqlite3DBManager();

    bool               open(const QString &dbPath, bool readOnly);
    bool               open(const std::string &dbPath, bool readOnly);
    bool               open(const char *dbPath, bool readOnly);
    [[nodiscard]] bool isOpened() const;
    bool               close();
    bool               save();
    bool               saveAs(const QString &newDbPath);
    bool               saveAs(const std::string &newDbPath);
    bool               saveAs(const char *newDbPath);
    bool               saveAndClose();
    bool               create(const QString &dbPath);
    bool               create(const std::string &dbPath);
    bool               create(const char *dbPath);
    bool               loadOrSaveInMemory(const QString &dbPath, bool isSave);

    Sqlite3Helper &engine();

private:
    sqlite3      *m_db {nullptr};
    Sqlite3Helper m_sqlite;
    std::string   m_savePoint;
    void          regenerateSavePoint();
    void          clearSavePoint();
    bool          setSavePoint();
    bool          releaseSavePoint();
};
