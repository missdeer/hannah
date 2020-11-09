#include <QFile>
#include <QtCore>

#include "sqlite3helper.h"

Sqlite3Helper::~Sqlite3Helper()
{
    if (isDatabaseOpened())
        closeDatabaseConnection();
    m_db = nullptr;
}

void Sqlite3Helper::bind(sqlite3_stmt *pVM, int nParam, const QString &sValue)
{
    bind(pVM, nParam, sValue.toStdString().c_str());
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, const char *szParam, const char *szValue)
{
    int nParam = bindParameterIndex(mpVM, szParam);
    bind(mpVM, nParam, szValue);
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, const char *szParam, int const nValue)
{
    int nParam = bindParameterIndex(mpVM, szParam);
    bind(mpVM, nParam, nValue);
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, const char *szParam, int64_t const nValue)
{
    int nParam = bindParameterIndex(mpVM, szParam);
    bind(mpVM, nParam, nValue);
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, const char *szParam, double const dwValue)
{
    int nParam = bindParameterIndex(mpVM, szParam);
    bind(mpVM, nParam, dwValue);
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, const char *szParam, unsigned char const *blobValue, int nLen)
{
    int nParam = bindParameterIndex(mpVM, szParam);
    bind(mpVM, nParam, blobValue, nLen);
}

void Sqlite3Helper::bindNull(sqlite3_stmt *mpVM, const char *szParam)
{
    int nParam = bindParameterIndex(mpVM, szParam);
    bindNull(mpVM, nParam);
}

sqlite3_stmt *Sqlite3Helper::compile(const QString &szSQL)
{
    return compile(szSQL.toStdString().c_str());
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, int nParam, const char *szValue)
{
    int nRes = sqlite3_bind_text(mpVM, nParam, szValue, -1, SQLITE_TRANSIENT);
    if (nRes != SQLITE_OK)
    {
        qDebug() << "Error binding string param";
    }
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, int nParam, const int nValue)
{
    int nRes = sqlite3_bind_int(mpVM, nParam, nValue);
    if (nRes != SQLITE_OK)
    {
        qDebug() << "Error binding int param";
    }
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, int nParam, const int64_t nValue)
{
    int nRes = sqlite3_bind_int64(mpVM, nParam, nValue);
    if (nRes != SQLITE_OK)
    {
        qDebug() << "Error binding int param";
    }
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, int nParam, const double dValue)
{
    int nRes = sqlite3_bind_double(mpVM, nParam, dValue);
    if (nRes != SQLITE_OK)
    {
        qDebug() << "Error binding double param";
    }
}

void Sqlite3Helper::bind(sqlite3_stmt *mpVM, int nParam, const unsigned char *blobValue, int nLen)
{
    int nRes = sqlite3_bind_blob(mpVM, nParam, (const void *)blobValue, nLen, SQLITE_TRANSIENT);
    if (nRes != SQLITE_OK)
    {
        qDebug() << "Error binding blob param";
    }
}

void Sqlite3Helper::bindNull(sqlite3_stmt *mpVM, int nParam)
{
    int nRes = sqlite3_bind_null(mpVM, nParam);
    if (nRes != SQLITE_OK)
    {
        qDebug() << "Error binding NULL param";
    }
}

int Sqlite3Helper::bindParameterIndex(sqlite3_stmt *mpVM, const char *szParam)
{
    int nParam = sqlite3_bind_parameter_index(mpVM, szParam);
    if (!nParam)
    {
        qDebug() << "Parameter '" << szParam << "' is not valid for this statement";
    }

    return nParam;
}

void Sqlite3Helper::bind(sqlite3_stmt *pVM, const char *szParam, const QString &sValue)
{
    bind(pVM, szParam, sValue.toStdString().c_str());
}

sqlite3_stmt *Sqlite3Helper::compile(const char *szSQL)
{
    sqlite3_stmt *pVM = nullptr;
    int           res = sqlite3_prepare(m_db, szSQL, -1, &pVM, nullptr);
    if (res != SQLITE_OK)
    {
        qDebug() << "prepare SQL statement " << szSQL << " failed:" << res << (const char *)sqlite3_errmsg(m_db);
    }

    return pVM;
}

int Sqlite3Helper::execDML(const QString &szSQL)
{
    return execDML(szSQL.toStdString().c_str());
}

int Sqlite3Helper::execDML(const char *szSQL)
{
    int nRet;

    do
    {
        sqlite3_stmt *pVM = compile(szSQL);
        if (!pVM)
        {
            qDebug() << "VM null pointer";
            return -1;
        }
        nRet = sqlite3_step(pVM);

        if (nRet == SQLITE_ERROR)
        {
            qDebug() << (const char *)sqlite3_errmsg(m_db);
            sqlite3_finalize(pVM);
            break;
        }
        nRet = sqlite3_finalize(pVM);
    } while (nRet == SQLITE_SCHEMA);

    return nRet;
}

