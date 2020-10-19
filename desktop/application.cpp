#include <QFileOpenEvent>

#include "application.h"

bool Application::event(QEvent *event)
{
    if (event->type() == QEvent::FileOpen)
    {
        QFileOpenEvent *openEvent = static_cast<QFileOpenEvent *>(event);
        emit            openUrl(openEvent->url());
    }

    return QApplication::event(event);
}
