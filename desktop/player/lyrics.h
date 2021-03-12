/*
 * 歌词解析、读取类
 * 支持乱序歌词、多时间标签歌词
 * 注意：<*不支持100分钟以上歌词*>
 */
#ifndef LYRICS_H
#define LYRICS_H
#include <QList>
#include <QMap>
#include <QString>

class Lyrics
{
public:
    Lyrics();
    bool    resolve(const QString &fileName, bool isLrc = false); //歌词解析 ,isLrc：此文件是否为歌词文件，失败返回false
    bool    loadFromLrcDir(const QString &fileName);
    bool    loadFromFileRelativePath(const QString &fileName, const QString &path);
    void updateTime(int ms, int totalms);//刷新时间
    double getTimePos(int ms);//返回指定时间在当前语句中的位置
    QString getLrcString(int offset);//取得歌词文本，参数：行数偏移量，负值表示提前
    bool isLrcEmpty();
private:
    QMap<int, QString> lrcMap;//存放歌词的QMap
    QList<int> timeList;//存放歌词时间的列表
    int                curLrcTime {0};  //目前将要显示的歌词时间
    int                nextLrcTime {0}; //下一句歌词的时间
};

#endif // LYRICS_H
