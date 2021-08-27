#include <QFileOpenEvent>

#include "application.h"

#include "qtlocalpeer.h"

Application::Application(int &argc, char **argv) : QApplication(argc, argv)
{
    peer = new QtLocalPeer(this);
    connect(peer, &QtLocalPeer::messageReceived, this, &Application::messageReceived);
}

bool Application::event(QEvent *event)
{
    if (event->type() == QEvent::FileOpen)
    {
        QFileOpenEvent *openEvent = static_cast<QFileOpenEvent *>(event);
        emit            openUrl(openEvent->url());
    }

    return QApplication::event(event);
}

bool Application::isRunning()
{
    return peer->isClient();
}

bool Application::sendMessage(const QString &message, int timeout)
{
    return peer->sendMessage(message, timeout);
}

QString Application::id() const
{
    return peer->applicationId();
}
