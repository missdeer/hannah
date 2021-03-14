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
    Ui::LrcBar *    ui;
    QTimer *        timer;
    Lyrics *        lyrics;
    Player *        player;
    QPoint          pos;
    bool            clickOnFrame;
    bool            mouseEnter = {false};
    QLinearGradient linearGradient;
    QLinearGradient maskLinearGradient;
    QFont           font;
    int             shadowMode {0};

protected:
    void paintEvent(QPaintEvent *) override;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *) override;
    void enterEvent(QEvent *) override;
    void leaveEvent(QEvent *) override;
    void contextMenuEvent(QContextMenuEvent *event) override;
};

#endif // LRCBAR_H
