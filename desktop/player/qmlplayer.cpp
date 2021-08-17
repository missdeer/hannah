#include <QCoreApplication>
#include <QTimer>

#include "qmlplayer.h"
#include "FlacPic.h"
#include "ID3v2Pic.h"
#include "configurationwindow.h"
#include "lrcbar.h"
#include "lyrics.h"
#include "osd.h"
#include "player.h"
#include "playlist.h"
#include "playlistmanagewindow.h"

QmlPlayer::QmlPlayer(QObject *parent)
    : QObject(parent),
      m_timer(new QTimer()),
      m_lrcTimer(new QTimer()),
      m_player(new Player()),
      m_lyrics(new Lyrics()),
      m_osd(new OSD()),
      m_lb(new LrcBar(m_lyrics, m_player))
{
    m_player->devInit();
    connect(m_timer, SIGNAL(timeout()), this, SLOT(onUpdateTime()));
    connect(m_lrcTimer, SIGNAL(timeout()), this, SLOT(onUpdateLrc()));

    m_timer->start(27);
    m_lrcTimer->start(70);

#if defined(Q_OS_WIN) && QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
    taskbarButton = new QWinTaskbarButton(this);
    taskbarButton->setWindow(windowHandle());
    taskbarProgress = taskbarButton->progress();
    taskbarProgress->setRange(0, 1000);

    thumbnailToolBar   = new QWinThumbnailToolBar(this);
    playToolButton     = new QWinThumbnailToolButton(thumbnailToolBar);
    stopToolButton     = new QWinThumbnailToolButton(thumbnailToolBar);
    backwardToolButton = new QWinThumbnailToolButton(thumbnailToolBar);
    forwardToolButton  = new QWinThumbnailToolButton(thumbnailToolBar);
    playToolButton->setToolTip(tr("Player"));
    playToolButton->setIcon(QIcon(":/rc/images/player/Play.png"));
    stopToolButton->setToolTip(tr("Stop"));
    stopToolButton->setIcon(QIcon(":/rc/images/player/Stop.png"));
    backwardToolButton->setToolTip(tr("Previous"));
    backwardToolButton->setIcon(QIcon(":/rc/images/player/Pre.png"));
    forwardToolButton->setToolTip(tr("Next"));
    forwardToolButton->setIcon(QIcon(":/rc/images/player/Next.png"));
    thumbnailToolBar->addButton(playToolButton);
    thumbnailToolBar->addButton(stopToolButton);
    thumbnailToolBar->addButton(backwardToolButton);
    thumbnailToolBar->addButton(forwardToolButton);
    connect(playToolButton, SIGNAL(clicked()), this, SLOT(onPlay()));
    connect(stopToolButton, SIGNAL(clicked()), this, SLOT(onPlayStop()));
    connect(backwardToolButton, SIGNAL(clicked()), this, SLOT(onPlayPrevious()));
    connect(forwardToolButton, SIGNAL(clicked()), this, SLOT(onPlayNext()));
#endif
}

void QmlPlayer::onQuit()
{
    QCoreApplication::quit();
}

void QmlPlayer::onShowPlaylists()
{
    Q_ASSERT(playlistManageWindow);
    if (playlistManageWindow->isHidden())
        playlistManageWindow->showNormal();
    playlistManageWindow->raise();
    playlistManageWindow->activateWindow();
}

void QmlPlayer::onSettings()
{
    Q_ASSERT(configurationWindow);
    configurationWindow->onShowConfiguration();
}

void QmlPlayer::onFilter() {}

void QmlPlayer::onMessage() {}

void QmlPlayer::onMusic() {}

void QmlPlayer::onCloud() {}

void QmlPlayer::onBluetooth() {}

void QmlPlayer::onCart() {}

