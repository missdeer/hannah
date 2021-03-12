/*
 * 重写的播放进度条，点击后立即跳转到指定位置
 */
#ifndef SPSLIDER_H
#define SPSLIDER_H

#include <QSlider>

QT_FORWARD_DECLARE_CLASS(QMouseEvent);

class SPSlider : public QSlider
{
    Q_OBJECT
public:
    explicit SPSlider(QWidget *parent = 0);

signals:

public slots:

protected:
    void mousePressEvent(QMouseEvent *event);
    void mouseMoveEvent(QMouseEvent *event);
};

#endif // SPSLIDER_H
