#ifndef PLAYLIST_H
#define PLAYLIST_H

#include <QWidget>
#include <QtCore>
#include <QtGui>
#include <QtWidgets>

#include "player.h"

namespace Ui
{
    class PlayList;
}

class PlayList : public QWidget
{
    Q_OBJECT

public:
    explicit PlayList(Player *player, QWidget *parent = 0);
    ~PlayList();
    bool           fixSuffix(const QString &uri);
    bool           isEmpty();
    void           add(const QString &fileName);
    void           insert(int index, const QString &fileName);
    void           remove(int index);
    void           clearAll();
    int            getLength();
    int            getIndex();
    QString        next(bool isLoop = false);
    QString        previous(bool isLoop = false);
    const QString &playIndex(int index);
    const QString &getFileNameForIndex(int index);
    const QString &getCurFile();
    QString        playLast();
    void           tableUpdate();
    void           saveToFile(const QString &fileName);
    void           readFromFile(const QString &fileName);

private slots:
    void on_deleteButton_clicked();
    void on_playListTable_cellDoubleClicked(int row, int);
    void on_clearButton_clicked();
    void on_insertButton_clicked();
    void on_addButton_clicked();
    void on_insertUrlButton_clicked();
    void on_addUrlButton_clicked();
    void on_searchButton_clicked();
    void on_searchNextButton_clicked();
    void on_setLenFilButton_clicked();

signals:
    void callPlayer();

private:
    Ui::PlayList * ui;
    QList<QString> m_trackList;
    QList<QString> m_timeList;
    Player *       m_player;
    int            m_curIndex {0};
    int            m_lengthFilter {0};

protected:
    void dragEnterEvent(QDragEnterEvent *event) override;
    void dropEvent(QDropEvent *event) override;
};

#endif // PLAYLIST_H
