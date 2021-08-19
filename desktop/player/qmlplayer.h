#ifndef QMLPLAYER_H
#define QMLPLAYER_H

#include <QObject>
#if defined(Q_OS_WIN) && QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
#    include <QtWinExtras>
#endif

QT_FORWARD_DECLARE_CLASS(QTimer);
QT_FORWARD_DECLARE_CLASS(QQmlApplicationEngine);

class PlayList;
class Lyrics;
class OSD;
class LrcBar;

class QmlPlayer : public QObject
{
    Q_OBJECT
    Q_PROPERTY(int primaryScreenWidth READ getPrimaryScreenWidth NOTIFY primaryScreenWidthChanged)
    Q_PROPERTY(int primaryScreenHeight READ getPrimaryScreenHeight NOTIFY primaryScreenHeightChanged)
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
    ~QmlPlayer();
    void showNormal();
    void loadAudio(const QString &uri);
    void addToListAndPlay(const QList<QUrl> &uris);
    void addToListAndPlay(const QStringList &uris);
    void addToListAndPlay(const QString &uri);
    void setTaskbarButtonWindow();

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

    int getPrimaryScreenWidth();
    int getPrimaryScreenHeight();

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

private slots:
    void onUpdateTime();
    void onUpdateLrc();
    void onPlay();
    void onPlayStop();
    void onPlayPrevious();
    void onPlayNext();

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
    void primaryScreenWidthChanged();
    void primaryScreenHeightChanged();

private:
    QTimer *m_timer {nullptr};
    QTimer *m_lrcTimer {nullptr};
    Lyrics *m_lyrics {nullptr};
    OSD *   m_osd {nullptr};
    LrcBar *m_lb {nullptr};
    float   m_fftData[2048];
    double  m_fftBarValue[29];
    double  m_fftBarPeakValue[29];
    int     m_oriFreq {0};
    bool    m_playing {false};

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

#if defined(Q_OS_WIN) && QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
    QWinTaskbarButton *  taskbarButton;
    QWinTaskbarProgress *taskbarProgress;

    QWinThumbnailToolBar *   thumbnailToolBar;
    QWinThumbnailToolButton *playToolButton;
    QWinThumbnailToolButton *stopToolButton;
    QWinThumbnailToolButton *backwardToolButton;
    QWinThumbnailToolButton *forwardToolButton;
#endif

    void  applyEQ();
    float arraySUM(int start, int end, float *array);
    void  updateFFT();
    void  showCoverPic(const QString &filePath);
    void  infoLabelAnimation();
    void  drawFFTBar(QWidget *parent, int x, int y, int width, int height, double percent);
    void  drawFFTBarPeak(QWidget *parent, int x, int y, int width, int height, double percent);
};

inline QmlPlayer *            gQmlPlayer            = nullptr;
inline QQmlApplicationEngine *gQmlApplicationEngine = nullptr;

#endif // QMLPLAYER_H
