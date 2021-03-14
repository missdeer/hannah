#include <QFileInfo>
#include <QRandomGenerator>

#include "shadowplayer.h"
#include "FlacPic.h"
#include "ID3v2Pic.h"
#include "lrcbar.h"
#include "lyrics.h"
#include "osd.h"
#include "player.h"
#include "playlist.h"
#include "ui_shadowplayer.h"

ShadowPlayer::ShadowPlayer(QWidget *parent)
    : QMainWindow(parent),
      ui(new Ui::ShadowPlayer),
      timer(new QTimer()),
      lrcTimer(new QTimer()),
      player(new Player()),
      lyrics(new Lyrics()),
      osd(new OSD()),
      lb(new LrcBar(lyrics, player, 0))
{
    ui->setupUi(this);
    // setFixedSize(width(), height());

    playList = new PlayList(player, ui->playerListArea);

    setWindowIcon(QIcon(":/rc/images/player/ShadowPlayer.ico"));
    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowSystemMenuHint | Qt::WindowMinimizeButtonHint);
    setAttribute(Qt::WA_TranslucentBackground, true);
    ui->coverLabel->setScaledContents(true);
    ui->coverLabel->setPixmap(QPixmap(":/rc/images/player/ShadowPlayer.png"));

    ui->tagLabel->setShadowMode(1);
    ui->mediaInfoLabel->setShadowMode(1);
    ui->tagLabel->setShadowColor(QColor(0, 0, 0, 80));
    ui->mediaInfoLabel->setShadowColor(QColor(0, 0, 0, 80));
    ui->curTimeLabel->setShadowMode(1);
    ui->totalTimeLabel->setShadowMode(1);
    ui->curTimeLabel->setShadowColor(QColor(0, 0, 0, 80));
    ui->totalTimeLabel->setShadowColor(QColor(0, 0, 0, 80));

    player->devInit();

    connect(timer, SIGNAL(timeout()), this, SLOT(UpdateTime()));
    connect(lrcTimer, SIGNAL(timeout()), this, SLOT(UpdateLrc()));

    connect(ui->eqSlider_1, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_2, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_3, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_4, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_5, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_6, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_7, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_8, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_9, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));
    connect(ui->eqSlider_10, SIGNAL(valueChanged(int)), this, SLOT(applyEQ()));

    connect(playList, SIGNAL(callPlayer()), this, SLOT(callFromPlayList()));

    timer->start(27);
    lrcTimer->start(70);

    setAcceptDrops(true);
    loadSkin(":/rc/images/Skin1.jpg", false);
    ui->playerListArea->setGeometry(QRect(370, 200, 331, 0));

    loadConfig();

    bgLinearGradient.setColorAt(0, QColor(255, 255, 255, 0));
    bgLinearGradient.setColorAt(1, QColor(255, 255, 255, 255));
    bgLinearGradient.setStart(0, 0);
    bgLinearGradient.setFinalStop(0, height());

#if defined(Q_OS_WIN)
    taskbarButton = new QWinTaskbarButton(this);
    taskbarButton->setWindow(windowHandle());
    taskbarProgress = taskbarButton->progress();
    taskbarProgress->setRange(0, 1000);
    connect(ui->playSlider, SIGNAL(valueChanged(int)), taskbarProgress, SLOT(setValue(int)));

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
    connect(playToolButton, SIGNAL(clicked()), this, SLOT(on_playButton_clicked()));
    connect(stopToolButton, SIGNAL(clicked()), this, SLOT(on_stopButton_clicked()));
    connect(backwardToolButton, SIGNAL(clicked()), this, SLOT(on_playPreButton_clicked()));
    connect(forwardToolButton, SIGNAL(clicked()), this, SLOT(on_playNextButton_clicked()));
#endif
    sizeSlideAnimation = new QPropertyAnimation(this, "geometry");
    sizeSlideAnimation->setDuration(700);
    sizeSlideAnimation->setEasingCurve(QEasingCurve::OutCirc);

    tagAnimation = new QPropertyAnimation(ui->tagLabel, "geometry");
    tagAnimation->setDuration(700);
    tagAnimation->setStartValue(QRect(130, 30, 0, 16));
    tagAnimation->setEndValue(QRect(130, 30, 221, 16));
    mediaInfoAnimation = new QPropertyAnimation(ui->mediaInfoLabel, "geometry");
    mediaInfoAnimation->setDuration(800);
    mediaInfoAnimation->setStartValue(QRect(130, 50, 0, 16));
    mediaInfoAnimation->setEndValue(QRect(130, 50, 221, 16));
    coverAnimation = new QPropertyAnimation(ui->coverLabel, "geometry");
    coverAnimation->setEasingCurve(QEasingCurve::OutCirc);
    coverAnimation->setDuration(600);
    coverAnimation->setStartValue(QRect(60, 60, 1, 1));
    coverAnimation->setEndValue(QRect(10, 10, 111, 111));

    eqHideAnimation = new QPropertyAnimation(ui->eqGroupBox, "geometry");
    eqHideAnimation->setDuration(600);
    eqHideAnimation->setStartValue(QRect(370, 30, 331, 171));
    eqHideAnimation->setEndValue(QRect(370, 30, 331, 0));

    eqShowAnimation = new QPropertyAnimation(ui->eqGroupBox, "geometry");
    eqShowAnimation->setDuration(600);
    eqShowAnimation->setStartValue(QRect(370, 30, 331, 0));
    eqShowAnimation->setEndValue(QRect(370, 30, 331, 171));

    lyricsHideAnimation = new QPropertyAnimation(ui->lyricsBox, "geometry");
    lyricsHideAnimation->setDuration(600);
    lyricsHideAnimation->setStartValue(QRect(370, 210, 331, 181));
    lyricsHideAnimation->setEndValue(QRect(370, 391, 331, 0));

    lyricsShowAnimation = new QPropertyAnimation(ui->lyricsBox, "geometry");
    lyricsShowAnimation->setDuration(600);
    lyricsShowAnimation->setStartValue(QRect(370, 391, 331, 0));
    lyricsShowAnimation->setEndValue(QRect(370, 210, 331, 181));

    playListHideAnimation = new QPropertyAnimation(ui->playerListArea, "geometry");
    playListHideAnimation->setDuration(600);
    playListHideAnimation->setStartValue(QRect(370, 30, 331, 361));
    playListHideAnimation->setEndValue(QRect(370, 205, 331, 0));

    playListShowAnimation = new QPropertyAnimation(ui->playerListArea, "geometry");
    playListShowAnimation->setDuration(600);
    playListShowAnimation->setStartValue(QRect(370, 205, 331, 0));
    playListShowAnimation->setEndValue(QRect(370, 30, 331, 361));

    connect(playListHideAnimation, SIGNAL(finished()), this, SLOT(update()));
    connect(playListShowAnimation, SIGNAL(finished()), this, SLOT(update())); //动画完成后刷新窗口

    fadeOutAnimation = new QPropertyAnimation(this, "windowOpacity");
    fadeOutAnimation->setDuration(400);
    fadeOutAnimation->setStartValue(1);
    fadeOutAnimation->setEndValue(0);
    connect(fadeOutAnimation, SIGNAL(finished()), this, SLOT(close()));

    fadeInAnimation = new QPropertyAnimation(this, "windowOpacity");
    fadeInAnimation->setDuration(400);
    fadeInAnimation->setStartValue(0);
    fadeInAnimation->setEndValue(1);
    fadeInAnimation->start();
    connect(fadeInAnimation, SIGNAL(finished()), this, SLOT(setTaskbarButtonWindow()));
}

