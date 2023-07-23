#include "BeastServerRunner.h"
#include "BeastServer.h"

void BeastServerRunner::stop()
{
    StopBeastServer();
}

void BeastServerRunner::run()
{
    StartBeastServer();
}

void BeastServerRunner::setPort(unsigned short port)
{
    SetListenPort(port);
}

void BeastServerRunner::setHttpProxy(const QString &proxy) {}

void BeastServerRunner::setSocks5Proxy(const QString &proxy) {}

void BeastServerRunner::setNetworkInterface(const QString &interface) {}

void BeastServerRunner::setAutoRedirect(bool checked) {}

void BeastServerRunner::setRedirect(bool checked) {}

void BeastServerRunner::loadConfigurations() {}
