
#include <sqlite3.h>

#include <QtCore>

#include "Sqlite3Statement.h"

Sqlite3Statement::Sqlite3Statement(sqlite3 *&database, sqlite3_stmt *pVM) : m_db(database), m_pVM(pVM) {}

void Sqlite3Statement::bind(int nParam, const QString &value) {}

void Sqlite3Statement::bind(const char *szParam, const char *szValue, int nLen /*= -1*/)
{
    int nParam = bindParameterIndex(szParam);
    bind(nParam, szValue, nLen);
}

void Sqlite3Statement::bind(const char *szParam, int nValue)
{
    int nParam = bindParameterIndex(szParam);
    bind(nParam, nValue);
}

void Sqlite3Statement::bind(const char *szParam, int64_t nValue)
{
    int nParam = bindParameterIndex(szParam);
    bind(nParam, nValue);
}

void Sqlite3Statement::bind(const char *szParam, double dwValue)
{
    int nParam = bindParameterIndex(szParam);
    bind(nParam, dwValue);
}

void Sqlite3Statement::bind(const char *szParam, const unsigned char *blobValue, int nLen)
{
    int nParam = bindParameterIndex(szParam);
    bind(nParam, blobValue, nLen);
}

void Sqlite3Statement::bind(const char *szParam, std::string_view szValue)
{
    bind(szParam, szValue.data(), static_cast<int>(szValue.size()));
}

void Sqlite3Statement::bind(const char *szParam, const std::string &szValue)
{
    bind(szParam, szValue.c_str(), static_cast<int>(szValue.size()));
}

void Sqlite3Statement::bindNull(const char *szParam)
{
    int nParam = bindParameterIndex(szParam);
    bindNull(nParam);
}

void Sqlite3Statement::bind(const std::string &szParam, const std::string &szValue)
{
    bind(szParam.c_str(), szValue.c_str(), static_cast<int>(szValue.size()));
}

void Sqlite3Statement::bind(const std::string &szParam, int nValue)
{
    bind(szParam.c_str(), nValue);
}

void Sqlite3Statement::bind(const std::string &szParam, int64_t nValue)
{
    bind(szParam.c_str(), nValue);
}

void Sqlite3Statement::bind(const std::string &szParam, double dwValue)
{
    bind(szParam.c_str(), dwValue);
}

void Sqlite3Statement::bind(const std::string &szParam, const unsigned char *blobValue, int nLen)
{
    bind(szParam.c_str(), blobValue, nLen);
}

void Sqlite3Statement::bind(const std::string &szParam, std::string_view szValue)
{
    bind(szParam.c_str(), szValue.data(), static_cast<int>(szValue.size()));
}

void Sqlite3Statement::bindNull(const std::string &szParam)
{
    bindNull(szParam.c_str());
}

int Sqlite3Statement::bindParameterIndex(const QString &szParam)
{
    return bindParameterIndex(szParam.toStdString().c_str());
}

void Sqlite3Statement::bind(const QString &szParam, const QString &value)
{
    bind(szParam.toStdString().c_str(), value.toStdString().c_str());
}

void Sqlite3Statement::bind(const QString &szParam, const std::string &szValue)
{
    bind(szParam.toStdString().c_str(), szValue.c_str());
}

void Sqlite3Statement::bind(const QString &szParam, std::string_view szValue)
{
    bind(szParam.toStdString().c_str(), szValue);
}

void Sqlite3Statement::bind(const QString &szParam, int nValue)
{
    bind(szParam.toStdString().c_str(), nValue);
}

void Sqlite3Statement::bind(const QString &szParam, int64_t nValue)
{
    bind(szParam.toStdString().c_str(), nValue);
}

void Sqlite3Statement::bind(const QString &szParam, double dwValue)
{
    bind(szParam.toStdString().c_str(), dwValue);
}

void Sqlite3Statement::bind(const QString &szParam, const unsigned char *blobValue, int nLen)
{
    bind(szParam.toStdString().c_str(), blobValue, nLen);
}

void Sqlite3Statement::bindNull(const QString &szParam)
{
    bindNull(szParam.toStdString().c_str());
}

void Sqlite3Statement::bind(int nParam, const char *szValue, int nLen /*= -1*/)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return;
    }
    int nRes = sqlite3_bind_text(m_pVM, nParam, szValue, nLen, SQLITE_TRANSIENT);
    if (nRes != SQLITE_OK)
    {
        qCritical() << "Error binding string param";
    }
}

void Sqlite3Statement::bind(int nParam, const std::string &szValue)
{
    bind(nParam, szValue.c_str(), static_cast<int>(szValue.size()));
}

void Sqlite3Statement::bind(int nParam, std::string_view szValue)
{
    bind(nParam, szValue.data(), static_cast<int>(szValue.size()));
}

void Sqlite3Statement::bind(int nParam, int nValue)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return;
    }
    int nRes = sqlite3_bind_int(m_pVM, nParam, nValue);
    if (nRes != SQLITE_OK)
    {
        qCritical() << "Error binding int param";
    }
}

void Sqlite3Statement::bind(int nParam, int64_t nValue)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return;
    }
    int nRes = sqlite3_bind_int64(m_pVM, nParam, nValue);
    if (nRes != SQLITE_OK)
    {
        qCritical() << "Error binding int param";
    }
}

void Sqlite3Statement::bind(int nParam, double dValue)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return;
    }
    int nRes = sqlite3_bind_double(m_pVM, nParam, dValue);
    if (nRes != SQLITE_OK)
    {
        qCritical() << "Error binding double param";
    }
}

void Sqlite3Statement::bind(int nParam, const unsigned char *blobValue, int nLen)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return;
    }
    int nRes = sqlite3_bind_blob(m_pVM, nParam, (const void *)blobValue, nLen, SQLITE_TRANSIENT);
    if (nRes != SQLITE_OK)
    {
        qCritical() << "Error binding blob param";
    }
}

