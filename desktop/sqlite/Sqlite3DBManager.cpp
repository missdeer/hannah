#include <filesystem>
#include <random>
#include <sqlite3.h>

#include <QtCore>

#include "Sqlite3DBManager.h"
#include "Sqlite3Constants.h"

void Sqlite3DBManager::regenerateSavePoint()
{
    static const std::string        chars("abcdefghijklmnopqrstuvwxyz"
                                          "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
                                          "1234567890");
    std::random_device              rng;
    std::uniform_int_distribution<> index_dist(0, static_cast<int>(chars.size()) - 1);
    std::string                     res;
    const int                       length = 16;
    for (int i = 0; i < length; i++)
    {
        res.append(1, chars[index_dist(rng)]);
    }
    m_savePoint = res;
}

void Sqlite3DBManager::clearSavePoint()
{
    m_savePoint.clear();
}

bool Sqlite3DBManager::setSavePoint()
{
    std::string sql = "SAVEPOINT '" + m_savePoint + "';";
    if (m_sqlite.execDML(sql) != SQLITE_OK)
    {
        qCritical() << "Cannot set savepoint.";
        return false;
    }
    qDebug() << "set savepoint" << QString::fromStdString(m_savePoint);
    return true;
}

bool Sqlite3DBManager::releaseSavePoint()
{
    std::string sql = "RELEASE '" + m_savePoint + "';";
    if (m_sqlite.execDML(sql) != SQLITE_OK)
    {
        qCritical() << "Cannot release savepoint.";
        return false;
    }
    qDebug() << "released savepoint" << QString::fromStdString(m_savePoint);
    return true;
}

bool Sqlite3DBManager::open(const QString &dbPath, bool readOnly)
{
    return open(dbPath.toUtf8().constData(), readOnly);
}

bool Sqlite3DBManager::open(const std::string &dbPath, bool readOnly)
{
    return open(dbPath.c_str(), readOnly);
}

bool Sqlite3DBManager::open(const char *dbPath, bool readOnly)
{
    if (sqlite3_open_v2(dbPath, &m_db, (readOnly ? SQLITE_OPEN_READONLY : SQLITE_OPEN_READWRITE), nullptr) != SQLITE_OK)
    {
        qCritical() << "Cannot open the database" << QString::fromUtf8(dbPath);
        return false;
    }

    Q_ASSERT(m_db);

    if (!readOnly)
    {
        regenerateSavePoint();
        setSavePoint();
    }

    qDebug() << "The database" << QString::fromStdString(dbPath) << "is opened";
    return true;
}

bool Sqlite3DBManager::close()
{
    if (!m_db)
    {
        qWarning() << "database is not opened";
        return true;
    }
    int result = sqlite3_close_v2(m_db);
    while (result == SQLITE_BUSY)
    {
        qDebug() << "Close db result: " << result << sqlite3_errmsg(m_db);
        sqlite3_stmt *stmt = sqlite3_next_stmt(m_db, nullptr);
        if (stmt)
        {
            result = sqlite3_finalize(stmt);
            if (result == SQLITE_OK)
            {
                result = sqlite3_close_v2(m_db);
                qDebug() << "Secondary try closing db result: " << result << sqlite3_errmsg(m_db);
            }
        }
    }

    if (result != SQLITE_OK)
    {
        qCritical() << "Closing db failed:" << result << sqlite3_errmsg(m_db);
        return false;
    }

    m_db = nullptr;
    qDebug() << "database closed";
    return true;
}

bool Sqlite3DBManager::save()
{
    releaseSavePoint();
    // always set a new save point
    regenerateSavePoint();
    return setSavePoint();
}

bool Sqlite3DBManager::saveAndClose()
{
    if (!releaseSavePoint())
    {
        return false;
    }

    return close();
}

bool Sqlite3DBManager::create(const QString &dbPath)
{
    return create(dbPath.toUtf8().constData());
}

bool Sqlite3DBManager::create(const std::string &dbPath)
{
    return create(dbPath.c_str());
}

bool Sqlite3DBManager::create(const char *dbPath)
{
    if (sqlite3_open_v2(dbPath, &m_db, SQLITE_OPEN_CREATE | SQLITE_OPEN_READWRITE, nullptr) != SQLITE_OK)
    {
        qCritical() << "Cannot create the database" << QString::fromUtf8(dbPath);
        return false;
    }

    Q_ASSERT(m_db);

    qInfo() << "The database is created at" << QString::fromStdString(dbPath);

    qDebug() << "set sqlite options for performance issue";

    if (m_sqlite.execDML(SQL_STATEMENT_CACHE_SIZE) != SQLITE_OK)
    {
        qCritical() << "Cannot set cache size.";
        return false;
    }

    if (m_sqlite.execDML(SQL_STATEMENT_ENABLE_FOREIGN_KEYS) != SQLITE_OK)
    {
        qCritical() << "Cannot enable the foreign keys.";
        return false;
    }

    regenerateSavePoint();
    setSavePoint();

    return true;
}

