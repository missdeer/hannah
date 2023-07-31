#include <QProcess>
#include <QSettings>
#if defined(Q_OS_WIN)
#    include <Windows.h>
#else
#    include <signal.h>
#endif

#include "externalreverseproxyrunner.h"

ExternalReverseProxyRunner::ExternalReverseProxyRunner(QObject *parent) : m_process(new QProcess) {}

ExternalReverseProxyRunner::~ExternalReverseProxyRunner()
{
    m_process->terminate();
    delete m_process;
}

void ExternalReverseProxyRunner::start()
{
    m_process->start();
}

void ExternalReverseProxyRunner::stop()
{
#if defined(Q_OS_WIN)
    // Send Ctrl+C
    GenerateConsoleCtrlEvent(CTRL_C_EVENT, m_process->processId());
#else
    // Send SIGINT
    kill(m_process->processId(), SIGINT);
#endif
}

void ExternalReverseProxyRunner::restart()
{
    stop();
    start();
}

void ExternalReverseProxyRunner::applySettings(QSettings &settings)
{
    m_process->setProgram(settings.value("externalReverseProxyPath").toString());
    QStringList args;
    bool        ok    = false;
    int         state = settings.value("reverseProxyAutoRedirect", 2).toInt(&ok);
    if (ok && state == Qt::Checked)
    {
        args << "--auto-redirect";
    }
    state = settings.value("reverseProxyRedirect", 2).toInt(&ok);
    if (ok)
    {
        args << "--redirect";
    }
#if !defined(Q_OS_WIN)
    auto networkInterface = settings.value("reverseProxyBindNetworkInterface").toString();
    if (!networkInterface.isEmpty())
    {
        args << "-i" << networkInterface;
    }
#endif
    auto port = settings.value("reverseProxyListenPort", 8090).toInt(&ok);
    if (ok)
    {
        args << "-b" << QStringLiteral("localhost:%1").arg(port);
    }
    auto proxyType = settings.value("reverseProxyProxyType").toString();
    if (proxyType == QStringLiteral("Http"))
    {
        auto proxyAddr = settings.value("reverseProxyProxyAddress").toString();
        args << "-t" << proxyAddr;
    }
    if (proxyType == QStringLiteral("Socks5"))
    {
        auto proxyAddr = settings.value("reverseProxyProxyAddress").toString();
        args << "-s" << proxyAddr;
    }

    m_process->setArguments(args);
}

bool ExternalReverseProxyRunner::isRunning()
{
    return m_process->state() == QProcess::Running;
}
