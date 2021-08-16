#ifndef QMLPLAYER_H
#define QMLPLAYER_H

#include <QObject>

QT_FORWARD_DECLARE_CLASS(QQmlApplicationEngine);

class QmlPlayer : public QObject
{
    Q_OBJECT
public:
    explicit QmlPlayer(QObject *parent = nullptr);

    void Show();

signals:

private:
    bool InitQmlApplicationEngine();
};

inline QmlPlayer *qmlPlayer = nullptr;

#endif // QMLPLAYER_H