Sqlite3Helper &Sqlite3DBManager::engine()
{
    return m_sqlite;
}

Sqlite3DBManager::~Sqlite3DBManager()
{
    if (isOpened())
    {
        close();
    }
}

bool Sqlite3DBManager::isOpened() const
{
    return m_db != nullptr;
}

bool Sqlite3DBManager::loadOrSaveInMemory(const QString &dbPath, bool isSave)
{
    sqlite3 *pFile = nullptr; /* Database connection opened on zFilename */

    /* Open the database file identified by zFilename. Exit early if this fails
    ** for any reason. */
#if defined(Q_OS_WIN)
    int result = sqlite3_open(dbPath.toUtf8().constData(), &pFile);
#else
    int result = sqlite3_open(dbPath.toStdString().c_str(), &pFile);
#endif
    if (result == SQLITE_OK)
    {
        /* If this is a 'load' operation (isSave==0), then data is copied
        ** from the database file just opened to database pInMemory.
        ** Otherwise, if this is a 'save' operation (isSave==1), then data
        ** is copied from pInMemory to pFile.  Set the variables pFrom and
        ** pTo accordingly. */
        auto *pFrom = (isSave ? m_db : pFile);
        auto *pTo   = (isSave ? pFile : m_db);

        /* Set up the backup procedure to copy from the "main" database of
        ** connection pFile to the main database of connection pInMemory.
        ** If something goes wrong, pBackup will be set to NULL and an error
        ** code and message left in connection pTo.
        **
        ** If the backup object is successfully created, call backup_step()
        ** to copy data from pFile to pInMemory. Then call backup_finish()
        ** to release resources associated with the pBackup object.  If an
        ** error occurred, then an error code and message will be left in
        ** connection pTo. If no error occurred, then the error code belonging
        ** to pTo is set to SQLITE_OK.
        */
        auto *pBackup = sqlite3_backup_init(pTo, "main", pFrom, "main");
        if (pBackup)
        {
            (void)sqlite3_backup_step(pBackup, -1);
            (void)sqlite3_backup_finish(pBackup);
        }
        result = sqlite3_errcode(pTo);
    }
    else
    {
        qCritical() << "cannot open file" << dbPath;
    }

    /* Close the database connection opened on database file zFilename
    ** and return the result of this function. */
    (void)sqlite3_close(pFile);
    return result == SQLITE_OK;
}

bool Sqlite3DBManager::saveAs(const QString &newDbPath)
{
    return saveAs(newDbPath.toUtf8().constData());
}

bool Sqlite3DBManager::saveAs(const std::string &newDbPath)
{
    return saveAs(newDbPath.c_str());
}

bool Sqlite3DBManager::saveAs(const char *newDbPath)
{
    releaseSavePoint();

    sqlite3 *newDb = nullptr;
    if (sqlite3_open_v2(newDbPath, &newDb, SQLITE_OPEN_CREATE | SQLITE_OPEN_READWRITE, nullptr) != SQLITE_OK)
    {
        qCritical() << "Cannot open the database" << QString::fromUtf8(newDbPath);
        return false;
    }

    Q_ASSERT(newDbPath);
    Sqlite3Helper newEngine(newDb);

    newEngine.beginTransaction();

    sqlite3_backup *backup = sqlite3_backup_init(newDb, "main", m_db, "main");
    if (backup == nullptr)
    {
        qCritical() << "sqlite3_backup_init failed:" << sqlite3_errmsg(newDb);
        return false;
    }

    int resultCode = sqlite3_backup_step(backup, -1);
    if (resultCode != SQLITE_DONE)
    {
        qCritical() << "sqlite3_backup_step failed:" << resultCode << "-" << sqlite3_errmsg(newDb);
        return false;
    }

    resultCode = sqlite3_backup_finish(backup);
    if (resultCode != SQLITE_OK)
    {
        qCritical() << "sqlite3_backup_finish failed:" << resultCode << "-" << sqlite3_errmsg(newDb);
        return false;
    }

    newEngine.endTransaction();

    close();

    m_db = newDb;

    regenerateSavePoint();
    setSavePoint();

    return true;
}
