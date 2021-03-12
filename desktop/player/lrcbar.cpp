#include <QApplication>
#include <QContextMenuEvent>
#include <QDesktopWidget>
#include <QFontDialog>
#include <QMenu>
#include <QMouseEvent>
#include <QPaintEvent>
#include <QPainter>
#include <QTimer>

#include "lrcbar.h"

#include "lyrics.h"
#include "player.h"
#include "ui_lrcbar.h"

LrcBar::LrcBar(Lyrics *lrc, Player *plr, QWidget *parent) : QWidget(parent), ui(new Ui::LrcBar), timer(new QTimer), lyrics(lrc), player(plr)
{
    ui->setupUi(this);

    font.setFamily("Microsoft Yahei");
    font.setPointSize(30);
    font.setStyleStrategy(QFont::PreferAntialias); //是否抗锯齿取决于字体

    //歌词渐变色彩
    linearGradient.setColorAt(0, QColor(0, 128, 0));
    linearGradient.setColorAt(1, QColor(0, 255, 0));
    //二层歌词渐变色彩
    maskLinearGradient.setColorAt(0, QColor(255, 255, 0));
    maskLinearGradient.setColorAt(0.5, QColor(255, 128, 0));
    maskLinearGradient.setColorAt(01, QColor(255, 255, 0));
    //计算渐变起点和终点
    QFontMetrics fm(font);
    linearGradient.setStart(0, height() / 2 - fm.height() / 2); //从字体顶部渐变到字体底部
    linearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);
    maskLinearGradient.setStart(0, height() / 2 - fm.height() / 2);
    maskLinearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);

    connect(timer, SIGNAL(timeout()), this, SLOT(UpdateTime()));
    timer->start(30);

    //设置起始位置
    setGeometry((QApplication::desktop()->screenGeometry().width() - width()) / 2,
                QApplication::desktop()->screenGeometry().height() - 130,
                width(),
                height());
    setFixedSize(width(), height());

    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowStaysOnTopHint | Qt::Tool);
    setAttribute(Qt::WA_TranslucentBackground); //背景透明
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

//窗口拖动相关
void LrcBar::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
    {
        pos = event->globalPos() - frameGeometry().topLeft();
        //也可以使用pos = event->pos();
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

    //以下内容为歌词文本绘制
    QString curLrc = lyrics->getLrcString(0);
    if (curLrc.isEmpty())
    {
        if (lyrics->isLrcEmpty())
        {
            curLrc = "Shadow Player";
        }
        else
        {
            curLrc = "（间奏……）";
        }
    }
    painter.setFont(font);
    painter.setRenderHint(QPainter::Antialiasing, true); //绘图抗锯齿

    QFontMetrics fm(font);                                                //字体参数，用于计算
    int          lrcWidth   = fm.width(curLrc);                           //歌词文本宽度
    double       curTimePos = lyrics->getTimePos(player->getCurTimeMS()); //当前时间点
    int          maskWidth  = lrcWidth * curTimePos;                      //计算出当前时间变色文本的宽度

    //原来的代码偶尔会有1像素点的误差，已删除，改用自己计算的中间位置
    //原先的实现方式是Qt::AlignCenter绘制固定文本，然后动态变色文本用自己计算的位置
    //某些特殊长度时，会出现1px的偏移

    if (fm.width(curLrc) < width())
    {
        //如果文本宽度没有超出窗口宽度
        int startXPos = width() / 2 - lrcWidth / 2; //歌词起始点（绘制在中央）

        //绘制底层文字作为阴影
        switch (shadowMode)
        {
        case 0:
            //右下双阴影
            painter.setPen(QColor(0, 0, 0, 80));
            painter.drawText(startXPos + 2, 2, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.setPen(QColor(0, 0, 0, 180));
            painter.drawText(startXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        case 1:
            //描边
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

        //绘制歌词文本
        painter.setPen(QPen(linearGradient, 0));
        painter.drawText(startXPos, 0, lrcWidth, height(), Qt::AlignVCenter, curLrc);

        //绘制变色文本
        painter.setPen(QPen(maskLinearGradient, 0));
        painter.drawText(startXPos, 0, maskWidth, height(), Qt::AlignVCenter, curLrc);
    }
    else
    {
        //如果文本宽度超出窗口宽度
        int rollXPos;

        //千千静听歌词滚动模式
        if (maskWidth < width() / 2)
        {
            //未达到中间宽度
            rollXPos = 0; //不移动
        }
        else if (lrcWidth - maskWidth < width() / 2)
        {
            //到达结尾部分
            rollXPos = width() - lrcWidth; //移动到最顶，不移动
        }
        else
        {
            //中间部分
            rollXPos = 0 - (maskWidth - width() / 2);
        }

        //这是原来的算法
        // int rollXPos = (width() - fm.width(curLrc)) * curTimePos;

        switch (shadowMode)
        {
        case 0:
            //右下双阴影
            painter.setPen(QColor(0, 0, 0, 80));
            painter.drawText(rollXPos + 2, 2, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.setPen(QColor(0, 0, 0, 180));
            painter.drawText(rollXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        case 1:
            //描边
            painter.setPen(QColor(0, 0, 0, 140));
            painter.drawText(rollXPos + 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(rollXPos - 1, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(rollXPos + 1, -1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            painter.drawText(rollXPos - 1, 1, lrcWidth, height(), Qt::AlignVCenter, curLrc);
            break;
        default:
            break;
        }

        //绘制歌词文本
        painter.setPen(QPen(linearGradient, 0));
        painter.drawText(rollXPos, 0, lrcWidth, height(), Qt::AlignVCenter, curLrc);

        //绘制变色文本
        painter.setPen(QPen(maskLinearGradient, 0));
        painter.drawText(rollXPos, 0, maskWidth, height(), Qt::AlignVCenter, curLrc);
    }
}

void LrcBar::enterEvent(QEvent *)
{
    mouseEnter = true;
    repaint();
}

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

//更改字体
void LrcBar::settingFont()
{
    bool  ok;
    QFont newFont = QFontDialog::getFont(&ok, font, this, tr("Select Font"));
    if (ok)
    {
        font = newFont;
    }
    //计算渐变起点和终点
    QFontMetrics fm(font);
    linearGradient.setStart(0, height() / 2 - fm.height() / 2); //从字体顶部渐变到字体底部
    linearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);
    maskLinearGradient.setStart(0, height() / 2 - fm.height() / 2);
    maskLinearGradient.setFinalStop(0, height() / 2 + fm.height() / 2);
}

void LrcBar::enableShadow()
{
    //阴影模式
    shadowMode = 0;
}

void LrcBar::enableStroke()
{
    //描边模式
    shadowMode = 1;
}
