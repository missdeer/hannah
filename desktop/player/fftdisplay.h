#ifndef FFTDISPLAY_H
#define FFTDISPLAY_H

#include <QGroupBox>

QT_FORWARD_DECLARE_CLASS(QPaintEvent);

class FFTDisplay : public QGroupBox
{
    Q_OBJECT
public:
    explicit FFTDisplay(QWidget *parent = 0);
    double fftBarValue[29];
    void   peakSlideDown();
    void   cutValue();

    double acc, maxSpeed;
    double speed;
    double forceD;
    double elasticCoefficient;
    double minElasticStep;

signals:

public slots:

protected:
    void paintEvent(QPaintEvent *event) override;

private:
    int    sx, sy, margin, width, height;
    double fftBarPeakValue[29];
    double peakSlideSpeed[29];
};

#endif // FFTDISPLAY_H
