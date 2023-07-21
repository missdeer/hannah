#ifndef NETWORKREPLYHELPER_H
#define NETWORKREPLYHELPER_H

#include <QAbstractItemModel>
#include <QNetworkReply>

QT_BEGIN_NAMESPACE
class QTimer;
class QIODevice;
QT_END_NAMESPACE

class NetworkReplyHelper : public QObject
{
    Q_OBJECT
public:
    explicit NetworkReplyHelper(QNetworkReply *reply, QIODevice *storage = nullptr, QObject *parent = nullptr);
    explicit NetworkReplyHelper(QNetworkReply *reply, QByteArray *content, QObject *parent = nullptr);
    NetworkReplyHelper(const NetworkReplyHelper &) = delete;
    void operator=(const NetworkReplyHelper &)     = delete;
    NetworkReplyHelper(NetworkReplyHelper &&)      = delete;
    void operator=(NetworkReplyHelper &&)          = delete;
    virtual ~NetworkReplyHelper();

    virtual void       postFinished();
    virtual QByteArray receivedData(QByteArray data);

    [[nodiscard]] QJsonDocument  json();
    [[nodiscard]] QByteArray    &content();
    [[nodiscard]] QNetworkReply *reply();
    [[nodiscard]] QIODevice     *storage();

    [[nodiscard]] QVariant       data() const;
    void                         setData(const QVariant &data);
    void                         setTimeout(int milliseconds);
    [[nodiscard]] const QString &getErrorMessage() const;

    void               waitForFinished();
    [[nodiscard]] bool isOk() const;

protected:
    void setErrorMessage(const QString &errMsg);

signals:
    void done();
    void cancel();
    void errorMessage(QNetworkReply::NetworkError, QString);
public slots:
    void downloadProgress(qint64 bytesReceived, qint64 bytesTotal);
    void error(QNetworkReply::NetworkError code);
    void finished();
    void sslErrors(const QList<QSslError> &errors);
    void uploadProgress(qint64 bytesSent, qint64 bytesTotal);
    void readyRead();
private slots:
    void timeout();

private:
    void readData(QNetworkReply *reply);

    QNetworkReply              *m_reply {nullptr};
    QIODevice                  *m_storage {nullptr};
    QByteArray                 *m_content {nullptr};
    QTimer                     *m_timeoutTimer {nullptr};
    QNetworkReply::NetworkError m_error {QNetworkReply::NoError};
    QVariant                    m_data;
    QString                     m_errMsg;
    bool                        m_ownContentObject {false};
};

#endif // NETWORKREPLYHELPER_H
