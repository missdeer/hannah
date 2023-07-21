#pragma once
#include <cstdint>
#include <memory>
#include <string>
#include <string_view>

#include <QByteArray>
#include <QString>

struct sqlite3;
struct sqlite3_stmt;

class Sqlite3Helper;

class Sqlite3Statement
{
    friend class Sqlite3Helper;

public:
    Sqlite3Statement(sqlite3 *&database, sqlite3_stmt *pVM);

    void bind(int nParam, const QString &value);
    void bind(int nParam, const char *szValue, int nLen = -1);
    void bind(int nParam, const std::string &szValue);
    void bind(int nParam, std::string_view szValue);
    void bind(int nParam, int nValue);
    void bind(int nParam, int64_t nValue);
    void bind(int nParam, double dValue);
    void bind(int nParam, const unsigned char *blobValue, int nLen);
    void bind(int nParam, const QByteArray &blobValue);
    void bindNull(int nParam);
    int  bindParameterIndex(const char *szParam);
    void bind(const char *szParam, const QString &value);
    void bind(const char *szParam, const char *szValue, int nLen = -1);
    void bind(const char *szParam, const std::string &szValue);
    void bind(const char *szParam, std::string_view szValue);
    void bind(const char *szParam, int nValue);
    void bind(const char *szParam, int64_t nValue);
    void bind(const char *szParam, double dwValue);
    void bind(const char *szParam, const unsigned char *blobValue, int nLen);
    void bind(const char *szParam, const QByteArray &blobValue);
    void bindNull(const char *szParam);
    int  bindParameterIndex(const std::string &szParam);
    void bind(const std::string &szParam, const QString &value);
    void bind(const std::string &szParam, const std::string &szValue);
    void bind(const std::string &szParam, std::string_view szValue);
    void bind(const std::string &szParam, int nValue);
    void bind(const std::string &szParam, int64_t nValue);
    void bind(const std::string &szParam, double dwValue);
    void bind(const std::string &szParam, const unsigned char *blobValue, int nLen);
    void bind(const std::string &szParam, const QByteArray &blobValue);
    void bindNull(const std::string &szParam);
    int  bindParameterIndex(const QString &szParam);
    void bind(const QString &szParam, const QString &value);
    void bind(const QString &szParam, const std::string &szValue);
    void bind(const QString &szParam, std::string_view szValue);
    void bind(const QString &szParam, int nValue);
    void bind(const QString &szParam, int64_t nValue);
    void bind(const QString &szParam, double dwValue);
    void bind(const QString &szParam, const unsigned char *blobValue, int nLen);
    void bind(const QString &szParam, const QByteArray &blobValue);
    void bindNull(const QString &szParam);

    int execDML();
    int execQuery(bool &eof);
    int nextRow(bool &eof);
    int endQuery();

    int          getInt(int column);
    double       getDouble(int column);
    std::int64_t getInt64(int column);
    std::string  getString(int column);
    QString      getQString(int column);
    QByteArray   getBlob(int column);

    std::string getLastErrorString();

    int countRow();

    bool isValid();

private:
    sqlite3     *&m_db;
    sqlite3_stmt *m_pVM;

    bool isDatabaseOpened();
};

using Sqlite3StatementPtr = std::shared_ptr<Sqlite3Statement>;