ShadowPlayer::~ShadowPlayer()
{
    delete ui;
    delete lb;
    delete player;
    delete timer;
    delete lyrics;
    delete playList;
    delete osd;
    delete sizeSlideAnimation;
    delete fadeInAnimation;
    delete tagAnimation;
    delete mediaInfoAnimation;
    delete coverAnimation;
    delete fadeOutAnimation;
    delete eqHideAnimation;
    delete eqShowAnimation;
    delete lyricsHideAnimation;
    delete lyricsShowAnimation;
    delete playListHideAnimation;
    delete playListShowAnimation;
}

void ShadowPlayer::dragEnterEvent(QDragEnterEvent *event)
{
    event->acceptProposedAction();
}

void ShadowPlayer::dropEvent(QDropEvent *event)
{
    QList<QUrl> urls = event->mimeData()->urls();
    if (urls.isEmpty())
        return;

    QString fileName = urls.first().toLocalFile();
    if (fileName.isEmpty())
        return;

    QFileInfo fi(fileName);
    QString   ext = fi.suffix();

    if (ext.compare("jpg", Qt::CaseInsensitive) == 0 || ext.compare("jpeg", Qt::CaseInsensitive) == 0 ||
        ext.compare("png", Qt::CaseInsensitive) == 0 || ext.compare("gif", Qt::CaseInsensitive) == 0 || ext.compare("bmp", Qt::CaseInsensitive) == 0)
    {
        loadSkin(fileName);
    }
    else if (ext.compare("lrc", Qt::CaseInsensitive) == 0)
    {
        lyrics->resolve(fileName, true);
    }
    else
    {
        addToListAndPlay(urls);
    }
}

void ShadowPlayer::paintEvent(QPaintEvent *)
{
    QPainter painter(this);

    if (!skin.isNull())
    {
        int topY = 0;

        switch (skinPos)
        {
        case 0:
            skinDrawPos = 0;
            break;
        case 1:
            skinDrawPos = 0.5;
            break;
        case 2:
            skinDrawPos = 1;
            break;
        default:
            break;
        }

        switch (skinMode)
        {
        case 0:
            topY = -(skin.height() - 400) * skinDrawPos;
            painter.drawPixmap(0, topY, skin);
            break;
        case 1:
            topY = -(skinLeft.height() - 400) * skinDrawPos;
            painter.drawPixmap(0, topY, skinLeft);
            break;
        case 2:
            topY = -(skinFull.height() - 400) * skinDrawPos;
            painter.drawPixmap(0, topY, skinFull);
            break;
        case 3:
            if (geometry().width() > 361)
            {
                topY = -(skinFull.height() - 400) * skinDrawPos;
                painter.drawPixmap(0, topY, skinFull);
            }
            else
            {
                topY = -(skinLeft.height() - 400) * skinDrawPos;
                painter.drawPixmap(0, topY, skinLeft);
            }
            break;
        case 4:
            if (geometry().width() <= 360)
            {
                topY = -(skinLeft.height() - 400) * skinDrawPos;
                painter.drawPixmap(0, topY, skinLeft);
            }
            else if (geometry().width() >= 710)
            {
                topY = -(skinFull.height() - 400) * skinDrawPos;
                painter.drawPixmap(0, topY, skinFull);
            }
            else
            {
                topY = -(width() * aspectRatio - 400) * skinDrawPos;
                painter.drawPixmap(0, topY, width(), width() * aspectRatio, skin);
            }
            break;
        default:
            break;
        }
    }

    painter.setBrush(QBrush(bgLinearGradient));
    painter.setPen(Qt::NoPen);
    painter.drawRect(0, 0, width(), height());
}

void ShadowPlayer::contextMenuEvent(QContextMenuEvent *)
{
    QCursor cur = cursor();

    QAction textAction1(tr("Drag left side of the window to adjust position accurately"), this);
    textAction1.setEnabled(false);

    QMenu menu;
    menu.addAction(tr("Load default skin"), this, SLOT(loadDefaultSkin()));
    menu.addSeparator();
    menu.addAction(tr("Fix skin to size side"), this, SLOT(fixSkinSizeLeft()));
    menu.addAction(tr("Fix skin to full window"), this, SLOT(fixSkinSizeFull()));
    menu.addAction(tr("Original skin size"), this, SLOT(originalSkinSize()));
    menu.addAction(tr("Auto resize skin"), this, SLOT(autoSkinSize()));
    menu.addAction(tr("Resize skin dynamically"), this, SLOT(dynamicSkinSize()));
    menu.addSeparator();
    menu.addAction(tr("Skin on top"), this, SLOT(skinOnTop()));
    menu.addAction(tr("Skin on center"), this, SLOT(skinOnCenter()));
    menu.addAction(tr("Skin on bottom"), this, SLOT(skinOnBottom()));
    menu.addAction(&textAction1);
    menu.addSeparator();
    menu.addAction(tr("Disable skin"), this, SLOT(skinDisable()));
    menu.addSeparator();
    menu.addAction(tr("Physical settings"), this, SLOT(physicsSetting()));
    menu.addAction(tr("Enable Physical FFT"), this, SLOT(enableFFTPhysics()));
    menu.addAction(tr("Disable Physical FFT"), this, SLOT(disableFFTPhysics()));
    menu.addSeparator();
    menu.addAction(tr("About"), this, SLOT(showDeveloperInfo()));
    menu.exec(cur.pos());
}

