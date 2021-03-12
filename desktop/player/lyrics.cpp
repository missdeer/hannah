#include <QCoreApplication>
#include <QFileInfo>
#include <QTextStream>

#include "lyrics.h"
//歌词解析，传入参数：音频文件名
bool Lyrics::resolve(const QString &fileName, bool isLrc)
{
    //解决乱码：http://blog.csdn.net/locky1218/article/details/10568261
    // Lrc文件读入QMap：http://www.cnblogs.com/tornadomeet/archive/2012/09/23/2699077.html
    //解析方式和我原想的方式非常接近……
    curLrcTime  = 0;
    nextLrcTime = 0;
    lrcMap.clear();
    timeList.clear();
    if (fileName.isEmpty())
        return false; //文件名不正确

    QFileInfo fileInfo(fileName);
    QString   lrcFileName;

    if (!isLrc)
    {
        lrcFileName = fileInfo.path() + "/" + fileInfo.completeBaseName() + ".lrc"; //打开同名lrc
    }
    else
    {
        lrcFileName = fileName;
    }

    QFile file(lrcFileName);

    if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
    {
        return false; //打开lrc文件失败（歌词不存在，歌词文件被独占）
    }

    QTextStream textIn(&file);                       //使用QTextStream读取文本，即可解决ANSI编码问题
    QString     allText = QString(textIn.readAll()); //全部读取
    file.close();
    QStringList lines = allText.split("\n"); //按行分割文本;

    QRegExp rx("\\[\\d{2}:\\d{2}\\.\\d{2}\\]"); //正则表达式匹配歌词时间（不支持[xxx:xx.xx]以上

    foreach (QString oneline, lines)
    {
        QString text = oneline;

        text.replace(rx, "");             //删除正则表达式匹配部分（时间标签），此时text内容为歌词文本
        int pos = rx.indexIn(oneline, 0); //返回匹配位置，-1表示匹配失败

        //分段读取，写入Map
        while (pos != -1)
        {
            QString cap = rx.cap(0); //返回第0个表达式匹配的内容
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
            pos = rx.indexIn(oneline, pos); //更改位置继续匹配
        }
    }

    if (!lrcMap.isEmpty())
    {
        timeList = lrcMap.uniqueKeys(); //得到时间列表以便搜索
        return true;
    }
    return false;
}

//载入程序目录下歌词文件夹中的歌词，文件存在返回true
bool Lyrics::loadFromLrcDir(const QString &fileName)
{
    QFileInfo fileInfo(fileName);
    QString   fn = QCoreApplication::applicationDirPath() + "/lyrics/" + fileInfo.completeBaseName() + ".lrc";

    if (QFile::exists(fn))
    {
        resolve(fn, true);
        return true; //文件存在
    }
    else
    {
        return false; //文件不存在
    }
}

bool Lyrics::loadFromFileRelativePath(const QString &fileName, const QString &path)
{
    QFileInfo fileInfo(fileName);
    QString   fn = fileInfo.path() + path + fileInfo.completeBaseName() + ".lrc";

    if (QFile::exists(fn))
    {
        resolve(fn, true);
        return true; //文件存在
    }
    else
    {
        return false; //文件不存在
    }
}

//更新歌词显示时间
void Lyrics::updateTime(int curms, int totalms)
{
    if (!lrcMap.isEmpty())
    {
        int time     = 0;
        int nextTime = 0;
        // keys()方法返回lrcMap列表
        foreach (int value, lrcMap.keys())
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
        curLrcTime = time; //记录要显示的歌词所在的时间

        if (nextTime != 0)
            nextLrcTime = nextTime; //如果有下一句歌词的话，设置时间
        else
            nextLrcTime = totalms; //否则设置下一句为总时间
    }
}

//取得curLrcTime下，指定偏移量的歌词
QString Lyrics::getLrcString(int offset)
{
    if (!lrcMap.isEmpty())
    {
        int showTime = 0;

        int index = timeList.indexOf(curLrcTime); //取得当前歌词索引
        if (index + offset >= 0 && index + offset < timeList.size())
        {
            showTime = timeList[index + offset]; //取得偏移后，要显示的歌词时间
        }
        else
        {
            return ""; //如果没有此索引，返回空字符串
        }

        return lrcMap.value(showTime); //返回要显示的歌词
    }
    else
    {
        return ""; //如果没有歌词，返回空字符串
    }
    return ""; //空Map返回空字符串
}

//返回输入时间在当前句子中的百分比，取值0~1
double Lyrics::getTimePos(int ms)
{
    if (!lrcMap.isEmpty())
    {
        if (ms < curLrcTime)
        {
            return 0;
        }
        else if (ms > nextLrcTime)
        {
            return 1;
        }
        else
        {
            return (double)(ms - curLrcTime) / (nextLrcTime - curLrcTime);
        }
    }
    else
    {
        return 0; //使“Shadow Player”文本未填充
    }
}

bool Lyrics::isLrcEmpty()
{
    return lrcMap.isEmpty();
}
