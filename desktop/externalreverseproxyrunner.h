#pragma once

#include <QObject>

QT_FORWARD_DECLARE_CLASS(QProcess);
QT_FORWARD_DECLARE_CLASS(QSettings);

class ExternalReverseProxyRunner : public QObject
{
    Q_OBJECT
public:
    explicit ExternalReverseProxyRunner(QObject *parent = nullptr);

    void start();
    void stop();
    void restart();
    void applySettings(QSettings& settings);

    [[nodiscard]] bool isRunning();

private:
    QProcess *m_process {nullptr};
};