void ShadowPlayer::loadDefaultSkin()
{
    loadSkin(":/image/image/Skin1.jpg", false);
    QFile skinFile(QCoreApplication::applicationDirPath() + "/skin.dat");
    skinFile.remove();
}

void ShadowPlayer::fixSkinSizeLeft()
{
    skinMode = 1;
    repaint();
}

void ShadowPlayer::fixSkinSizeFull()
{
    skinMode = 2;
    repaint();
}

void ShadowPlayer::originalSkinSize()
{
    skinMode = 0;
    repaint();
}

void ShadowPlayer::autoSkinSize()
{
    skinMode = 3;
    repaint();
}

void ShadowPlayer::dynamicSkinSize()
{
    skinMode = 4;
    repaint();
}

void ShadowPlayer::skinOnTop()
{
    skinPos = 0;
    repaint();
}

void ShadowPlayer::skinOnCenter()
{
    skinPos = 1;
    repaint();
}

void ShadowPlayer::skinOnBottom()
{
    skinPos = 2;
    repaint();
}

void ShadowPlayer::skinDisable()
{
    skin     = QPixmap();
    skinLeft = skin;
    skinFull = skin;
    repaint();
}

void ShadowPlayer::infoLabelAnimation()
{
    tagAnimation->stop();
    mediaInfoAnimation->stop();
    coverAnimation->stop();

    tagAnimation->start();
    mediaInfoAnimation->start();
    coverAnimation->start();
    update();
}

void ShadowPlayer::loadFile(QString file)
{
    if (player->openFile(file) != "err")
    {
        QFileInfo fileinfo(file);
        showCoverPic(file);

        ui->mediaInfoLabel->setText(player->getNowPlayInfo());
        ui->totalTimeLabel->setText(player->getTotalTime());
        oriFreq = player->getFreq();
        on_freqSlider_valueChanged(ui->freqSlider->value());
        player->setVol(ui->volSlider->value());
        applyEQ();

        player->setReverse(isReverse);

        on_eqEnableCheckBox_clicked(ui->eqEnableCheckBox->isChecked());

        player->updateReverb(ui->reverbDial->value());

        player->play();
        if (player->isPlaying())
        {
            playing = true;
            ui->playButton->setIcon(QIcon(":/rc/images/player/Pause.png"));
            ui->playButton->setToolTip(tr("Pause"));
            taskbarProgress->show();
            taskbarProgress->resume();
            taskbarButton->setOverlayIcon(QIcon(":/rc/images/player/Play.png"));
            playToolButton->setIcon(QIcon(":/rc/images/player/Pause.png"));
            playToolButton->setToolTip(tr("Pause"));
        }

        if (!lyrics->resolve(file))
            if (!lyrics->loadFromLrcDir(file))
                if (!lyrics->loadFromFileRelativePath(file, "/Lyrics/"))
                    lyrics->loadFromFileRelativePath(file, "/../Lyrics/");

        if (player->getTags() == tr("Show_File_Name"))
            ui->tagLabel->setText(fileinfo.fileName());
        else
            ui->tagLabel->setText(player->getTags());

        infoLabelAnimation();
        osd->showOSD(ui->tagLabel->text(), player->getTotalTime());
    }
}

void ShadowPlayer::loadSkin(QString image, bool save)
{
    skin        = QPixmap(image);
    skinLeft    = skin.scaledToWidth(360, Qt::SmoothTransformation);
    skinFull    = skin.scaledToWidth(710, Qt::SmoothTransformation);
    aspectRatio = (double)skin.height() / skin.width();
    if (save)
        saveSkinData();
    update();
}

void ShadowPlayer::UpdateTime()
{
    ui->leftLevel->setValue(LOWORD(player->getLevel()));
    ui->rightLevel->setValue(HIWORD(player->getLevel()));
    ui->leftLevel->update();
    ui->rightLevel->update();

    double ldB = 20 * log((double)ui->leftLevel->value() / 32768);
    double rdB = 20 * log((double)ui->rightLevel->value() / 32768);
    if (ldB < -60)
        ldB = -60;
    if (rdB < -60)
        rdB = -60;
    ui->leftdB->setValue(ldB);
    ui->rightdB->setValue(rdB);
    ui->leftdB->update();
    ui->rightdB->update();

    if (skinPos == 3)
        update();

    ui->curTimeLabel->setText(player->getCurTime());
    if (isPlaySliderPress == false)
        ui->playSlider->setSliderPosition(player->getPos());

    updateFFT();

    if (!player->isPlaying())
    {
        switch (playMode)
        {
        case 0:
            playing = false;
            ui->playButton->setIcon(QIcon(":/rc/images/player/Play.png"));
            ui->playButton->setToolTip(tr("Play"));
            taskbarProgress->hide();
            taskbarButton->setOverlayIcon(QIcon(":/rc/images/player/Stop.png"));
            playToolButton->setIcon(QIcon(":/rc/images/player/Play.png"));
            playToolButton->setToolTip(tr("Play"));
            break;
        case 1:
            if (playing)
            {
                player->stop();
                player->play();
            }
            break;
        case 2:
            if (playing)
            {
                QString nextFile = playList->next(false);
                if (nextFile == "stop")
                {
                    playing = false;
                    ui->playButton->setIcon(QIcon(":/rc/images/player/Play.png"));
                    ui->playButton->setToolTip(tr("Play"));
                    taskbarProgress->hide();
                    taskbarButton->setOverlayIcon(QIcon(":/rc/images/player/Stop.png"));
                    playToolButton->setIcon(QIcon(":/rc/images/player/Play.png"));
                    playToolButton->setToolTip(tr("Play"));
                }
                else
                {
                    loadFile(nextFile);
                }
            }
            break;
        case 3:
            if (playing)
            {
                loadFile(playList->next(true));
            }
            break;
        case 4:
            if (playing)
            {
                int index = QRandomGenerator::global()->bounded(playList->getLength());
                loadFile(playList->playIndex(index));
            }
        default:
            break;
        }
    }
}

