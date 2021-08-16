#include <QQmlApplicationEngine>

#include "qmlplayer.h"

static QQmlApplicationEngine *gQmlApplicationEngine = nullptr;

QmlPlayer::QmlPlayer(QObject *parent) : QObject(parent)
{
    
}

void QmlPlayer::Show()
{
    if (!InitQmlApplicationEngine())
    {
    }
}

bool QmlPlayer::InitQmlApplicationEngine()
{
    if (!gQmlApplicationEngine)
    {
        gQmlApplicationEngine = new QQmlApplicationEngine;
        gQmlApplicationEngine->load(QUrl("qrc:/rc/qml/musicplayer.qml"));
        return true;
    }

    return false;
}
