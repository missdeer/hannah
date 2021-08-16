#include <QQmlApplicationEngine>
#include <QQmlContext>

#include "qmlplayer.h"
#include "player.h"

static QQmlApplicationEngine *gQmlApplicationEngine = nullptr;

QmlPlayer::QmlPlayer(QObject *parent) : QObject(parent), m_player(new Player) {}

void QmlPlayer::Show()
{
    if (!InitQmlApplicationEngine())
    {
        emit showPlayer();
    }
}

bool QmlPlayer::InitQmlApplicationEngine()
{
    if (!gQmlApplicationEngine)
    {
        gQmlApplicationEngine = new QQmlApplicationEngine;
        gQmlApplicationEngine->load(QUrl("qrc:/rc/qml/musicplayer.qml"));

        QQmlContext *context = gQmlApplicationEngine->rootContext();
        context->setContextProperty("playerCore", this);
        return true;
    }

    return false;
}
