#include <cstring>

#include <QBrush>
#include <QColor>
#include <QPainter>
#include <QPen>

#include "fftdisplay.h"

FFTDisplay::FFTDisplay(QWidget *parent) : QGroupBox(parent)
{
    height             = 120;
    acc                = 0.35;
    maxSpeed           = 9;
    speed              = 0.025;
    forceD             = 6;
    elasticCoefficient = 0.6;
    minElasticStep     = 0.02;
    setAttribute(Qt::WA_TransparentForMouseEvents, true);

    std::memset(fftBarValue, 0, sizeof(fftBarValue));
    std::memset(fftBarPeakValue, 0, sizeof(fftBarPeakValue));
    std::memset(peakSlideSpeed, 0, sizeof(peakSlideSpeed));
}

void FFTDisplay::paintEvent(QPaintEvent * /*event*/)
{
    QPainter painter(this);

    QPen noPen(Qt::NoPen);
    QPen peakPen(QColor(110, 139, 61));
    peakPen.setWidth(2);
    QBrush stickBrush(QColor(70, 130, 180, 180));

    for (int i = 0; i < arraySize; i++)
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
    for (int i = 0; i < arraySize; i++)
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
    for (double &i : fftBarValue)
    {
        if (i > 1)
        {
            i = 1;
        }
        else if (i < 0)
        {
            i = 0;
        }
    }
}
