#ifndef LRCBAR_H
#define LRCBAR_H

#include <QWidget>

QT_FORWARD_DECLARE_CLASS(QTimer);
QT_FORWARD_DECLARE_CLASS(QPaintEvent);
QT_FORWARD_DECLARE_CLASS(QMouseEvent);
QT_FORWARD_DECLARE_CLASS(QContextMenuEvent);

class Lyrics;
class Player;

namespace Ui
{
    class LrcBar;
}

class LrcBar : public QWidget
{
    Q_OBJECT

public:
    LrcBar(Lyrics *lrc, Player *plr, QWidget *parent = nullptr);
    ~LrcBar();

private slots:
    void UpdateTime();
    void settingFont();
    void enableShadow();
    void enableStroke();

private:
    Ui::LrcBar *ui;
    QTimer *timer;
    Lyrics *lyrics;
    Player *player;
    QPoint pos;//用于窗口拖动，存储鼠标坐标
    bool clickOnFrame;
    bool            mouseEnter = {false};
    QLinearGradient linearGradient;
    QLinearGradient maskLinearGradient;
    QFont font;//字体
    int             shadowMode {0}; //阴影模式'

protected:
    void paintEvent(QPaintEvent *);
    void mousePressEvent(QMouseEvent *event);//窗体拖动相关
    void mouseMoveEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *);
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);
    void contextMenuEvent(QContextMenuEvent *event);//右键菜单
};

#endif // LRCBAR_H