void ShadowPlayer::UpdateLrc()
{
    lyrics->updateTime(player->getCurTimeMS(), player->getTotalTimeMS());
    double pos        = 0;
    double curTimePos = lyrics->getTimePos(player->getCurTimeMS());

    ui->lrcLabel_1->setText(lyrics->getLrcString(-3));
    ui->lrcLabel_2->setText(lyrics->getLrcString(-2));
    ui->lrcLabel_3->setText(lyrics->getLrcString(-1));
    ui->lrcLabel_4->setText(lyrics->getLrcString(0));
    ui->lrcLabel_5->setText(lyrics->getLrcString(1));
    ui->lrcLabel_6->setText(lyrics->getLrcString(2));
    ui->lrcLabel_7->setText(lyrics->getLrcString(3));

    ui->lrcLabel_1->setToolTip(ui->lrcLabel_1->text());
    ui->lrcLabel_2->setToolTip(ui->lrcLabel_2->text());
    ui->lrcLabel_3->setToolTip(ui->lrcLabel_3->text());
    ui->lrcLabel_4->setToolTip(ui->lrcLabel_4->text());
    ui->lrcLabel_5->setToolTip(ui->lrcLabel_5->text());
    ui->lrcLabel_6->setToolTip(ui->lrcLabel_6->text());
    ui->lrcLabel_7->setToolTip(ui->lrcLabel_7->text());

    if (curTimePos >= 0.2 && curTimePos <= 0.8)
    {
        pos = 0.5;
    }
    else if (curTimePos < 0.2)
    {
        pos = curTimePos * 2.5; // 0~0.5
    }
    else if (curTimePos > 0.8)
    {
        pos = (curTimePos - 0.8) * 2.5 + 0.5; // 0.5~1
    }
    ui->lrcLabel_1->setGeometry(10, 35 - 20 * pos, 311, 16);
    // ui->lrcLabel_1->setStyleSheet(QString("color: rgba(0, 0, 0, %1)").arg(245 - 235 * curTimePos));
    ui->lrcLabel_2->setGeometry(10, 55 - 20 * pos, 311, 16);
    ui->lrcLabel_3->setGeometry(10, 75 - 20 * pos, 311, 16);
    ui->lrcLabel_4->setGeometry(10, 95 - 20 * pos, 311, 16);
    ui->lrcLabel_5->setGeometry(10, 115 - 20 * pos, 311, 16);
    ui->lrcLabel_6->setGeometry(10, 135 - 20 * pos, 311, 16);
    ui->lrcLabel_7->setGeometry(10, 155 - 20 * pos, 311, 16);
    // ui->lrcLabel_7->setStyleSheet(QString("color: rgba(0, 0, 0, %1)").arg(235 * curTimePos + 10));
}

void ShadowPlayer::showCoverPic(QString filePath)
{
    QFileInfo fileinfo(filePath);
    QString   path = fileinfo.path();
    if (spID3::loadPictureData(filePath.toLocal8Bit().data()))
    {
        QByteArray picData((const char *)spID3::getPictureDataPtr(), spID3::getPictureLength());
        ui->coverLabel->setPixmap(QPixmap::fromImage(QImage::fromData(picData)));
        spID3::freePictureData();
    }
    else if (spFLAC::loadPictureData(filePath.toLocal8Bit().data()))
    {
        QByteArray picData((const char *)spFLAC::getPictureDataPtr(), spFLAC::getPictureLength());
        ui->coverLabel->setPixmap(QPixmap::fromImage(QImage::fromData(picData)));
        spFLAC::freePictureData();
    }
    else if (QFileInfo(path + "/cover.jpg").exists())
        ui->coverLabel->setPixmap(QPixmap(path + "/cover.jpg"));
    else if (QFileInfo(path + "/cover.jpeg").exists())
        ui->coverLabel->setPixmap(QPixmap(path + "/cover.jpeg"));
    else if (QFileInfo(path + "/cover.png").exists())
        ui->coverLabel->setPixmap(QPixmap(path + "/cover.png"));
    else if (QFileInfo(path + "/cover.gif").exists())
        ui->coverLabel->setPixmap(QPixmap(path + "/cover.gif"));
    else
        ui->coverLabel->setPixmap(QPixmap(":image/image/ShadowPlayer.png"));
}

void ShadowPlayer::on_openButton_clicked()
{
    QStringList files = QFileDialog::getOpenFileNames(
        this, tr("Open"), 0, tr("Audio file (*.mp3 *.mp2 *.mp1 *.wav *.aiff *.ogg *.ape *.mp4 *.m4a *.m4v *.aac *.alac *.tta *.flac *.wma *.wv)"));
    int newIndex = playList->getLength();
    int length   = files.length();
    if (!files.isEmpty())
    {
        for (int i = 0; i < length; i++)
        {
            playList->add(files[i]);
        }
        if (playList->getLength() > newIndex)
        {
            loadFile(playList->playIndex(newIndex));
        }
    }
}

void ShadowPlayer::on_playButton_clicked()
{
    if (!playing)
    {
        player->play();

        if (player->isPlaying())
        {
            playing = true;
            ui->playButton->setIcon(QIcon(":/rc/images/player/Pause.png"));
            ui->playButton->setToolTip(tr("Pause"));
            taskbarProgress->show();
            taskbarProgress->resume();
            taskbarButton->setOverlayIcon(QIcon(":/rc/images/player/Play.png"));
            playToolButton->setIcon(QIcon(":/rc/images/player/Pause.png"));
            playToolButton->setToolTip(tr("Pause"));
        }
        else
        {
            if (playList->getLength() > 0)
                loadFile(playList->playIndex(playList->getIndex()));
        }
    }
    else
    {
        player->pause();
        playing = false;
        ui->playButton->setIcon(QIcon(":/rc/images/player/Play.png"));
        ui->playButton->setToolTip(tr("Play"));
        taskbarProgress->show();
        taskbarProgress->pause();
        taskbarButton->setOverlayIcon(QIcon(":/rc/images/player/Pause.png"));
        playToolButton->setIcon(QIcon(":/rc/images/player/Play.png"));
        playToolButton->setToolTip(tr("Play"));
    }
}

void ShadowPlayer::on_stopButton_clicked()
{
    player->stop();
    playing = false;
    ui->playButton->setIcon(QIcon(":/rc/images/player/Play.png"));
    ui->playButton->setToolTip(tr("Play"));
    taskbarProgress->hide();
    taskbarButton->setOverlayIcon(QIcon(":/rc/images/player/Stop.png"));
    playToolButton->setIcon(QIcon(":/rc/images/player/Play.png"));
    playToolButton->setToolTip(tr("Play"));
}

