#include <QSettings>

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

void BeastServerRunner::applySettings(QSettings &settings) {}
