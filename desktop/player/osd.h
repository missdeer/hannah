#ifndef OSD_H
#define OSD_H

#include <QPixmap>
#include <QWidget>

QT_FORWARD_DECLARE_CLASS(QTimer);
QT_FORWARD_DECLARE_CLASS(QPropertyAnimation);
QT_FORWARD_DECLARE_CLASS(QPaintEvent);
QT_FORWARD_DECLARE_CLASS(QMouseEvent);

namespace Ui {
class OSD;
}

class OSD : public QWidget
{
    Q_OBJECT

public:
    explicit OSD(QWidget *parent = 0);
    ~OSD();
    void showOSD(QString tags, QString totalTime);

private:
    Ui::OSD *ui;
    QPixmap backGround;
    QTimer *timer;
    int timeleft;
    QPropertyAnimation *hideAnimation;
    QPropertyAnimation *titleAnimation;
    QPropertyAnimation *timeAnimation;

private slots:
    void timeRoll();

protected:
    void paintEvent(QPaintEvent *);
    void mousePressEvent(QMouseEvent *event);
};

#endif // OSD_H