void ShadowPlayer::on_volSlider_valueChanged(int value)
{
    player->setVol(value);
    isMute = false;
    ui->muteButton->setIcon(QIcon(":/rc/images/player/Vol.png"));
}

void ShadowPlayer::on_muteButton_clicked()
{
    if (isMute == false)
    {
        lastVol = ui->volSlider->value();
        ui->volSlider->setValue(0);
        ui->muteButton->setIcon(QIcon(":/rc/images/player/Mute.png"));
        isMute = true;
    }
    else
    {
        ui->volSlider->setValue(lastVol);
        ui->muteButton->setIcon(QIcon(":/rc/images/player/Vol.png"));
        isMute = false;
    }
}

float ShadowPlayer::arraySUM(int start, int end, float *array)
{
    float sum = 0;
    for (int i = start; i <= end; i++)
    {
        sum += array[i];
    }
    return sum;
}

void ShadowPlayer::fullZero(int length, float *array)
{
    for (int i = 0; i < length; i++)
    {
        array[i] = 0;
    }
}

void ShadowPlayer::updateFFT()
{
    if (player->isPlaying())
    {
        if (ui->leftLevel->value() > 6 || ui->rightLevel->value() > 6)
            player->getFFT(fftData);
        else
            fullZero(2048, fftData);

        double start = 5;
        for (int i = 0; i < 29; i++)
        {
            double end                      = start * 1.23048;
            ui->FFTGroupBox->fftBarValue[i] = sqrt(arraySUM((int)start, (int)end, fftData));
            start                           = end;
        }

        ui->FFTGroupBox->cutValue();
        ui->FFTGroupBox->peakSlideDown();
        ui->FFTGroupBox->update();
    }
    else
    {
        ui->FFTGroupBox->peakSlideDown();
        ui->FFTGroupBox->update();
    }
}

void ShadowPlayer::on_playSlider_sliderPressed()
{
    isPlaySliderPress = true;
}

void ShadowPlayer::on_playSlider_sliderReleased()
{
    isPlaySliderPress = false;
    player->setPos(ui->playSlider->sliderPosition());
}

void ShadowPlayer::on_resetFreqButton_clicked()
{
    ui->freqSlider->setSliderPosition(0);
    player->setFreq(oriFreq);
}

void ShadowPlayer::applyEQ()
{
    player->setEQ(0, ui->eqSlider_1->value());
    player->setEQ(1, ui->eqSlider_2->value());
    player->setEQ(2, ui->eqSlider_3->value());
    player->setEQ(3, ui->eqSlider_4->value());
    player->setEQ(4, ui->eqSlider_5->value());
    player->setEQ(5, ui->eqSlider_6->value());
    player->setEQ(6, ui->eqSlider_7->value());
    player->setEQ(7, ui->eqSlider_8->value());
    player->setEQ(8, ui->eqSlider_9->value());
    player->setEQ(9, ui->eqSlider_10->value());
}

void ShadowPlayer::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton && event->x() > 10)
    {
        pos          = event->pos();
        clickOnFrame = true;
        event->accept();
    }
    else if (event->button() == Qt::LeftButton && event->x() <= 10)
    {
        clickOnLeft = true;
    }
}

void ShadowPlayer::mouseReleaseEvent(QMouseEvent *)
{
    clickOnFrame = false;
    clickOnLeft  = false;
}

void ShadowPlayer::mouseMoveEvent(QMouseEvent *event)
{
    if (event->buttons() & Qt::LeftButton && clickOnFrame && event->x() > 10)
    {
        move(event->globalPos() - pos);
        event->accept();
    }
    else if (event->buttons() & Qt::LeftButton && clickOnLeft && event->y() >= 0 && event->y() <= 400)
    {
        skinDrawPos = (double)event->y() / 400;
        skinPos     = -1;
        update();
    }
}

void ShadowPlayer::on_extraButton_clicked()
{
    if (geometry().width() < 535)
    {
        sizeSlideAnimation->stop();
        sizeSlideAnimation->setStartValue(QRect(geometry().x(), geometry().y(), 360, 402));
        sizeSlideAnimation->setEndValue(QRect(geometry().x() - 175, geometry().y(), 710, 400));
        sizeSlideAnimation->start();
        ui->extraButton->setText("<-");
    }
    else
    {
        sizeSlideAnimation->stop();
        sizeSlideAnimation->setStartValue(QRect(geometry().x(), geometry().y(), 710, 402));
        sizeSlideAnimation->setEndValue(QRect(geometry().x() + 175, geometry().y(), 360, 400));
        sizeSlideAnimation->start();
        ui->extraButton->setText("->");
    }
}

void ShadowPlayer::on_closeButton_clicked()
{
    fadeOutAnimation->start();
}

void ShadowPlayer::closeEvent(QCloseEvent *event)
{
#if defined(Q_OS_MACOS)
    if (!event->spontaneous() || !isVisible())
    {
        return;
    }
#endif
    hide();
    event->ignore();
}

void ShadowPlayer::on_setSkinButton_clicked()
{
    QString image = QFileDialog::getOpenFileName(this, tr("Choose skin"), 0, tr("Image file (*.bmp *.jpg *.jpeg *.png *.gif)"));
    if (!image.isEmpty())
    {
        loadSkin(image);
    }
}

void ShadowPlayer::on_miniSizeButton_clicked()
{
    showMinimized();
}

void ShadowPlayer::on_playModeButton_clicked()
{
    playMode = ++playMode % 5;

    switch (playMode)
    {
    case 0:
        ui->playModeButton->setIcon(QIcon(":/rc/images/player/Single.png"));
        ui->playModeButton->setToolTip(tr("Track Play"));
        break;
    case 1:
        ui->playModeButton->setIcon(QIcon(":/rc/images/player/Repeat.png"));
        ui->playModeButton->setToolTip(tr("Track Repeat"));
        break;
    case 2:
        ui->playModeButton->setIcon(QIcon(":/rc/images/player/Order.png"));
        ui->playModeButton->setToolTip(tr("Playlist Order"));
        break;
    case 3:
        ui->playModeButton->setIcon(QIcon(":/rc/images/player/AllRepeat.png"));
        ui->playModeButton->setToolTip(tr("Playlist Repeat"));
        break;
    case 4:
        ui->playModeButton->setIcon(QIcon(":/rc/images/player/Shuffle.png"));
        ui->playModeButton->setToolTip(tr("Shuffle"));
        break;
    default:
        break;
    }
    QToolTip::showText(QCursor::pos(), ui->playModeButton->toolTip());
}

