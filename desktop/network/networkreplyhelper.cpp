#include <zlib.h>

#include <QEventLoop>
#include <QIODevice>
#include <QJsonDocument>
#include <QJsonParseError>
#include <QTimer>

#include "networkreplyhelper.h"

static QByteArray gUncompress(const QByteArray &data)
{
    if (data.size() <= 4)
    {
        qWarning("gUncompress: Input data is truncated");
        return data;
    }

    QByteArray result;

    int              ret = 0;
    z_stream         strm;
    static const int CHUNK_SIZE = 1024;
    char             out[CHUNK_SIZE];

    /* allocate inflate state */
    strm.zalloc   = Z_NULL;
    strm.zfree    = Z_NULL;
    strm.opaque   = Z_NULL;
    strm.avail_in = data.size();
    strm.next_in  = (Bytef *)(data.data());

    ret = inflateInit2(&strm, 15 + 32); // gzip decoding
    if (ret != Z_OK)
    {
        return data;
    }

    // run inflate()
    do
    {
        strm.avail_out = CHUNK_SIZE;
        strm.next_out  = (Bytef *)(out);

        ret = inflate(&strm, Z_NO_FLUSH);
        Q_ASSERT(ret != Z_STREAM_ERROR); // state not clobbered

        switch (ret)
        {
        case Z_NEED_DICT:
            ret = Z_DATA_ERROR; // and fall through
        case Z_DATA_ERROR:
        case Z_MEM_ERROR:
            (void)inflateEnd(&strm);
            return data;
        }

        result.append(out, CHUNK_SIZE - strm.avail_out);
    } while (strm.avail_out == 0);

    // clean up and return
    inflateEnd(&strm);
    return result;
}

NetworkReplyHelper::NetworkReplyHelper(QNetworkReply *reply, QByteArray *content, QObject *parent)
    : QObject(parent), m_reply(reply), m_content(content)
{
    Q_ASSERT(reply);
    connect(reply, &QNetworkReply::downloadProgress, this, &NetworkReplyHelper::downloadProgress);
#if QT_VERSION >= QT_VERSION_CHECK(5, 15, 0)
    connect(reply, &QNetworkReply::errorOccurred, this, &NetworkReplyHelper::error);
#else
    connect(reply, static_cast<void (QNetworkReply::*)(QNetworkReply::NetworkError)>(&QNetworkReply::error), this, &NetworkReplyHelper::error);
#endif
    connect(reply, &QNetworkReply::finished, this, &NetworkReplyHelper::finished);
    connect(reply, &QNetworkReply::sslErrors, this, &NetworkReplyHelper::sslErrors);
    connect(reply, &QNetworkReply::uploadProgress, this, &NetworkReplyHelper::uploadProgress);
    connect(reply, &QNetworkReply::readyRead, this, &NetworkReplyHelper::readyRead);
    if (!m_content)
    {
        m_content          = new QByteArray;
        m_ownContentObject = true;
    }
}

NetworkReplyHelper::NetworkReplyHelper(QNetworkReply *reply, QIODevice *storage, QObject *parent)
    : QObject(parent), m_reply(reply), m_storage(storage)
{
    Q_ASSERT(reply);
    connect(reply, &QNetworkReply::downloadProgress, this, &NetworkReplyHelper::downloadProgress);
#if QT_VERSION >= QT_VERSION_CHECK(5, 15, 0)
    connect(reply, &QNetworkReply::errorOccurred, this, &NetworkReplyHelper::error);
#else
    connect(reply, static_cast<void (QNetworkReply::*)(QNetworkReply::NetworkError)>(&QNetworkReply::error), this, &NetworkReplyHelper::error);
#endif
    connect(reply, &QNetworkReply::finished, this, &NetworkReplyHelper::finished);
    connect(reply, &QNetworkReply::sslErrors, this, &NetworkReplyHelper::sslErrors);
    connect(reply, &QNetworkReply::uploadProgress, this, &NetworkReplyHelper::uploadProgress);
    connect(reply, &QNetworkReply::readyRead, this, &NetworkReplyHelper::readyRead);
    if (!m_storage)
    {
        m_content          = new QByteArray;
        m_ownContentObject = true;
    }
}

NetworkReplyHelper::~NetworkReplyHelper()
{
    if (m_reply)
    {
        m_reply->disconnect(this);
        m_reply->deleteLater();
        m_reply = nullptr;
    }

    if (m_timeoutTimer)
    {
        if (m_timeoutTimer->isActive())
        {
            m_timeoutTimer->stop();
        }
        delete m_timeoutTimer;
        m_timeoutTimer = nullptr;
    }

    if (m_ownContentObject)
    {
        delete m_content;
        m_content = nullptr;
    }
}

QJsonDocument NetworkReplyHelper::json()
{
    QJsonParseError err {};
    QJsonDocument   jsonDocument = QJsonDocument::fromJson(*m_content, &err);
    if (err.error == QJsonParseError::NoError)
    {
        return jsonDocument;
    }
    return {};
}

