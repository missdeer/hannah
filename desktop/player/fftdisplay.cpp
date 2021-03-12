#include <QBrush>
#include <QColor>
#include <QPainter>
#include <QPen>

#include "fftdisplay.h"

FFTDisplay::FFTDisplay(QWidget *parent) : QGroupBox(parent)
{
    //相关参数
    // sx=起始x，sy=起始y
    sx                 = 10;
    sy                 = 190;
    margin             = 1;
    width              = 10;
    height             = 120;
    acc                = 0.35;
    maxSpeed           = 9;                               //加速度，最大速度
    speed              = 0.025;                           //基准速度
    forceD             = 6;                               //抛力系数
    elasticCoefficient = 0.6;                             //弹力系数
    minElasticStep     = 0.02;                            //弹力禁用阀值
    setAttribute(Qt::WA_TransparentForMouseEvents, true); //鼠标事件交给父窗口

    //初始化数组
    for (int i = 0; i < 29; i++)
    {
        fftBarValue[i]     = 0;
        fftBarPeakValue[i] = 0;
        peakSlideSpeed[i]  = 0;
    }
}

void FFTDisplay::paintEvent(QPaintEvent * /*event*/)
{
    //绘制组合框
    // QGroupBox::paintEvent(event);

    //绘制添加的内容
    QPainter painter(this);

    QPen noPen(Qt::NoPen);
    QPen peakPen(QColor(110, 139, 61));
    peakPen.setWidth(2);
    QBrush stickBrush(QColor(70, 130, 180, 180));

    for (int i = 0; i < 29; i++)
    {
        int drawX = sx + (i * (width + margin)); //计算x坐标

        stickBrush.setColor(QColor(70, 130, 180, (int)(140 * fftBarValue[i]) + 90));

        //绘制条
        painter.setPen(noPen);
        painter.setBrush(stickBrush);
        // y轴正方向向下，起点坐标是矩形的左上角，高度向下。
        // height * (1 - fftBarValue[i])：顶部留白的高度。
        painter.drawRect(drawX, sy + (height * (1 - fftBarValue[i])), width, (height - 0.000000001) * fftBarValue[i]);

        //绘制下滑顶峰
        painter.setPen(peakPen);
        //计算实际坐标
        int y = sy + ((double)height * (1 - fftBarPeakValue[i])) - 1; //减去1确保位置在条形上方
        painter.drawLine(drawX + 1, y, drawX + width - 1, y);
    }
}

//下滑（计算“面饼”下一帧的位置，包含物理运算）
void FFTDisplay::peakSlideDown()
{
    //计算峰值
    for (int i = 0; i < 29; i++)
    {
        double realityDown = speed * peakSlideSpeed[i]; //计算下落速度
        if (fftBarPeakValue[i] - realityDown <= fftBarValue[i])
        {
            // fftBarPeakValue[i] - realityDown指的是下落之后的位置。
            //如果下落后将会比矩形低，则落在矩形上。
            //撞击
            double elasticForce = 0;
            if (realityDown > minElasticStep)
                elasticForce = peakSlideSpeed[i] * elasticCoefficient; //计算弹力

            peakSlideSpeed[i]  = (fftBarPeakValue[i] - fftBarValue[i]) * forceD - elasticForce; //计算撞击+弹力产生的加速度
            fftBarPeakValue[i] = fftBarValue[i]; //此次存在误差，下落过程产生的动量被忽略，导致动量损失
        }
        else
        {
            //下落
            if (fftBarPeakValue[i] > realityDown) //下落的情况，当前值是否比每次下落量还小
            {
                fftBarPeakValue[i] -= realityDown; //如果不是，可以下落（不会小于0）
            }
            else
            {
                fftBarPeakValue[i] = 0; //下落后将变为负值，不再下落（赋值0）
            }

            //下落过程加速
            if (peakSlideSpeed[i] + acc <= maxSpeed)
                peakSlideSpeed[i] += acc; //加速下落
            else
                peakSlideSpeed[i] = maxSpeed; //达到最大速度
        }
    }
}

//频谱条 限值
void FFTDisplay::cutValue()
{
    for (int i = 0; i < 29; i++)
    {
        //限制条取值范围
        if (fftBarValue[i] > 1)
            fftBarValue[i] = 1;
        else if (fftBarValue[i] < 0)
            fftBarValue[i] = 0;
    }
}
