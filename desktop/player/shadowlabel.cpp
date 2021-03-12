#include <QPaintEvent>
#include <QPainter>

#include "shadowlabel.h"

ShadowLabel::ShadowLabel(QWidget *parent) : QLabel(parent) {}

void ShadowLabel::paintEvent(QPaintEvent * /*event*/)
{
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true); //绘图抗锯齿
    painter.setFont(font());
    painter.setPen(shadowColor); //取得阴影颜色
    //绘制阴影
    switch (shadowMode)
    {
    case 0:
        painter.drawText(2, 2, width() - 1, height() - 1, alignment(), text());
        break;
    case 1:
        painter.drawText(0, 0, width() - 1, height() - 1, alignment(), text()); //左上
        painter.drawText(0, 1, width() - 1, height() - 1, alignment(), text()); //左
        painter.drawText(0, 2, width() - 1, height() - 1, alignment(), text()); //左下
        painter.drawText(1, 0, width() - 1, height() - 1, alignment(), text()); //上
        painter.drawText(1, 2, width() - 1, height() - 1, alignment(), text()); //下
        painter.drawText(2, 0, width() - 1, height() - 1, alignment(), text()); //右上
        painter.drawText(2, 1, width() - 1, height() - 1, alignment(), text()); //右
        painter.drawText(2, 2, width() - 1, height() - 1, alignment(), text()); //右下
        break;
    default:
        painter.drawText(1, 1, width(), height(), alignment(), text());
        break;
    }
    painter.setPen(palette().color(QPalette::WindowText));                  //取得文本颜色
    painter.drawText(1, 1, width() - 1, height() - 1, alignment(), text()); //绘制文本
}

void ShadowLabel::setShadowColor(QColor color)
{
    shadowColor = color;
}

void ShadowLabel::setShadowMode(int mode)
{
    shadowMode = mode;
}