void NetworkReplyHelper::waitForFinished()
{
    QEventLoop loop;
    QObject::connect(m_reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    QObject::disconnect(m_reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
}

void NetworkReplyHelper::downloadProgress(qint64 bytesReceived, qint64 bytesTotal)
{
#if defined(NETWORK_LOG)
    qDebug() << __FUNCTION__ << __LINE__ << bytesReceived << bytesTotal;
#else
    Q_UNUSED(bytesReceived);
    Q_UNUSED(bytesTotal);
#endif
    if (m_timeoutTimer && m_timeoutTimer->isActive())
    {
        m_timeoutTimer->start();
    }
}

void NetworkReplyHelper::error(QNetworkReply::NetworkError code)
{
    m_error     = code;
    auto *reply = qobject_cast<QNetworkReply *>(sender());
    m_errMsg.append(reply->errorString() + "\n");

#if defined(NETWORK_LOG)
    qCritical() << __FUNCTION__ << __LINE__ << m_errMsg;
#endif

    emit errorMessage(code, m_errMsg);
}

void NetworkReplyHelper::finished()
{
    if (m_timeoutTimer)
    {
        m_timeoutTimer->stop();
    }

    auto *reply           = qobject_cast<QNetworkReply *>(sender());
    auto  contentEncoding = reply->rawHeader("Content-Encoding");

#if defined(NETWORK_LOG)
    auto headerList = reply->rawHeaderList();
    for (const auto &header : headerList)
    {
        qDebug() << __FUNCTION__ << __LINE__ << QString(header) << ":" << QString(reply->rawHeader(header));
    }
    qDebug() << __FUNCTION__ << __LINE__ << contentEncoding << m_content->length() << "bytes received";
#endif
    readData(reply);
    if (contentEncoding == "gzip" || contentEncoding == "deflate")
    {
        auto content = gUncompress(*m_content);
        if (m_storage)
        {
            m_storage->write(content);
        }
        if (!content.isEmpty())
        {
            m_content->swap(content);
        }
    }

#if defined(NETWORK_LOG)
    qDebug() << __FUNCTION__ << __LINE__ << m_content->length() << "bytes uncompressed: " << QString(*m_content).left(256) << "\n";
#endif
    postFinished();
}

void NetworkReplyHelper::sslErrors(const QList<QSslError> &errors)
{
    for (const auto &err : errors)
    {
#if defined(NETWORK_LOG)
        qCritical() << __FUNCTION__ << __LINE__ << err.errorString();
#endif
        m_errMsg.append(err.errorString() + "\n");
    }
}

void NetworkReplyHelper::uploadProgress(qint64 bytesSent, qint64 bytesTotal)
{
    Q_UNUSED(bytesSent);
    Q_UNUSED(bytesTotal);
    if (m_timeoutTimer && m_timeoutTimer->isActive())
    {
        m_timeoutTimer->start();
    }
}

void NetworkReplyHelper::readyRead()
{
    auto *reply = qobject_cast<QNetworkReply *>(sender());
#if defined(NETWORK_LOG)
    int statusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
    qDebug() << __FUNCTION__ << __LINE__ << statusCode << reply;
#endif
    // if (statusCode >= 200 && statusCode < 300)
    {
        readData(reply);
    }
}

void NetworkReplyHelper::timeout()
{
#if defined(NETWORK_LOG)
    qDebug() << __FUNCTION__ << __LINE__;
#endif
    if (m_reply && m_reply->isRunning())
    {
        m_reply->abort();
    }
}

void NetworkReplyHelper::readData(QNetworkReply *reply)
{
    auto resp            = receivedData(reply->readAll());
    auto contentEncoding = reply->rawHeader("Content-Encoding");
#if defined(NETWORK_LOG)
    qDebug() << __FUNCTION__ << __LINE__ << contentEncoding << resp.length() << QString(resp);
#endif
    if (contentEncoding != "gzip" && contentEncoding != "deflate")
    {
        if (m_storage)
        {
            m_storage->write(resp);
        }
    }
    m_content->append(resp);
}

QVariant NetworkReplyHelper::data() const
{
    return m_data;
}

void NetworkReplyHelper::setData(const QVariant &data)
{
    m_data = data;
}

void NetworkReplyHelper::setTimeout(int milliseconds)
{
    if (!m_reply->isRunning())
    {
        return;
    }
    if (!m_timeoutTimer)
    {
        m_timeoutTimer = new QTimer;
        m_timeoutTimer->setSingleShot(true);
        connect(m_timeoutTimer, &QTimer::timeout, this, &NetworkReplyHelper::timeout);
    }
    m_timeoutTimer->start(milliseconds);
}

const QString &NetworkReplyHelper::getErrorMessage() const
{
    return m_errMsg;
}

QNetworkReply *NetworkReplyHelper::reply()
{
    return m_reply;
}

QIODevice *NetworkReplyHelper::storage()
{
    return m_storage;
}

QByteArray &NetworkReplyHelper::content()
{
    return *m_content;
}

bool NetworkReplyHelper::isOk() const
{
    return m_error == QNetworkReply::NoError;
}

void NetworkReplyHelper::postFinished()
{
    emit done();
}

void NetworkReplyHelper::setErrorMessage(const QString &errMsg)
{
    m_errMsg = errMsg;
}

QByteArray NetworkReplyHelper::receivedData(QByteArray data)
{
    return data;
}