void ShadowPlayer::on_showDskLrcButton_clicked()
{
    if (lb->isVisible())
        lb->hide();
    else
        lb->show();
}

void ShadowPlayer::on_loadLrcButton_clicked()
{
    QString file = QFileDialog::getOpenFileName(this, tr("Load lyric"), 0, tr("Lyric file (*.lrc)"));
    if (QFile::exists(file))
        lyrics->resolve(file, true);
}

void ShadowPlayer::on_playSlider_valueChanged(int value)
{
    // avoiding setValue being triggered recursively
    if (qAbs(player->getPos() - value) > 2)
    {
        player->setPos(value);
    }
}

void ShadowPlayer::on_freqSlider_valueChanged(int value)
{
    player->setFreq(oriFreq + (oriFreq * value * 0.0001));
    ui->textLabel1->setText(tr("Playback speed (x%1)").arg(value * 0.0001 + 1));
}

void ShadowPlayer::on_eqComboBox_currentIndexChanged(int index)
{
    switch (index)
    {
    case 0:
        ui->eqSlider_1->setSliderPosition(0);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(0);
        ui->eqSlider_4->setSliderPosition(0);
        ui->eqSlider_5->setSliderPosition(0);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(0);
        ui->eqSlider_10->setSliderPosition(0);
        applyEQ();
        break;
    case 1:
        ui->eqSlider_1->setSliderPosition(3);
        ui->eqSlider_2->setSliderPosition(1);
        ui->eqSlider_3->setSliderPosition(0);
        ui->eqSlider_4->setSliderPosition(-2);
        ui->eqSlider_5->setSliderPosition(-4);
        ui->eqSlider_6->setSliderPosition(-4);
        ui->eqSlider_7->setSliderPosition(-2);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(1);
        ui->eqSlider_10->setSliderPosition(2);
        applyEQ();
        break;
    case 2:
        ui->eqSlider_1->setSliderPosition(-2);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(2);
        ui->eqSlider_4->setSliderPosition(4);
        ui->eqSlider_5->setSliderPosition(-2);
        ui->eqSlider_6->setSliderPosition(-2);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(4);
        ui->eqSlider_10->setSliderPosition(4);
        applyEQ();
        break;
    case 3:
        ui->eqSlider_1->setSliderPosition(-6);
        ui->eqSlider_2->setSliderPosition(1);
        ui->eqSlider_3->setSliderPosition(4);
        ui->eqSlider_4->setSliderPosition(-2);
        ui->eqSlider_5->setSliderPosition(-2);
        ui->eqSlider_6->setSliderPosition(-4);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(6);
        ui->eqSlider_10->setSliderPosition(6);
        applyEQ();
        break;
    case 4:
        ui->eqSlider_1->setSliderPosition(0);
        ui->eqSlider_2->setSliderPosition(8);
        ui->eqSlider_3->setSliderPosition(8);
        ui->eqSlider_4->setSliderPosition(4);
        ui->eqSlider_5->setSliderPosition(0);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(2);
        ui->eqSlider_10->setSliderPosition(2);
        applyEQ();
        break;
    case 5:
        ui->eqSlider_1->setSliderPosition(-6);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(0);
        ui->eqSlider_4->setSliderPosition(0);
        ui->eqSlider_5->setSliderPosition(0);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(4);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(4);
        ui->eqSlider_10->setSliderPosition(0);
        applyEQ();
        break;
    case 6:
        ui->eqSlider_1->setSliderPosition(-2);
        ui->eqSlider_2->setSliderPosition(3);
        ui->eqSlider_3->setSliderPosition(4);
        ui->eqSlider_4->setSliderPosition(1);
        ui->eqSlider_5->setSliderPosition(-2);
        ui->eqSlider_6->setSliderPosition(-2);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(4);
        ui->eqSlider_10->setSliderPosition(4);
        applyEQ();
        break;
    case 7:
        ui->eqSlider_1->setSliderPosition(-2);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(0);
        ui->eqSlider_4->setSliderPosition(2);
        ui->eqSlider_5->setSliderPosition(2);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(4);
        ui->eqSlider_10->setSliderPosition(4);
        applyEQ();
        break;
    case 8:
        ui->eqSlider_1->setSliderPosition(0);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(0);
        ui->eqSlider_4->setSliderPosition(4);
        ui->eqSlider_5->setSliderPosition(4);
        ui->eqSlider_6->setSliderPosition(4);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(2);
        ui->eqSlider_9->setSliderPosition(3);
        ui->eqSlider_10->setSliderPosition(4);
        applyEQ();
        break;
    case 9:
        ui->eqSlider_1->setSliderPosition(-2);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(2);
        ui->eqSlider_4->setSliderPosition(1);
        ui->eqSlider_5->setSliderPosition(0);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(-2);
        ui->eqSlider_10->setSliderPosition(-4);
        applyEQ();
        break;
    case 10:
        ui->eqSlider_1->setSliderPosition(-4);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(2);
        ui->eqSlider_4->setSliderPosition(1);
        ui->eqSlider_5->setSliderPosition(0);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(-4);
        ui->eqSlider_10->setSliderPosition(-6);
        applyEQ();
        break;
    case 11:
        ui->eqSlider_1->setSliderPosition(0);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(0);
        ui->eqSlider_4->setSliderPosition(4);
        ui->eqSlider_5->setSliderPosition(5);
        ui->eqSlider_6->setSliderPosition(3);
        ui->eqSlider_7->setSliderPosition(6);
        ui->eqSlider_8->setSliderPosition(3);
        ui->eqSlider_9->setSliderPosition(0);
        ui->eqSlider_10->setSliderPosition(0);
        applyEQ();
        break;
    case 12:
        ui->eqSlider_1->setSliderPosition(-4);
        ui->eqSlider_2->setSliderPosition(0);
        ui->eqSlider_3->setSliderPosition(2);
        ui->eqSlider_4->setSliderPosition(0);
        ui->eqSlider_5->setSliderPosition(0);
        ui->eqSlider_6->setSliderPosition(0);
        ui->eqSlider_7->setSliderPosition(0);
        ui->eqSlider_8->setSliderPosition(0);
        ui->eqSlider_9->setSliderPosition(-4);
        ui->eqSlider_10->setSliderPosition(-6);
        applyEQ();
        break;
    default:
        break;
    }
}