void QmlPlayer::presetEQChanged(int index)
{
    QVector<QVector<int>> presets = {{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
                                     {3, 1, 0, -2, -4, -4, -2, 0, 1, 2},
                                     {-2, 0, 2, 4, -2, -2, 0, 0, 4, 4},
                                     {-6, 1, 4, -2, -2, -4, 0, 0, 6, 6},
                                     {0, 8, 8, 4, 0, 0, 0, 0, 2, 2},
                                     {-6, 0, 0, 0, 0, 0, 4, 0, 4, 0},
                                     {-2, 3, 4, 1, -2, -2, 0, 0, 4, 4},
                                     {-2, 0, 0, 2, 2, 0, 0, 0, 4, 4},
                                     {0, 0, 0, 4, 4, 4, 0, 2, 3, 4},
                                     {-2, 0, 2, 1, 0, 0, 0, 0, -2, -4},
                                     {-4, 0, 2, 1, 0, 0, 0, 0, -4, -6},
                                     {0, 0, 0, 4, 5, 3, 6, 3, 0, 0},
                                     {-4, 0, 2, 0, 0, 0, 0, 0, -4, -6}};
    if (index >= 0 && index < presets.length())
    {
        auto &preset = presets[index];
        m_eq0        = preset[0];
        m_eq1        = preset[1];
        m_eq2        = preset[2];
        m_eq3        = preset[3];
        m_eq4        = preset[4];
        m_eq5        = preset[5];
        m_eq6        = preset[6];
        m_eq7        = preset[7];
        m_eq8        = preset[8];
        m_eq9        = preset[9];
        emit eq0Changed();
        emit eq1Changed();
        emit eq2Changed();
        emit eq3Changed();
        emit eq4Changed();
        emit eq5Changed();
        emit eq6Changed();
        emit eq7Changed();
        emit eq8Changed();
        emit eq9Changed();
        applyEQ();
    }
}

void QmlPlayer::onOpenPreset() {}

void QmlPlayer::onSavePreset() {}

void QmlPlayer::onFavorite() {}

void QmlPlayer::onStop() {}

void QmlPlayer::onPrevious() {}

void QmlPlayer::onPause() {}

void QmlPlayer::onNext() {}

void QmlPlayer::onRepeat() {}

void QmlPlayer::onShuffle() {}

void QmlPlayer::onSwitchFiles() {}

void QmlPlayer::onSwitchPlaylists() {}

void QmlPlayer::onSwitchFavourites() {}

void QmlPlayer::onOpenFile() {}

void QmlPlayer::Show()
{
    emit showPlayer();
}

qreal QmlPlayer::getEq0() const
{
    return m_eq0;
}

qreal QmlPlayer::getEq1() const
{
    return m_eq1;
}

qreal QmlPlayer::getEq2() const
{
    return m_eq2;
}

qreal QmlPlayer::getEq3() const
{
    return m_eq3;
}

qreal QmlPlayer::getEq4() const
{
    return m_eq4;
}

qreal QmlPlayer::getEq5() const
{
    return m_eq5;
}

qreal QmlPlayer::getEq6() const
{
    return m_eq6;
}

qreal QmlPlayer::getEq7() const
{
    return m_eq7;
}

qreal QmlPlayer::getEq8() const
{
    return m_eq8;
}

qreal QmlPlayer::getEq9() const
{
    return m_eq9;
}

qreal QmlPlayer::getVolumn() const
{
    return m_volumn;
}

qreal QmlPlayer::getProgress() const
{
    return m_progress;
}

const QString &QmlPlayer::getCoverUrl() const
{
    return m_coverUrl;
}

const QString &QmlPlayer::getSongName() const
{
    return m_songName;
}

void QmlPlayer::setEq0(qreal value)
{
    if (m_eq0 == value)
        return;
    m_eq0 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(0, (int)m_eq0);
}

void QmlPlayer::setEq1(qreal value)
{
    if (m_eq1 == value)
        return;
    m_eq1 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(1, (int)m_eq1);
}

void QmlPlayer::setEq2(qreal value)
{
    if (m_eq2 == value)
        return;
    m_eq2 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(2, (int)m_eq2);
}

void QmlPlayer::setEq3(qreal value)
{
    if (m_eq3 == value)
        return;
    m_eq3 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(3, (int)m_eq3);
}

void QmlPlayer::setEq4(qreal value)
{
    if (m_eq4 == value)
        return;
    m_eq4 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(4, (int)m_eq4);
}

void QmlPlayer::setEq5(qreal value)
{
    if (m_eq5 == value)
        return;
    m_eq5 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(5, (int)m_eq6);
}

void QmlPlayer::setEq6(qreal value)
{
    if (m_eq6 == value)
        return;
    m_eq6 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(6, (int)m_eq6);
}

void QmlPlayer::setEq7(qreal value)
{
    if (m_eq7 == value)
        return;
    m_eq7 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(7, (int)m_eq7);
}

void QmlPlayer::setEq8(qreal value)
{
    if (m_eq8 == value)
        return;
    m_eq8 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(8, (int)m_eq8);
}

void QmlPlayer::setEq9(qreal value)
{
    if (m_eq9 == value)
        return;
    m_eq9 = value;
    Q_ASSERT(m_player);
    m_player->setEQ(9, (int)m_eq9);
}

void QmlPlayer::setVolumn(qreal value)
{
    if (m_volumn == value)
        return;
    m_volumn = value;
    Q_ASSERT(m_player);
    m_player->setVol((int)m_volumn);
}

void QmlPlayer::setProgress(qreal progress)
{
    if (m_progress == progress)
        return;
    m_progress = progress;
}

void QmlPlayer::setCoverUrl(const QString &u)
{
    if (m_coverUrl == u)
        return;
    m_coverUrl = u;
}

void QmlPlayer::setSongName(const QString &n)
{
    if (m_songName == n)
        return;
    m_songName = n;
}

void QmlPlayer::onPlay() {}

void QmlPlayer::onPlayStop() {}

void QmlPlayer::onPlayPrevious() {}

void QmlPlayer::onPlayNext() {}

void QmlPlayer::onUpdateTime() {}

void QmlPlayer::onUpdateLrc() {}

void QmlPlayer::applyEQ()
{
    Q_ASSERT(m_player);
    m_player->setEQ(0, (int)m_eq0);
    m_player->setEQ(1, (int)m_eq1);
    m_player->setEQ(2, (int)m_eq2);
    m_player->setEQ(3, (int)m_eq3);
    m_player->setEQ(4, (int)m_eq4);
    m_player->setEQ(5, (int)m_eq5);
    m_player->setEQ(6, (int)m_eq6);
    m_player->setEQ(7, (int)m_eq7);
    m_player->setEQ(8, (int)m_eq8);
    m_player->setEQ(9, (int)m_eq9);
}
