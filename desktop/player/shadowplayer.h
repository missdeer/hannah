#ifndef SHADOWPLAYER_H
#define SHADOWPLAYER_H
#include <math.h>

#include <QLinearGradient>
#include <QList>
#include <QMainWindow>
#include <QUrl>

#if defined(Q_OS_WIN)
#    include <QtWinExtras>
#endif
#include "bass.h"

QT_FORWARD_DECLARE_CLASS(QTimer);
QT_FORWARD_DECLARE_CLASS(QPropertyAnimation);

class Player;
class PlayList;
class Lyrics;
class OSD;
class LrcBar;

namespace Ui
{
    class ShadowPlayer;
}

class ShadowPlayer : public QMainWindow
{
    Q_OBJECT

public:
    explicit ShadowPlayer(QWidget *parent = 0);
    ~ShadowPlayer();
    void loadAudio(const QString &uri);
    void loadSkin(const QString &image, bool save = true);
    int  cutInt(int i);
    void addToListAndPlay(const QList<QUrl> &files);
    void addToListAndPlay(const QStringList &files);
    void addToListAndPlay(const QString &file);
    void showPlayer();
private slots:
    void UpdateTime();
    void UpdateLrc();
    void callFromPlayList();
    void on_openButton_clicked();
    void on_playButton_clicked();
    void on_stopButton_clicked();
    void on_volSlider_valueChanged(int value);
    void on_muteButton_clicked();
    void on_playSlider_sliderPressed();
    void on_playSlider_sliderReleased();
    void on_resetFreqButton_clicked();
    void applyEQ();
    void on_extraButton_clicked();
    void on_closeButton_clicked();
    void on_setSkinButton_clicked();
    void loadDefaultSkin();
    void fixSkinSizeLeft();
    void fixSkinSizeFull();
    void originalSkinSize();
    void autoSkinSize();
    void dynamicSkinSize();
    void skinOnTop();
    void skinOnCenter();
    void skinOnBottom();
    void skinDisable();
    void physicsSetting();
    void enableFFTPhysics();
    void disableFFTPhysics();
    void on_miniSizeButton_clicked();
    void on_playModeButton_clicked();
    void on_loadLrcButton_clicked();
    void on_playSlider_valueChanged(int value);
    void on_freqSlider_valueChanged(int value);
    void on_eqComboBox_currentIndexChanged(int index);
    void on_playPreButton_clicked();
    void on_playNextButton_clicked();
    void on_playListButton_clicked();
    void on_reverseButton_clicked();
    void showDeveloperInfo();
    void on_reverbDial_valueChanged(int value);

#if defined(Q_OS_WIN)
    void setTaskbarButtonWindow();
#endif

    void on_eqEnableCheckBox_clicked(bool checked);

public slots:
    void on_showDskLrcButton_clicked();

private:
    Ui::ShadowPlayer *ui;

    float arraySUM(int start, int end, float *array);
    void  updateFFT();
    void  showCoverPic(const QString &filePath);
    void  infoLabelAnimation();
    void  drawFFTBar(QWidget *parent, int x, int y, int width, int height, double percent);
    void  drawFFTBarPeak(QWidget *parent, int x, int y, int width, int height, double percent);
    void  saveConfig();
    void  loadConfig();
    void  saveSkinData();
    void  loadSkinData();

    QTimer *  timer;
    QTimer *  lrcTimer;
    Player *  player;
    Lyrics *  lyrics;
    OSD *     osd;
    LrcBar *  lb;
    PlayList *playList;
    bool      isPlaySliderPress {false};
    bool      isMute {false};
    int       lastVol;
    float     fftData[2048];
    double    fftBarValue[29];
    double    fftBarPeakValue[29];
    int       oriFreq;
    bool      playing {false};

    QPoint pos;
    bool   clickOnFrame {false};
    bool   clickOnLeft {false};

    QPixmap skin;
    QPixmap skinLeft;
    QPixmap skinFull;
    double  aspectRatio {0.0};
    int     skinMode {2};
    int     playMode {2};
    int     skinPos {1};
    double  skinDrawPos {0.0};
    bool    isReverse {false};

#if defined(Q_OS_WIN)
    QWinTaskbarButton *  taskbarButton;
    QWinTaskbarProgress *taskbarProgress;

    QWinThumbnailToolBar *   thumbnailToolBar;
    QWinThumbnailToolButton *playToolButton;
    QWinThumbnailToolButton *stopToolButton;
    QWinThumbnailToolButton *backwardToolButton;
    QWinThumbnailToolButton *forwardToolButton;
#endif

    QPropertyAnimation *sizeSlideAnimation;
    QPropertyAnimation *fadeInAnimation;
    QPropertyAnimation *tagAnimation;
    QPropertyAnimation *mediaInfoAnimation;
    QPropertyAnimation *coverAnimation;
    QPropertyAnimation *fadeOutAnimation;

    QPropertyAnimation *eqHideAnimation;
    QPropertyAnimation *eqShowAnimation;
    QPropertyAnimation *lyricsHideAnimation;
    QPropertyAnimation *lyricsShowAnimation;
    QPropertyAnimation *playListHideAnimation;
    QPropertyAnimation *playListShowAnimation;

    QLinearGradient bgLinearGradient;

protected:
    void dragEnterEvent(QDragEnterEvent *event) override;
    void dropEvent(QDropEvent *event) override;
    void paintEvent(QPaintEvent *) override;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void contextMenuEvent(QContextMenuEvent *) override;
    void closeEvent(QCloseEvent *) override;
    void resizeEvent(QResizeEvent *) override;
#if defined(Q_OS_WIN)
    bool nativeEvent(const QByteArray &eventType, void *message, long *result) override;
#endif
};

inline ShadowPlayer *shadowPlayer = nullptr;

#endif // SHADOWPLAYER_H
