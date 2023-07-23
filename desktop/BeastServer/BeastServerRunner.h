#pragma once

#include <QObject>
#include <QThread>

class BeastServerRunner : public QThread
{
    Q_OBJECT

public:
    void stop();
    void setPort(unsigned short port);
    void setHttpProxy(const QString &proxy);
    void setSocks5Proxy(const QString &proxy);
    void setNetworkInterface(const QString &interface);
    void setAutoRedirect(bool checked);
    void setRedirect(bool checked);
    void loadConfigurations();

protected:
    void run() override;
};
