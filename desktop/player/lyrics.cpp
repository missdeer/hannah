#include <QCoreApplication>
#include <QFileInfo>
#include <QTextStream>

#include "lyrics.h"

bool Lyrics::resolve(const QString &fileName, bool isLrc)
{
    if (fileName.isEmpty())
        return false;

    curLrcTime  = 0;
    nextLrcTime = 0;
    lrcMap.clear();
    timeList.clear();
    QFileInfo fileInfo(fileName);
    QString   lrcFileName;

    if (!isLrc)
    {
        lrcFileName = fileInfo.path() + "/" + fileInfo.completeBaseName() + ".lrc";
    }
    else
    {
        lrcFileName = fileName;
    }

    QFile file(lrcFileName);

    if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
    {
        return false;
    }

    QTextStream textIn(&file);
    QString     allText = textIn.readAll();
    file.close();
    QStringList lines = allText.split("\n");

    QRegExp rx("\\[\\d{2}:\\d{2}\\.\\d{2}\\]");

    for (const auto &oneline : lines)
    {
        QString text = oneline;

        text.replace(rx, "");
        int pos = rx.indexIn(oneline, 0);

        while (pos != -1)
        {
            QString cap = rx.cap(0);
            QRegExp rx2;
            rx2.setPattern("\\d{2}(?=:)");
            rx2.indexIn(cap);
            int minute = rx2.cap(0).toInt();
            rx2.setPattern("\\d{2}(?=\\.)");
            rx2.indexIn(cap);
            int second = rx2.cap(0).toInt();
            rx2.setPattern("\\d{2}(?=\\])");
            rx2.indexIn(cap);
            int millisecond = rx2.cap(0).toInt();
            int totalTime   = minute * 60000 + second * 1000 + millisecond * 10;
            lrcMap.insert(totalTime, text);
            pos += rx.matchedLength();
            pos = rx.indexIn(oneline, pos);
        }
    }

    if (!lrcMap.isEmpty())
    {
        timeList = lrcMap.keys();
        return true;
    }
    return false;
}

bool Lyrics::loadFromLrcDir(const QString &fileName)
{
    QFileInfo fileInfo(fileName);
    QString   fn = QCoreApplication::applicationDirPath() + "/lyrics/" + fileInfo.completeBaseName() + ".lrc";

    if (QFile::exists(fn))
    {
        resolve(fn, true);
        return true;
    }
    return false;
}

bool Lyrics::loadFromFileRelativePath(const QString &fileName, const QString &path)
{
    QFileInfo fileInfo(fileName);
    QString   fn = fileInfo.path() + path + fileInfo.completeBaseName() + ".lrc";

    if (QFile::exists(fn))
    {
        resolve(fn, true);
        return true;
    }
    return false;
}

void Lyrics::updateTime(int curms, int totalms)
{
    if (!lrcMap.isEmpty())
    {
        int time     = 0;
        int nextTime = 0;

        auto keys = lrcMap.keys();
        for (int value : keys)
        {
            if (curms >= value)
            {
                time = value;
            }
            else
            {
                nextTime = value;
                break;
            }
        }
        curLrcTime = time;

        if (nextTime != 0)
            nextLrcTime = nextTime;
        else
            nextLrcTime = totalms;
    }
}

QString Lyrics::getLrcString(int offset)
{
    if (!lrcMap.isEmpty())
    {
        int showTime = 0;

        int index = timeList.indexOf(curLrcTime);
        if (index + offset >= 0 && index + offset < timeList.size())
        {
            showTime = timeList[index + offset];
            return lrcMap.value(showTime);
        }
    }
    return "";
}

double Lyrics::getTimePos(int ms)
{
    if (!lrcMap.isEmpty())
    {
        if (ms < curLrcTime)
        {
            return 0;
        }
        if (ms > nextLrcTime)
        {
            return 1;
        }

        return (double)(ms - curLrcTime) / (nextLrcTime - curLrcTime);
    }
    return 0;
}

bool Lyrics::isLrcEmpty()
{
    return lrcMap.isEmpty();
}