void Sqlite3Statement::bindNull(int nParam)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return;
    }
    int nRes = sqlite3_bind_null(m_pVM, nParam);

    if (nRes != SQLITE_OK)
    {
        qCritical() << "Error binding NULL param";
    }
}

int Sqlite3Statement::bindParameterIndex(const char *szParam)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return -1;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return -1;
    }
    int nParam = sqlite3_bind_parameter_index(m_pVM, szParam);

    if (!nParam)
    {
        qCritical() << "parameter is not valid for this statement";
    }

    return nParam;
}

void Sqlite3Statement::bind(const char *szParam, const QString &value)
{
    bind(szParam, value.toStdString().c_str());
}

int Sqlite3Statement::bindParameterIndex(const std::string &szParam)
{
    return bindParameterIndex(szParam.c_str());
}

void Sqlite3Statement::bind(const std::string &szParam, const QString &value)
{
    bind(szParam.c_str(), value.toStdString().c_str());
}

int Sqlite3Statement::execDML()
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return -1;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return -1;
    }

    if (sqlite3_step(m_pVM) == SQLITE_DONE)
    {
        int nRowsChanged = sqlite3_changes(m_db);

        if (sqlite3_finalize(m_pVM) != SQLITE_OK)
        {
            qCritical() << sqlite3_errmsg(m_db);
            return -1;
        }
        return nRowsChanged;
    }

    sqlite3_finalize(m_pVM);
    qCritical() << sqlite3_errmsg(m_db);

    return -1;
}

int Sqlite3Statement::execQuery(bool &eof)
{
    int nRet = sqlite3_step(m_pVM);

    if (nRet == SQLITE_DONE)
    {
        // no rows
        sqlite3_finalize(m_pVM);
        eof = true;
        return nRet;
    }

    if (nRet == SQLITE_ROW)
    {
        // at least 1 row
        eof = false;
        return nRet;
    }
    nRet = sqlite3_finalize(m_pVM);

    if (nRet == SQLITE_SCHEMA)
    {
        return nRet;
    }

    qCritical() << sqlite3_errmsg(m_db);
    return -1;
}

int Sqlite3Statement::nextRow(bool &eof)
{
    int nRet = sqlite3_step(m_pVM);

    if (nRet == SQLITE_DONE)
    {
        // no rows
        sqlite3_finalize(m_pVM);
        eof = true;
    }
    else if (nRet == SQLITE_ROW)
    {
        // more rows, nothing to do
        eof = false;
    }
    else
    {
        nRet = sqlite3_finalize(m_pVM);

        qCritical() << sqlite3_errmsg(m_db);
        return nRet;
    }
    return nRet;
}

int Sqlite3Statement::endQuery()
{
    int nRet = sqlite3_finalize(m_pVM);
    if (nRet != SQLITE_OK && nRet != SQLITE_DONE)
    {
        qCritical() << sqlite3_errmsg(m_db);
    }
    return nRet;
}

int Sqlite3Statement::getInt(int column)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return -1;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return -1;
    }
    return sqlite3_column_int(m_pVM, column);
}

double Sqlite3Statement::getDouble(int column)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return -1;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return -1;
    }
    return sqlite3_column_double(m_pVM, column);
}

std::int64_t Sqlite3Statement::getInt64(int column)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return -1;
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return -1;
    }
    return sqlite3_column_int64(m_pVM, column);
}

std::string Sqlite3Statement::getString(int column)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return {};
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return {};
    }
    return {(const char *)sqlite3_column_text(m_pVM, column), static_cast<std::string::size_type>(sqlite3_column_bytes(m_pVM, column))};
}

QString Sqlite3Statement::getQString(int column)
{
    return QString::fromStdString(getString(column));
}

std::string Sqlite3Statement::getLastErrorString()
{
    return sqlite3_errmsg(m_db);
}

int Sqlite3Statement::countRow()
{
    int  count = 0;
    int  nRet  = 0;
    bool eof   = false;
    do
    {
        nRet = execQuery(eof);
        if (nRet == SQLITE_DONE || nRet == SQLITE_ROW)
        {
            while (!eof)
            {
                count = sqlite3_column_int(m_pVM, 0);
                nextRow(eof);
            }
            break;
        }
    } while (nRet == SQLITE_SCHEMA);
    return count;
}

bool Sqlite3Statement::isValid()
{
    return m_pVM != nullptr;
}

bool Sqlite3Statement::isDatabaseOpened()
{
    return m_db != nullptr;
}
QByteArray Sqlite3Statement::getBlob(int column)
{
    if (!isDatabaseOpened())
    {
        qWarning() << "database is not opened";
        return {};
    }
    if (!m_pVM)
    {
        qCritical() << "VM null pointer";
        return {};
    }

    return {static_cast<const char *>(sqlite3_column_blob(m_pVM, column)), sqlite3_column_bytes(m_pVM, column)};
}
void Sqlite3Statement::bind(int nParam, const QByteArray &blobValue)
{
    bind(nParam, (const unsigned char *)blobValue.data(), blobValue.length());
}
void Sqlite3Statement::bind(const char *szParam, const QByteArray &blobValue)
{
    bind(szParam, (const unsigned char *)blobValue.data(), blobValue.length());
}
void Sqlite3Statement::bind(const std::string &szParam, const QByteArray &blobValue)
{
    bind(szParam, (const unsigned char *)blobValue.data(), blobValue.length());
}
void Sqlite3Statement::bind(const QString &szParam, const QByteArray &blobValue)
{
    bind(szParam, (const unsigned char *)blobValue.data(), blobValue.length());
}
