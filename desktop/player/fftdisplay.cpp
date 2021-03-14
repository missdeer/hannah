#include <QBrush>
#include <QColor>
#include <QPainter>
#include <QPen>

#include "fftdisplay.h"

FFTDisplay::FFTDisplay(QWidget *parent) : QGroupBox(parent)
{
    sx                 = 10;
    sy                 = 190;
    margin             = 1;
    width              = 10;
    height             = 120;
    acc                = 0.35;
    maxSpeed           = 9;
    speed              = 0.025;
    forceD             = 6;
    elasticCoefficient = 0.6;
    minElasticStep     = 0.02;
    setAttribute(Qt::WA_TransparentForMouseEvents, true);

    for (int i = 0; i < 29; i++)
    {
        fftBarValue[i]     = 0;
        fftBarPeakValue[i] = 0;
        peakSlideSpeed[i]  = 0;
    }
}

void FFTDisplay::paintEvent(QPaintEvent * /*event*/)
{
    QPainter painter(this);

    QPen noPen(Qt::NoPen);
    QPen peakPen(QColor(110, 139, 61));
    peakPen.setWidth(2);
    QBrush stickBrush(QColor(70, 130, 180, 180));

    for (int i = 0; i < 29; i++)
    {
        int drawX = sx + (i * (width + margin));

        stickBrush.setColor(QColor(70, 130, 180, (int)(140 * fftBarValue[i]) + 90));

        painter.setPen(noPen);
        painter.setBrush(stickBrush);
        painter.drawRect(drawX, sy + (height * (1 - fftBarValue[i])), width, (height - 0.000000001) * fftBarValue[i]);

        painter.setPen(peakPen);
        int y = sy + ((double)height * (1 - fftBarPeakValue[i])) - 1;
        painter.drawLine(drawX + 1, y, drawX + width - 1, y);
    }
}

void FFTDisplay::peakSlideDown()
{
    for (int i = 0; i < 29; i++)
    {
        double realityDown = speed * peakSlideSpeed[i];
        if (fftBarPeakValue[i] - realityDown <= fftBarValue[i])
        {
            double elasticForce = 0;
            if (realityDown > minElasticStep)
                elasticForce = peakSlideSpeed[i] * elasticCoefficient;

            peakSlideSpeed[i]  = (fftBarPeakValue[i] - fftBarValue[i]) * forceD - elasticForce;
            fftBarPeakValue[i] = fftBarValue[i];
        }
        else
        {
            //下落
            if (fftBarPeakValue[i] > realityDown)
            {
                fftBarPeakValue[i] -= realityDown;
            }
            else
            {
                fftBarPeakValue[i] = 0;
            }

            if (peakSlideSpeed[i] + acc <= maxSpeed)
                peakSlideSpeed[i] += acc;
            else
                peakSlideSpeed[i] = maxSpeed;
        }
    }
}

void FFTDisplay::cutValue()
{
    for (int i = 0; i < 29; i++)
    {
        if (fftBarValue[i] > 1)
            fftBarValue[i] = 1;
        else if (fftBarValue[i] < 0)
            fftBarValue[i] = 0;
    }
}
