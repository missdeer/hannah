#ifndef FFTDISPLAY_H
#define FFTDISPLAY_H

#include <QGroupBox>

QT_FORWARD_DECLARE_CLASS(QPaintEvent);

const int arraySize = 29;

class FFTDisplay : public QGroupBox
{
    Q_OBJECT
public:
    explicit FFTDisplay(QWidget *parent = 0);
    void   peakSlideDown();
    void   cutValue();

    double acc {0.35};
    double maxSpeed {9};
    double speed {0.025};
    double forceD {6};
    double elasticCoefficient {0.6};
    double minElasticStep {0.02};
    double fftBarValue[arraySize];

signals:

public slots:

protected:
    void paintEvent(QPaintEvent *event) override;

private:
    int    sx {10};
    int    sy {190};
    int    margin {1};
    int    width {10};
    int    height {120};
    double fftBarPeakValue[arraySize];
    double peakSlideSpeed[arraySize];
};

#endif // FFTDISPLAY_H
