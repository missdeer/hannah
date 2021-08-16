#ifndef QMLPLAYER_H
#define QMLPLAYER_H

#include <QObject>

QT_FORWARD_DECLARE_CLASS(QQmlApplicationEngine);

class Player;

class QmlPlayer : public QObject
{
    Q_OBJECT
public:
    explicit QmlPlayer(QObject *parent = nullptr);

    void Show();

signals:
    void showPlayer();

private:
    Player *m_player {nullptr};

    bool InitQmlApplicationEngine();
};

inline QmlPlayer *qmlPlayer = nullptr;

#endif // QMLPLAYER_H