int Sqlite3Helper::execDML(sqlite3_stmt *pVM)
{
    if (!isDatabaseOpened())
    {
        qDebug() << "database is not opened";
        return -1;
    }
    if (!pVM)
    {
        qDebug() << "VM null pointer";
        return -1;
    }

    if (sqlite3_step(pVM) == SQLITE_DONE)
    {
        int nRowsChanged = sqlite3_changes(m_db);

        if (sqlite3_finalize(pVM) != SQLITE_OK)
        {
            qDebug() << (const char *)sqlite3_errmsg(m_db);
            return -1;
        }
        return nRowsChanged;
    }
    else
    {
        if (sqlite3_finalize(pVM) != SQLITE_OK)
        {
            qDebug() << (const char *)sqlite3_errmsg(m_db);
        }
    }
    return -1;
}

int Sqlite3Helper::execQuery(sqlite3_stmt *pVM, bool &eof)
{
    int nRet = sqlite3_step(pVM);

    if (nRet == SQLITE_DONE)
    { // no rows
        sqlite3_finalize(pVM);
        eof = true;
        return nRet;
    }
    else if (nRet == SQLITE_ROW)
    { // at least 1 row
        eof = false;
        return nRet;
    }
    nRet = sqlite3_finalize(pVM);

    if (nRet == SQLITE_SCHEMA)
        return nRet;

    qDebug() << (const char *)sqlite3_errmsg(m_db);
    return -1;
}

bool Sqlite3Helper::nextRow(sqlite3_stmt *pVM, bool &eof)
{
    int nRet = sqlite3_step(pVM);

    if (nRet == SQLITE_DONE)
    {
        // no rows
        if (sqlite3_finalize(pVM) != SQLITE_OK)
        {
            qDebug() << (const char *)sqlite3_errmsg(m_db);
        }
        eof = true;
    }
    else if (nRet == SQLITE_ROW)
    {
        // more rows, nothing to do
        eof = false;
    }
    else
    {
        if (sqlite3_finalize(pVM) != SQLITE_OK)
        {
            qDebug() << (const char *)sqlite3_errmsg(m_db);
            return false;
        }
    }
    return true;
}

bool Sqlite3Helper::isDatabaseOpened()
{
    return m_db != nullptr;
}

bool Sqlite3Helper::closeDatabaseConnection()
{
    if (!m_db)
    {
        qDebug() << "Database is not opened!";
        return false;
    }

    int result = sqlite3_close_v2(m_db);

    if (result != SQLITE_OK)
    {
        qDebug() << "Close DB failed: " << sqlite3_errmsg(m_db);
        return false;
    }

    qDebug() << "Database is closed!";
    m_db = nullptr;
    return true;
}

int Sqlite3Helper::checkTableIndexExists(const QString &field, const QString &name)
{
    int nRet = 0;
    do
    {
        sqlite3_stmt *pVM = compile("SELECT COUNT(*) FROM sqlite_master WHERE type = ? AND name = ?;");
        if (!pVM)
        {
            return -1;
        }
        bind(pVM, 1, field.toStdString().c_str()); // "table", "index"
        bind(pVM, 2, name.toStdString().c_str());
        bool eof = false;
        nRet     = execQuery(pVM, eof);
        if (nRet == SQLITE_DONE || nRet == SQLITE_ROW)
        {
            int size = 0;
            while (!eof)
            {
                size = (int)sqlite3_column_int(pVM, 0);
                if (!nextRow(pVM, eof))
                    break;
            }

            if (size > 0)
            {
                qDebug() << "found expected " << field << ":" << name;
                return size;
            }

            qDebug() << "not found expected " << field << ":" << name;
            return 0;
        }
    } while (nRet == SQLITE_SCHEMA);

    return 0;
}

bool Sqlite3Helper::openDatabase(const QString &name)
{
    if (sqlite3_open_v2(name.toStdString().c_str(), &m_db, SQLITE_OPEN_READWRITE, nullptr) != SQLITE_OK)
    {
        qDebug() << "Cannot open the database.";
        return false;
    }

    m_dbFile = name;
    qDebug() << "The database is opened.";
    return true;
}

bool Sqlite3Helper::createDatabase(const QString &name)
{
    if (sqlite3_open_v2(name.toStdString().c_str(), &m_db, SQLITE_OPEN_CREATE | SQLITE_OPEN_READWRITE, nullptr) != SQLITE_OK)
    {
        qDebug() << "Cannot create the database.";
        return false;
    }

    m_dbFile = name;
    qDebug() << "The database is created.";
    return true;
}

bool Sqlite3Helper::beginTransaction()
{
    return execDML("BEGIN TRANSACTION;") == SQLITE_OK;
}

bool Sqlite3Helper::endTransaction()
{
    return execDML("COMMIT TRANSACTION;") == SQLITE_OK;
}

bool Sqlite3Helper::rollbackTransaction()
{
    return execDML("ROLLBACK TRANSACTION;") == SQLITE_OK;
}

bool Sqlite3Helper::vacuum()
{
    return execDML("VACUUM;") == SQLITE_OK;
}