void ShadowPlayer::addToListAndPlay(QList<QUrl> files)
{
    int newIndex = playList->getLength();
    int length   = files.length();
    if (!files.isEmpty())
    {
        for (int i = 0; i < length; i++)
        {
            playList->add(files[i].toLocalFile());
        }
        if (playList->getLength() > newIndex)
            loadFile(playList->playIndex(newIndex));
    }
}

void ShadowPlayer::addToListAndPlay(QStringList files)
{
    int newIndex = playList->getLength();
    int length   = files.length();
    if (length > 0)
    {
        for (int i = 0; i < length; i++)
        {
            playList->add(files[i]);
        }
        if (playList->getLength() > newIndex)
            loadFile(playList->playIndex(newIndex));
    }
}

void ShadowPlayer::addToListAndPlay(QString file)
{
    int newIndex = playList->getLength();
    playList->add(file);
    if (playList->getLength() > newIndex)
        loadFile(playList->playIndex(newIndex));
}

void ShadowPlayer::on_playPreButton_clicked()
{
    if (playMode == 3)
    {
        loadFile(playList->previous(true));
    }
    else if (playMode == 4)
    {
        if (!playList->isEmpty())
        {
            int index = QRandomGenerator::global()->bounded(playList->getLength());
            loadFile(playList->playIndex(index));
        }
    }
    else
    {
        loadFile(playList->previous(false));
    }
}

void ShadowPlayer::on_playNextButton_clicked()
{
    if (playMode == 3)
    {
        loadFile(playList->next(true));
    }
    else if (playMode == 4)
    {
        if (!playList->isEmpty())
        {
            int index = QRandomGenerator::global()->bounded(playList->getLength());
            loadFile(playList->playIndex(index));
        }
    }
    else
    {
        loadFile(playList->next(false));
    }
}

void ShadowPlayer::on_playListButton_clicked()
{
    if (width() < 370)
        on_extraButton_clicked();

    if (playListHideAnimation->state() != QAbstractAnimation::Running && playListShowAnimation->state() != QAbstractAnimation::Running)
    {
        if (ui->playerListArea->height() < 190)
        {
            eqHideAnimation->stop();
            lyricsHideAnimation->stop();
            playListShowAnimation->stop();
            eqHideAnimation->start();
            lyricsHideAnimation->start();
            playListShowAnimation->start();
        }
        else
        {
            if (width() > 370)
            {
                eqShowAnimation->stop();
                lyricsShowAnimation->stop();
                playListHideAnimation->stop();
                eqShowAnimation->start();
                lyricsShowAnimation->start();
                playListHideAnimation->start();
            }
        }
    }
}

void ShadowPlayer::callFromPlayList()
{
    loadFile(playList->getCurFile());
}

void ShadowPlayer::on_reverseButton_clicked()
{
    if (isReverse)
    {
        isReverse = false;
        ui->reverseButton->setText(tr("Play"));
        player->setReverse(false);
    }
    else
    {
        isReverse = true;
        ui->reverseButton->setText(tr("Reverse Play"));
        player->setReverse(true);
    }
}

void ShadowPlayer::on_reverbDial_valueChanged(int value)
{
    player->updateReverb(value);
}

void ShadowPlayer::physicsSetting()
{
    double temp = 0;
    bool   ok   = false;

    temp = QInputDialog::getDouble(this,
                                   tr("Acceleration"),
                                   tr("Gravitational acceleration \n [parameters are ratios, 1 = total length of spectrum bar"),
                                   ui->FFTGroupBox->acc,
                                   0,
                                   2147483647,
                                   3,
                                   &ok);
    if (ok)
    {
        ui->FFTGroupBox->acc = temp;
    }

    temp = QInputDialog::getDouble(
        this, tr("Maximum drop speed"), tr("The maximum velocity of the object falling"), ui->FFTGroupBox->maxSpeed, 0, 2147483647, 3, &ok);
    if (ok)
    {
        ui->FFTGroupBox->maxSpeed = temp;
    }

    temp = QInputDialog::getDouble(this,
                                   tr("Base speed"),
                                   tr("Overall speed, force multiplier \n This value affects the elasticity, throwing force \n The initial version "
                                      "was used for falling speed, when there was no physical effect \n [modified with caution]"),
                                   ui->FFTGroupBox->speed,
                                   0,
                                   2147483647,
                                   3,
                                   &ok);
    if (ok)
    {
        ui->FFTGroupBox->speed = temp;
    }

    temp = QInputDialog::getDouble(this,
                                   tr("Throwing force multiplication factor"),
                                   tr("The coefficient of multiplication of the intensity of the throwing force of the spectrum bar on the object \n "
                                      "the throwing force will be doubled by this multiplier"),
                                   ui->FFTGroupBox->forceD,
                                   0,
                                   2147483647,
                                   3,
                                   &ok);
    if (ok)
    {
        ui->FFTGroupBox->forceD = temp;
    }

    temp = QInputDialog::getDouble(
        this,
        tr("Elasticity factor"),
        tr("The amount of kinetic energy retained after landing \n the object falls to the ground, part of the kinetic energy into sound and heat \n "
           "this parameter determines the percentage of kinetic energy remaining after the collision \n [Note: due to algorithm bugs, each fall will "
           "lose the potential energy of the last frame height from the ground, so kinetic energy has been lost] \n because it did not do air "
           "resistance, so it should make kinetic energy loss more \n1 = complete rebound"),
        ui->FFTGroupBox->elasticCoefficient,
        0,
        2147483647,
        3,
        &ok);
    if (ok)
    {
        ui->FFTGroupBox->elasticCoefficient = temp;
    }

    temp = QInputDialog::getDouble(this,
                                   tr("Elasticity Threshold"),
                                   tr("When the percentage of objects falling in one frame (considered to be penetrable) is less than this value \n "
                                      "does not calculate the elasticity"),
                                   ui->FFTGroupBox->minElasticStep,
                                   0,
                                   2147483647,
                                   3,
                                   &ok);
    if (ok)
    {
        ui->FFTGroupBox->minElasticStep = temp;
    }
}

