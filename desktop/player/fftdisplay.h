#ifndef FFTDISPLAY_H
#define FFTDISPLAY_H

#include <QGroupBox>

QT_FORWARD_DECLARE_CLASS(QPaintEvent);

class FFTDisplay : public QGroupBox
{
    Q_OBJECT
public:
    explicit FFTDisplay(QWidget *parent = 0);
    double fftBarValue[29];//存放频谱分析条的值
    void peakSlideDown();//下滑
    void cutValue();//限制最大最小值

    double acc, maxSpeed;//重力加速度，最大速度
    double speed;//基准速度。加速度为1时，按此速度下落/上升
    double forceD;//抛力系数
    double elasticCoefficient;//弹力系数
    double minElasticStep;//弹力禁用阀值

signals:

public slots:

protected:
    void paintEvent(QPaintEvent *event);

private:
    //数值除了绘图坐标，均为比值。
    int sx, sy, margin, width, height;//绘图相关，起始x，起始y，间距，宽度，高度
    double fftBarPeakValue[29];//存放频谱分析顶峰的值
    double peakSlideSpeed[29];//下滑速度
};

#endif // FFTDISPLAY_H
