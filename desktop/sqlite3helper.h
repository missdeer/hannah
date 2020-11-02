#ifndef SQLITE3HELPER_H
#define SQLITE3HELPER_H

#include <QSharedPointer>
#include <QString>
#include <sqlite3.h>

class Sqlite3Helper
{
public:
    Sqlite3Helper();
    ~Sqlite3Helper();
    void bind(sqlite3_stmt *pVM, int nParam, const QString &sValue);
    void bind(sqlite3_stmt *pVM, int nParam, const char *szValue);
    void bind(sqlite3_stmt *pVM, int nParam, const int nValue);
    void bind(sqlite3_stmt *pVM, int nParam, const int64_t nValue);
    void bind(sqlite3_stmt *pVM, int nParam, const double dValue);
    void bind(sqlite3_stmt *pVM, int nParam, const unsigned char *blobValue, int nLen);
    void bindNull(sqlite3_stmt *pVM, int nParam);
    int  bindParameterIndex(sqlite3_stmt *pVM, const char *szParam);
    void bind(sqlite3_stmt *pVM, const char *szParam, const QString &sValue);
    void bind(sqlite3_stmt *pVM, const char *szParam, const char *szValue);
    void bind(sqlite3_stmt *pVM, const char *szParam, const int nValue);
    void bind(sqlite3_stmt *pVM, const char *szParam, const int64_t nValue);
    void bind(sqlite3_stmt *pVM, const char *szParam, const double dwValue);
    void bind(sqlite3_stmt *pVM, const char *szParam, const unsigned char *blobValue, int nLen);
    void bindNull(sqlite3_stmt *pVM, const char *szParam);

    sqlite3_stmt *compile(const QString &szSQL);
    sqlite3_stmt *compile(const char *szSQL);
    int           execDML(const QString &szSQL);
    int           execDML(const char *szSQL);
    int           execDML(sqlite3_stmt *pVM);
    int           execQuery(sqlite3_stmt *pVM, bool &eof);
    bool          nextRow(sqlite3_stmt *pVM, bool &eof);

    bool isDatabaseOpened();
    bool closeDatabaseConnection();

    int  checkExists(const QString &field, const QString &name);
    bool openDatabase(const QString &name);
    bool createDatabase(const QString &name);

    bool beginTransaction();
    bool endTransaction();
    bool rollbackTransaction();

    bool vacuum();

private:
    sqlite3 *m_db {nullptr};
    QString  m_dbFile;
};

typedef QSharedPointer<Sqlite3Helper> Sqlite3HelperPtr;

#endif // SQLITE3HELPER_H
