#pragma once

#include <cstdint>
#include <map>
#include <string>

#include "Sqlite3Statement.h"

struct sqlite3;
struct sqlite3_stmt;
struct sqlite3_context;
struct sqlite3_value;

class Sqlite3Helper
{
public:
    explicit Sqlite3Helper(sqlite3 *&db);

    Sqlite3StatementPtr compile(const char *szSQL);
    Sqlite3StatementPtr compile(const std::string &sql);
    Sqlite3StatementPtr compile(const QString &sql);
    int                 execDML(const char *szSQL);
    int                 execDML(const std::string &sql);
    int                 execDML(const QString &sql);
    static bool         isQueryOk(int result);
    static bool         isOk(int result);
    static bool         canQueryLoop(int result);

    int getQueryResult(const char *szSQL);
    int getQueryResult(const std::string &sql);
    int getQueryResult(const QString &sql);

    bool isDatabaseOpened();

    int checkTableOrIndexExists(const std::string &field, const std::string &name);
    int checkTableOrIndexExists(const QString &field, const QString &name);

    bool checkTableColumnExists(const QString &tableName, const QString &field);
    bool addTableColumn(const QString &tableName, const QString &fieldName, const QString &fieldType);
    int  getTableColumnCount(const QString &tableName);
    int  getTableRowCount(const QString &tableName);

    bool createTablesAndIndexes(std::map<std::string, const char *> &tablesMap, std::map<std::string, const char *> &indexesMap);

    bool beginTransaction();
    bool endTransaction();
    bool rollbackTransaction();

    bool vacuum();

    std::int64_t lastInsertRowId();

    void registerCustomFunctions();

private:
    sqlite3 *&m_db;
};
