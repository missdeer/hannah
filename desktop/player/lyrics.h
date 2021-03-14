#ifndef LYRICS_H
#define LYRICS_H
#include <QList>
#include <QMap>
#include <QString>

class Lyrics
{
public:
    bool    resolve(const QString &fileName, bool isLrc = false);
    bool    loadFromLrcDir(const QString &fileName);
    bool    loadFromFileRelativePath(const QString &fileName, const QString &path);
    void    updateTime(int ms, int totalms);
    double  getTimePos(int ms);
    QString getLrcString(int offset);
    bool    isLrcEmpty();

private:
    QMap<int, QString> lrcMap;
    QList<int>         timeList;
    int                curLrcTime {0};
    int                nextLrcTime {0};
};

#endif // LYRICS_H
