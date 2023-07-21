
#include <chrono>
#include <vector>
#include <sqlite3.h>

#include <QFile>
#include <QtCore/qstringliteral.h>
#include <QtCore>

#include "Sqlite3Helper.h"
#include "Sqlite3Constants.h"

Sqlite3Helper::Sqlite3Helper(sqlite3 *&db) : m_db(db) {}

Sqlite3StatementPtr Sqlite3Helper::compile(const char *szSQL)
{
    if (!m_db)
    {
        qCritical() << szSQL << "null db pointer";
        return nullptr;
    }

    sqlite3_stmt *pVM = nullptr;
    int           res = sqlite3_prepare(m_db, szSQL, -1, &pVM, nullptr);
    if (res != SQLITE_OK)
    {
        const char *szError = sqlite3_errmsg(m_db);
        qCritical() << "prepare SQL statement" << szSQL << "failed:" << res << "-" << szError;
        return nullptr;
    }

    return std::make_shared<Sqlite3Statement>(m_db, pVM);
}

Sqlite3StatementPtr Sqlite3Helper::compile(const std::string &sql)
{
    return compile(sql.c_str());
}

Sqlite3StatementPtr Sqlite3Helper::compile(const QString &sql)
{
    return compile(sql.toStdString().c_str());
}

int Sqlite3Helper::execDML(const char *szSQL)
{
    int resultCode = SQLITE_OK;

    do
    {
        auto stmt = compile(szSQL);
        if (!stmt || !stmt->isValid())
        {
            qCritical() << szSQL << "null pVM, quit";
            return SQLITE_ERROR;
        }

        resultCode = sqlite3_step(stmt->m_pVM);

        if (resultCode == SQLITE_ERROR)
        {
            qCritical() << "step SQL statement " << szSQL << "failed:" << resultCode << "-" << sqlite3_errmsg(m_db);
            sqlite3_finalize(stmt->m_pVM);
            break;
        }
        resultCode = sqlite3_finalize(stmt->m_pVM);
    } while (resultCode == SQLITE_SCHEMA);

    return resultCode;
}

int Sqlite3Helper::execDML(const std::string &sql)
{
    return execDML(sql.c_str());
}

int Sqlite3Helper::execDML(const QString &sql)
{
    return execDML(sql.toStdString().c_str());
}

bool Sqlite3Helper::isQueryOk(int result)
{
    return (result == SQLITE_DONE || result == SQLITE_ROW);
}

bool Sqlite3Helper::isOk(int result)
{
    return result == SQLITE_OK;
}

bool Sqlite3Helper::canQueryLoop(int result)
{
    return result == SQLITE_SCHEMA;
}

int Sqlite3Helper::getQueryResult(const char *szSQL)
{
    int  count = 0;
    int  nRet  = 0;
    bool eof   = false;
    do
    {
        auto stmt = compile(szSQL);
        nRet      = stmt->execQuery(eof);
        if (nRet == SQLITE_DONE || nRet == SQLITE_ROW)
        {
            while (!eof)
            {
                count = stmt->getInt(0);
                stmt->nextRow(eof);
            }
            break;
        }
    } while (nRet == SQLITE_SCHEMA);
    return count;
}

int Sqlite3Helper::getQueryResult(const std::string &sql)
{
    return getQueryResult(sql.c_str());
}

int Sqlite3Helper::getQueryResult(const QString &sql)
{
    return getQueryResult(sql.toStdString().c_str());
}

bool Sqlite3Helper::isDatabaseOpened()
{
    return m_db != nullptr;
}

int Sqlite3Helper::checkTableOrIndexExists(const std::string &field, const std::string &name)
{
    bool eof  = false;
    int  nRet = 0;
    do
    {
        auto stmt = compile(SQL_STATEMENT_SELECT_SQLITE_MASTER);
        if (!stmt || !stmt->isValid())
        {
            return -1;
        }
        stmt->bind(1, field.c_str());
        stmt->bind(2, name.c_str());
        nRet = stmt->execQuery(eof);
        if (nRet == SQLITE_DONE || nRet == SQLITE_ROW)
        {
            int size = 0;
            while (!eof)
            {
                size = 1;
                if (!stmt->nextRow(eof))
                {
                    return -2;
                }
            }

            if (size > 0)
            {
                qDebug() << "found expected" << QString::fromStdString(field) << ":" << QString::fromStdString(name);
                return size;
            }

            qDebug() << "not found expected" << QString::fromStdString(field) << ":" << QString::fromStdString(name);
            return 0;
        }
    } while (nRet == SQLITE_SCHEMA);

    return 0;
}