void ShadowPlayer::enableFFTPhysics()
{
    ui->FFTGroupBox->acc                = 0.35;
    ui->FFTGroupBox->maxSpeed           = 9;
    ui->FFTGroupBox->speed              = 0.025;
    ui->FFTGroupBox->forceD             = 6;
    ui->FFTGroupBox->elasticCoefficient = 0.6;
    ui->FFTGroupBox->minElasticStep     = 0.02;
}

void ShadowPlayer::disableFFTPhysics()
{
    ui->FFTGroupBox->acc                = 0.15;
    ui->FFTGroupBox->maxSpeed           = 2;
    ui->FFTGroupBox->speed              = 0.025;
    ui->FFTGroupBox->forceD             = 0;
    ui->FFTGroupBox->elasticCoefficient = 0;
    ui->FFTGroupBox->minElasticStep     = 0.02;
}

void ShadowPlayer::resizeEvent(QResizeEvent *)
{
    ui->extraButton->setGeometry(width() - 20, 0, 20, 20);
    ui->closeButton->setGeometry(width() - 65, 0, 40, 20);
    ui->miniSizeButton->setGeometry(width() - 90, 0, 25, 20);
    ui->setSkinButton->setGeometry(width() - 115, 0, 25, 20);
}

void ShadowPlayer::showDeveloperInfo()
{
    QMessageBox::about(this, tr("Hannah"), tr("Simple Music Player"));
}

void ShadowPlayer::saveConfig()
{
    QFile file(QCoreApplication::applicationDirPath() + "/config.dat");
    file.open(QIODevice::WriteOnly);
    QDataStream stream(&file);
    stream << (quint32)0x61727480 << ui->freqSlider->value() << ui->volSlider->value() << isMute << ui->reverbDial->value() << playMode << skinMode
           << skinPos << skinDrawPos << ui->eqComboBox->currentIndex() << ui->eqEnableCheckBox->isChecked() << ui->eqSlider_1->value()
           << ui->eqSlider_2->value() << ui->eqSlider_3->value() << ui->eqSlider_4->value() << ui->eqSlider_5->value() << ui->eqSlider_6->value()
           << ui->eqSlider_7->value() << ui->eqSlider_8->value() << ui->eqSlider_9->value() << ui->eqSlider_10->value() << lb->isVisible();
    file.close();
}

void ShadowPlayer::loadConfig()
{
    QFile file(QCoreApplication::applicationDirPath() + "/config.dat");
    file.open(QIODevice::ReadOnly);
    QDataStream stream(&file);
    quint32     magic;
    stream >> magic;
    if (magic == 0x61727480)
    {
        int dataInt = 0;
        stream >> dataInt;
        ui->freqSlider->setValue(dataInt);
        stream >> dataInt;
        ui->volSlider->setValue(dataInt);
        bool dataBool = false;
        stream >> dataBool;
        isMute = dataBool;
        if (isMute)
        {
            ui->muteButton->setIcon(QIcon(":/rc/images/player/Mute.png"));
        }
        else
        {
            ui->muteButton->setIcon(QIcon(":/rc/images/player/Vol.png"));
        }
        stream >> dataInt;
        ui->reverbDial->setValue(dataInt);
        stream >> playMode;
        switch (playMode)
        {
        case 0:
            ui->playModeButton->setIcon(QIcon(":/rc/images/player/Single.png"));
            ui->playModeButton->setToolTip(tr("Track"));
            break;
        case 1:
            ui->playModeButton->setIcon(QIcon(":/rc/images/player/Repeat.png"));
            ui->playModeButton->setToolTip(tr("Track Repeat"));
            break;
        case 2:
            ui->playModeButton->setIcon(QIcon(":/rc/images/player/Order.png"));
            ui->playModeButton->setToolTip(tr("Order"));
            break;
        case 3:
            ui->playModeButton->setIcon(QIcon(":/rc/images/player/AllRepeat.png"));
            ui->playModeButton->setToolTip(tr("Playlist Repeat"));
            break;
        case 4:
            ui->playModeButton->setIcon(QIcon(":/rc/images/player/Shuffle.png"));
            ui->playModeButton->setToolTip(tr("Shuffle"));
            break;
        default:
            break;
        }
        stream >> skinMode;
        stream >> skinPos;
        stream >> skinDrawPos;
        stream >> dataInt;
        ui->eqComboBox->setCurrentIndex(dataInt);
        stream >> dataBool;
        ui->eqEnableCheckBox->setChecked(dataBool);
        stream >> dataInt;
        ui->eqSlider_1->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_2->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_3->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_4->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_5->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_6->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_7->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_8->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_9->setValue(dataInt);
        stream >> dataInt;
        ui->eqSlider_10->setValue(dataInt);
        stream >> dataBool;
        lb->setVisible(dataBool);
    }
    file.close();
}

void ShadowPlayer::saveSkinData()
{
    QFile file(QCoreApplication::applicationDirPath() + "/skin.dat");
    file.open(QIODevice::WriteOnly);
    QDataStream stream(&file);
    stream << (quint32)0x61727481 << skin;
    file.close();
}

void ShadowPlayer::loadSkinData()
{
    QFile file(QCoreApplication::applicationDirPath() + "/skin.dat");
    file.open(QIODevice::ReadOnly);
    QDataStream stream(&file);
    quint32     magic;
    stream >> magic;
    if (magic == 0x61727481)
    {
        stream >> skin;
        skinLeft    = skin.scaledToWidth(360, Qt::SmoothTransformation);
        skinFull    = skin.scaledToWidth(710, Qt::SmoothTransformation);
        aspectRatio = (double)skin.height() / skin.width();
    }
    file.close();
}

#if defined(Q_OS_WIN)
bool ShadowPlayer::nativeEvent(const QByteArray &eventType, void *message, long * /*result*/)
{
    Q_UNUSED(eventType);
    MSG *msg = reinterpret_cast<MSG *>(message);
    if (msg->message == WM_COPYDATA)
    {
        COPYDATASTRUCT *p = reinterpret_cast<COPYDATASTRUCT *>(msg->lParam);
        addToListAndPlay(QString::fromUtf8((LPCSTR)(p->lpData)));
        return true;
    }
    return false;
}

void ShadowPlayer::setTaskbarButtonWindow()
{
    taskbarButton->setWindow(windowHandle());
    thumbnailToolBar->setWindow(windowHandle());
}

#endif

void ShadowPlayer::on_eqEnableCheckBox_clicked(bool checked)
{
    if (checked)
    {
        player->eqReady();
        applyEQ();
    }
    else
    {
        player->disableEQ();
    }
}
