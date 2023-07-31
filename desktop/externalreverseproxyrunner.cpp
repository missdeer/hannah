#include <QProcess>
#include <QSettings>
#if defined(Q_OS_WIN)
#    include <Windows.h>
#else
#    include <signal.h>
#endif

#include "externalreverseproxyrunner.h"

ExternalReverseProxyRunner::ExternalReverseProxyRunner(QObject *parent) {}

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

void ExternalReverseProxyRunner::applySettings(QSettings &settings) {}

bool ExternalReverseProxyRunner::isRunning()
{
    return m_process->state() == QProcess::Running;
}
