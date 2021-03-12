#include <QPaintEvent>
#include <QPainter>

#include "shadowlabel.h"

ShadowLabel::ShadowLabel(QWidget *parent) : QLabel(parent) {}

void ShadowLabel::paintEvent(QPaintEvent * /*event*/)
{
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setFont(font());
    painter.setPen(shadowColor);

    switch (shadowMode)
    {
    case 0:
        painter.drawText(2, 2, width() - 1, height() - 1, alignment(), text());
        break;
    case 1:
        painter.drawText(0, 0, width() - 1, height() - 1, alignment(), text());
        painter.drawText(0, 1, width() - 1, height() - 1, alignment(), text());
        painter.drawText(0, 2, width() - 1, height() - 1, alignment(), text());
        painter.drawText(1, 0, width() - 1, height() - 1, alignment(), text());
        painter.drawText(1, 2, width() - 1, height() - 1, alignment(), text());
        painter.drawText(2, 0, width() - 1, height() - 1, alignment(), text());
        painter.drawText(2, 1, width() - 1, height() - 1, alignment(), text());
        painter.drawText(2, 2, width() - 1, height() - 1, alignment(), text());
        break;
    default:
        painter.drawText(1, 1, width(), height(), alignment(), text());
        break;
    }
    painter.setPen(palette().color(QPalette::WindowText));
    painter.drawText(1, 1, width() - 1, height() - 1, alignment(), text());
}

void ShadowLabel::setShadowColor(QColor color)
{
    shadowColor = color;
}

void ShadowLabel::setShadowMode(int mode)
{
    shadowMode = mode;
}
