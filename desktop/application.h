#ifndef APPLICATION_H
#define APPLICATION_H

#include <QApplication>
#include <QUrl>

QT_BEGIN_NAMESPACE
class QEvent;
class QtLocalPeer;
QT_END_NAMESPACE

class Application : public QApplication
{
    Q_OBJECT

public:
    Application(int &argc, char **argv);

    bool event(QEvent *event) override;

    bool    isRunning();
    QString id() const;

public slots:
    bool sendMessage(const QString &message, int timeout = 5000);

signals:
    void openUrl(QUrl);

    void messageReceived(const QString &message);

private:
    QtLocalPeer *peer;
};

#endif // APPLICATION_H
