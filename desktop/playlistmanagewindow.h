#ifndef PLAYLISTMANAGEWINDOW_H
#define PLAYLISTMANAGEWINDOW_H

#include <QMainWindow>

namespace Ui {
    class PlaylistManageWindow;
}
class QCloseEvent;
class PlaylistModel;
class SonglistModel;

class PlaylistManageWindow : public QMainWindow
{
    Q_OBJECT
    
public:
    explicit PlaylistManageWindow(QWidget *parent = nullptr);
    ~PlaylistManageWindow();

    void onAppendToPlaylist(const QStringList &s);
    void onClearAndAddToPlaylist(const QStringList &s);
    void onAppendToPlaylistFile(const QStringList &s);
    void onClearAndAddToPlaylistFile(const QStringList &s);

protected:
    void closeEvent(QCloseEvent *event);

private slots:
    void on_edtPlaylistFilter_textChanged(const QString &s);

    void on_tblSongs_activated(const QModelIndex &index);

    void on_btnAddPlaylist_triggered(QAction *);

    void on_btnDeletePlaylist_triggered(QAction *);

    void on_btnImportPlaylist_triggered(QAction *);

    void on_btnSavePlaylist_triggered(QAction *);

    void on_btnAddSongs_triggered(QAction *);

    void on_btnDeleteSongs_triggered(QAction *);

    void on_btnImportSongs_triggered(QAction *);

private:
    Ui::PlaylistManageWindow *ui;
    PlaylistModel *           m_playlistModel;
    SonglistModel *           m_songlistModel;
};

inline PlaylistManageWindow *playlistManageWindow = nullptr;

#endif // PLAYLISTMANAGEWINDOW_H
