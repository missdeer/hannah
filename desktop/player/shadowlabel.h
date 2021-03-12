#ifndef SHADOWLABEL_H
#define SHADOWLABEL_H

#include <QLabel>

QT_FORWARD_DECLARE_CLASS(QPaintEvent);
QT_FORWARD_DECLARE_CLASS(QColor);

class ShadowLabel : public QLabel
{
    Q_OBJECT
public:
    explicit ShadowLabel(QWidget *parent = 0);
    void setShadowColor(QColor color);
    void setShadowMode(int mode);

signals:

public slots:

private:
    QColor shadowColor {QColor(255, 255, 255, 128)};
    int    shadowMode {0};

protected:
    void paintEvent(QPaintEvent *event);

};

#endif // SHADOWLABEL_H
