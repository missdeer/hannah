#include <QContextMenuEvent>
#include <QFontDialog>
#include <QGuiApplication>
#include <QMenu>
#include <QMouseEvent>
#include <QPaintEvent>
#include <QPainter>
#include <QScreen>
#include <QTimer>

#include "lrcbar.h"
#include "bassplayer.h"
#include "lyrics.h"
#include "ui_lrcbar.h"

LrcBar::LrcBar(Lyrics *lrc, BassPlayer *plr, QWidget *parent) : QWidget(parent), ui(new Ui::LrcBar), timer(new QTimer), lyrics(lrc), player(plr)
{
    ui->setupUi(this);

    font.setFamily("Microsoft Yahei");
    font.setPointSize(30);
    font.setStyleStrategy(QFont::PreferAntialias);
    linearGradient.setColorAt(0, QColor(0, 128, 0));
    linearGradient.setColorAt(1, QColor(0, 255, 0));
    maskLinearGradient.setColorAt(0, QColor(255, 255, 0));
    maskLinearGradient.setColorAt(0.5, QColor(255, 128, 0));
    maskLinearGradient.setColorAt(01, QColor(255, 255, 0));
    QFontMetrics fm(font);
    linearGradient.setStart(0, height() / 2 - fm.height() / 2);
    linearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);
    maskLinearGradient.setStart(0, height() / 2 - fm.height() / 2);
    maskLinearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);

    connect(timer, SIGNAL(timeout()), this, SLOT(UpdateTime()));
    timer->start(30);

    setGeometry((QGuiApplication::primaryScreen()->availableGeometry().width() - width()) / 2,
                QGuiApplication::primaryScreen()->availableGeometry().height() - 130,
                width(),
                height());
    setFixedSize(width(), height());

    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowStaysOnTopHint | Qt::Tool);
    setAttribute(Qt::WA_TranslucentBackground);
}

void LrcBar::UpdateTime()
{
    if (isVisible())
    {
        lyrics->updateTime(player->getCurTimeMS(), player->getTotalTimeMS());
        ;
        repaint();
    }
}

LrcBar::~LrcBar()
{
    delete ui;
}

void LrcBar::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
    {
        pos = event->globalPos() - frameGeometry().topLeft();

        clickOnFrame = true;
        event->accept();
    }
}

void LrcBar::mouseReleaseEvent(QMouseEvent *)
{
    clickOnFrame = false;
}

void LrcBar::mouseMoveEvent(QMouseEvent *event)
{
    if (event->buttons() & Qt::LeftButton && clickOnFrame)
    {
        move(event->globalPos() - pos);
        event->accept();
    }
}

void LrcBar::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    if (mouseEnter)
    {
        painter.setBrush(QBrush(QColor(255, 255, 255, 120)));
        painter.setPen(Qt::NoPen);
        painter.drawRect(0, 0, width(), height());
    }

    QString curLrc = lyrics->getLrcString(0);
    if (curLrc.isEmpty())
    {
        if (lyrics->isLrcEmpty())
        {
            curLrc = tr("Hannah");
        }
        else
        {
            curLrc = tr("Interlude...");
        }
    }
    painter.setFont(font);
    painter.setRenderHint(QPainter::Antialiasing, true);

    QFontMetrics fm(font);
    int          lrcWidth   = fm.boundingRect(curLrc).width();
    double       curTimePos = lyrics->getTimePos(player->getCurTimeMS());
    int          maskWidth  = lrcWidth * curTimePos;

    if (fm.boundingRect(curLrc).width() < width())
    {
        int startXPos = width() / 2 - lrcWidth / 2;
        switch (shadowMode)
        {
        case 0:
            painter.setPen(QColor(0, 0, 0, 80));
            painter.drawText(startXPos + 2, 2, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.setPen(QColor(0, 0, 0, 180));
            painter.drawText(startXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        case 1:
            painter.setPen(QColor(0, 0, 0, 140));
            painter.drawText(startXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos + 1, 0, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos + 1, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos - 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos - 1, 0, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(startXPos - 1, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        default:
            break;
        }

        painter.setPen(QPen(linearGradient, 0));
        painter.drawText(startXPos, 0, lrcWidth, height(), Qt::AlignVCenter, curLrc);

        painter.setPen(QPen(maskLinearGradient, 0));
        painter.drawText(startXPos, 0, maskWidth, height(), Qt::AlignVCenter, curLrc);
    }
    else
    {
        int rollXPos;

        if (maskWidth < width() / 2)
        {
            rollXPos = 0;
        }
        else if (lrcWidth - maskWidth < width() / 2)
        {
            rollXPos = width() - lrcWidth;
        }
        else
        {
            rollXPos = 0 - (maskWidth - width() / 2);
        }

        switch (shadowMode)
        {
        case 0:
            painter.setPen(QColor(0, 0, 0, 80));
            painter.drawText(rollXPos + 2, 2, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.setPen(QColor(0, 0, 0, 180));
            painter.drawText(rollXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        case 1:
            painter.setPen(QColor(0, 0, 0, 140));
            painter.drawText(rollXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(rollXPos - 1, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(rollXPos + 1, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(rollXPos - 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        default:
            break;
        }

        painter.setPen(QPen(linearGradient, 0));
        painter.drawText(rollXPos, 0, lrcWidth, height(), Qt::AlignVCenter, curLrc);

        painter.setPen(QPen(maskLinearGradient, 0));
        painter.drawText(rollXPos, 0, maskWidth, height(), Qt::AlignVCenter, curLrc);
    }
}

#if QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
void LrcBar::enterEvent(QEvent *)
{
    mouseEnter = true;
    repaint();
}
#endif

void LrcBar::leaveEvent(QEvent *)
{
    mouseEnter = false;
    repaint();
}

void LrcBar::contextMenuEvent(QContextMenuEvent *event)
{
    QMenu menu;
    menu.addAction(tr("Font Settings"), this, SLOT(settingFont()));
    menu.addAction(tr("Shadow Mode(TTPlayer)"), this, SLOT(enableShadow()));
    menu.addAction(tr("Stroke Mode(Kugou Music)"), this, SLOT(enableStroke()));
    menu.exec(event->globalPos());
}

void LrcBar::settingFont()
{
    bool  ok;
    QFont newFont = QFontDialog::getFont(&ok, font, this, tr("Select Font"));
    if (ok)
    {
        font = newFont;
    }
    QFontMetrics fm(font);
    linearGradient.setStart(0, height() / 2 - fm.height() / 2);
    linearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);
    maskLinearGradient.setStart(0, height() / 2 - fm.height() / 2);
    maskLinearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);
}

void LrcBar::enableShadow()
{
    shadowMode = 0;
}

void LrcBar::enableStroke()
{
    shadowMode = 1;
}