int Sqlite3Helper::checkTableOrIndexExists(const QString &field, const QString &name)
{
    return checkTableOrIndexExists(field.toStdString(), name.toStdString());
}

bool Sqlite3Helper::createTablesAndIndexes(std::map<std::string, const char *> &tablesMap, std::map<std::string, const char *> &indexesMap)
{
    for (auto &[tableName, sql] : tablesMap)
    {
        int res = checkTableOrIndexExists("table", tableName);
        if (res < 0)
        {
            // do repair database file
            return false;
        }

        if (!res)
        {
            if (execDML(sql) != SQLITE_OK)
            {
                qWarning() << sqlite3_errmsg(m_db);
                return false;
            }
        }
    }

    for (auto &[indexName, sql] : indexesMap)
    {
        int res = checkTableOrIndexExists("index", indexName);
        if (res < 0)
        {
            // do repair database file
            return false;
        }

        if (!res)
        {
            if (execDML(sql) != SQLITE_OK)
            {
                qWarning() << sqlite3_errmsg(m_db);
                return false;
            }
        }
    }

    return true;
}

bool Sqlite3Helper::beginTransaction()
{
    return execDML(SQL_STATEMENT_BEGIN_TRANSACTION) == SQLITE_OK;
}

bool Sqlite3Helper::endTransaction()
{
    return execDML(SQL_STATEMENT_COMMIT_TRANSACTION) == SQLITE_OK;
}

bool Sqlite3Helper::rollbackTransaction()
{
    return execDML(SQL_STATEMENT_ROLLBACK_TRANSACTION) == SQLITE_OK;
}

bool Sqlite3Helper::vacuum()
{
    return execDML(SQL_STATEMENT_VACUUM) == SQLITE_OK;
}

std::int64_t Sqlite3Helper::lastInsertRowId()
{
    if (m_db)
    {
        return sqlite3_last_insert_rowid(m_db);
    }
    return -1;
}

bool Sqlite3Helper::addTableColumn(const QString &tableName, const QString &fieldName, const QString &fieldType)
{
    QString sql = QStringLiteral("ALTER TABLE %1 ADD COLUMN %2 %3;").arg(tableName, fieldName, fieldType);
    return execDML(sql) == SQLITE_OK;
}

bool Sqlite3Helper::checkTableColumnExists(const QString &tableName, const QString &fieldName)
{
    sqlite3_stmt *stmt = nullptr;
    QString       sql  = QStringLiteral("SELECT %2 FROM %1 LIMIT 1;").arg(tableName, fieldName);

    int resultCode = sqlite3_prepare_v2(m_db, sql.toStdString().c_str(), -1, &stmt, nullptr);

    if (resultCode != SQLITE_OK)
    {
        qCritical() << "Error preparing statement: " << sqlite3_errmsg(m_db);
        return false;
    }

    int columnIndex = sqlite3_column_count(stmt);
    sqlite3_finalize(stmt);

    return (columnIndex >= 1);
}

int Sqlite3Helper::getTableColumnCount(const QString &tableName)
{
    sqlite3_stmt *stmt = nullptr;
    QString       sql  = QStringLiteral("SELECT * FROM %1;").arg(tableName);

    int resultCode = sqlite3_prepare_v2(m_db, sql.toStdString().c_str(), -1, &stmt, nullptr);

    if (resultCode != SQLITE_OK)
    {
        qCritical() << "Error preparing statement: " << sqlite3_errmsg(m_db);
        return false;
    }

    int columnCount = sqlite3_column_count(stmt);
    sqlite3_finalize(stmt);

    return columnCount;
}

int Sqlite3Helper::getTableRowCount(const QString &tableName)
{
    return getQueryResult(QStringLiteral("SELECT COUNT(*) FROM %1;").arg(tableName));
}
