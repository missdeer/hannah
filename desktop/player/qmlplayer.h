#ifndef QMLPLAYER_H
#define QMLPLAYER_H

#include <QObject>

class Player;

class QmlPlayer : public QObject
{
    Q_OBJECT
    Q_PROPERTY(qreal eq0 READ getEq0 WRITE setEq0 NOTIFY eq0Changed)
    Q_PROPERTY(qreal eq1 READ getEq1 WRITE setEq1 NOTIFY eq1Changed)
    Q_PROPERTY(qreal eq2 READ getEq2 WRITE setEq2 NOTIFY eq2Changed)
    Q_PROPERTY(qreal eq3 READ getEq3 WRITE setEq3 NOTIFY eq3Changed)
    Q_PROPERTY(qreal eq4 READ getEq4 WRITE setEq4 NOTIFY eq4Changed)
    Q_PROPERTY(qreal eq5 READ getEq5 WRITE setEq5 NOTIFY eq5Changed)
    Q_PROPERTY(qreal eq6 READ getEq6 WRITE setEq6 NOTIFY eq6Changed)
    Q_PROPERTY(qreal eq7 READ getEq7 WRITE setEq7 NOTIFY eq7Changed)
    Q_PROPERTY(qreal eq8 READ getEq8 WRITE setEq8 NOTIFY eq8Changed)
    Q_PROPERTY(qreal eq9 READ getEq9 WRITE setEq9 NOTIFY eq9Changed)
    Q_PROPERTY(qreal volumn READ getVolumn WRITE setVolumn NOTIFY volumnChanged)
    Q_PROPERTY(qreal progress READ getProgress WRITE setProgress NOTIFY progressChanged)
    Q_PROPERTY(QString coverUrl READ getCoverUrl WRITE setCoverUrl NOTIFY coverUrlChanged)
    Q_PROPERTY(QString songName READ getSongName WRITE setSongName NOTIFY songNameChanged)
public:
    explicit QmlPlayer(QObject *parent = nullptr);

    Q_INVOKABLE void onQuit();
    Q_INVOKABLE void onShowPlaylists();
    Q_INVOKABLE void onSettings();
    Q_INVOKABLE void onFilter();
    Q_INVOKABLE void onMessage();
    Q_INVOKABLE void onMusic();
    Q_INVOKABLE void onCloud();
    Q_INVOKABLE void onBluetooth();
    Q_INVOKABLE void onCart();
    Q_INVOKABLE void presetEQChanged(int index);
    Q_INVOKABLE void onOpenPreset();
    Q_INVOKABLE void onSavePreset();
    Q_INVOKABLE void onFavorite();
    Q_INVOKABLE void onStop();
    Q_INVOKABLE void onPrevious();
    Q_INVOKABLE void onPause();
    Q_INVOKABLE void onNext();
    Q_INVOKABLE void onRepeat();
    Q_INVOKABLE void onShuffle();
    Q_INVOKABLE void onSwitchFiles();
    Q_INVOKABLE void onSwitchPlaylists();
    Q_INVOKABLE void onSwitchFavourites();
    Q_INVOKABLE void onOpenFile();

    void Show();

    qreal          getEq0() const;
    qreal          getEq1() const;
    qreal          getEq2() const;
    qreal          getEq3() const;
    qreal          getEq4() const;
    qreal          getEq5() const;
    qreal          getEq6() const;
    qreal          getEq7() const;
    qreal          getEq8() const;
    qreal          getEq9() const;
    qreal          getVolumn() const;
    qreal          getProgress() const;
    const QString &getCoverUrl() const;
    const QString &getSongName() const;

    void setEq0(qreal value);
    void setEq1(qreal value);
    void setEq2(qreal value);
    void setEq3(qreal value);
    void setEq4(qreal value);
    void setEq5(qreal value);
    void setEq6(qreal value);
    void setEq7(qreal value);
    void setEq8(qreal value);
    void setEq9(qreal value);
    void setVolumn(qreal value);
    void setProgress(qreal progress);
    void setCoverUrl(const QString &u);
    void setSongName(const QString &n);
signals:
    void showPlayer();
    void eq0Changed();
    void eq1Changed();
    void eq2Changed();
    void eq3Changed();
    void eq4Changed();
    void eq5Changed();
    void eq6Changed();
    void eq7Changed();
    void eq8Changed();
    void eq9Changed();
    void volumnChanged();
    void progressChanged();
    void coverUrlChanged();
    void songNameChanged();

private:
    Player *m_player {nullptr};
    qreal   m_eq0 {0.0};
    qreal   m_eq1 {0.0};
    qreal   m_eq2 {0.0};
    qreal   m_eq3 {0.0};
    qreal   m_eq4 {0.0};
    qreal   m_eq5 {0.0};
    qreal   m_eq6 {0.0};
    qreal   m_eq7 {0.0};
    qreal   m_eq8 {0.0};
    qreal   m_eq9 {0.0};
    qreal   m_volumn {0.0};
    qreal   m_progress {0.0};
    QString m_coverUrl;
    QString m_songName;

    void applyEQ();
};

inline QmlPlayer *qmlPlayer = nullptr;

#endif // QMLPLAYER_H
