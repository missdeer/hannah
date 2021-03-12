#include <QMouseEvent>

#include "spslider.h"

SPSlider::SPSlider(QWidget *parent) : QSlider(parent)
{
    //关闭分段移动
    setSingleStep(0);
    setPageStep(0);
}

//点击Slider即可调节Value
void SPSlider::mousePressEvent(QMouseEvent *event)
{
    if (event->buttons() & Qt::LeftButton)
    {
        QSlider::mousePressEvent(event);
        double pos = orientation() == Qt::Horizontal ? (event->pos().x() / (double)width()) : (event->pos().y() / (double)height());
        setValue(pos * (maximum() - minimum()) + minimum());
    }
}

void SPSlider::mouseMoveEvent(QMouseEvent *event)
{
    if (event->buttons() & Qt::LeftButton)
    {
        QSlider::mousePressEvent(event);
        double pos = orientation() == Qt::Horizontal ? (event->pos().x() / (double)width()) : (event->pos().y() / (double)height());
        setValue(pos * (maximum() - minimum()) + minimum());
    }
}
