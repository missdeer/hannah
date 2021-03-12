#include <QDesktopWidget>
#include <QMouseEvent>
#include <QPaintEvent>
#include <QPainter>
#include <QPropertyAnimation>
#include <QTimer>

#include "osd.h"

#include "ui_osd.h"

OSD::OSD(QWidget *parent) : QWidget(parent), ui(new Ui::OSD)
{
    ui->setupUi(this);
    timer = new QTimer(this);
    connect(timer, SIGNAL(timeout()), this, SLOT(timeRoll()));
    backGround = QPixmap(":/rc/images/OSD.png");

    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowStaysOnTopHint | Qt::Tool);
    setAttribute(Qt::WA_TranslucentBackground);

    setGeometry(QApplication::desktop()->screenGeometry().width() - width(), 80, width(), height());
    setFixedSize(width(), height());

    hideAnimation = new QPropertyAnimation(this, "windowOpacity");
    hideAnimation->setDuration(600);
    hideAnimation->setStartValue(1);
    hideAnimation->setEndValue(0);
    connect(hideAnimation, SIGNAL(finished()), this, SLOT(hide()));

    titleAnimation = new QPropertyAnimation(ui->titleLabel, "geometry");
    titleAnimation->setEasingCurve(QEasingCurve::OutCirc);
    titleAnimation->setDuration(1200);
    titleAnimation->setStartValue(QRect(30, 30, 0, 61));
    titleAnimation->setEndValue(QRect(40, 30, 281, 61));

    timeAnimation = new QPropertyAnimation(ui->timeLabel, "geometry");
    timeAnimation->setEasingCurve(QEasingCurve::OutCirc);
    timeAnimation->setDuration(1200);
    timeAnimation->setStartValue(QRect(331, 80, 0, 41));
    timeAnimation->setEndValue(QRect(40, 80, 281, 41));
}

OSD::~OSD()
{
    delete ui;
    delete timer;
    delete hideAnimation;
    delete titleAnimation;
    delete timeAnimation;
}

void OSD::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    painter.drawPixmap(0, 0, backGround);
}

void OSD::timeRoll()
{
    timeleft--;
    if (timeleft <= 0)
    {
        hideAnimation->start();
        // setVisible(false);
        timer->stop();
    }
}

void OSD::showOSD(QString tags, QString totalTime)
{
    timeleft = 4;
    hideAnimation->stop();
    titleAnimation->stop();
    titleAnimation->start();
    timeAnimation->stop();
    timeAnimation->start();
    setWindowOpacity(1);

    QFontMetrics titleFontMetrics(ui->titleLabel->font());
    if (titleFontMetrics.width(tags) > 281)
    {
        ui->titleLabel->setAlignment(Qt::AlignTop | Qt::AlignLeft);
        ui->titleLabel->setWordWrap(true);
    }
    else
    {
        ui->titleLabel->setAlignment(Qt::AlignCenter);
        ui->titleLabel->setWordWrap(false);
    }

    ui->titleLabel->setText(tags);
    ui->timeLabel->setText(totalTime);
    timer->start(1000);
    show();
    update();
}

void OSD::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
    {
        timer->stop();
        hide();
    }
}
