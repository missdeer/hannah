#pragma once

#include <QObject>
#include <QThread>

QT_FORWARD_DECLARE_CLASS(QSettings);

class BeastServerRunner : public QThread
{
    Q_OBJECT

public:
    void stop();
    
    void applySettings(QSettings& settings);

protected:
    void run() override;
};
