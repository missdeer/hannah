#ifndef PLAYLISTMANAGEWINDOW_H
#define PLAYLISTMANAGEWINDOW_H

#include <QMainWindow>

namespace Ui {
    class PlaylistManageWindow;
}
class QCloseEvent;
class PlaylistModel;
class SonglistModel;
class Sqlite3Helper;

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

    void on_btnAddPlaylist_clicked(bool checked);

    void on_btnDeletePlaylist_clicked(bool checked);

    void on_btnImportPlaylist_clicked(bool checked);

    void on_btnExportPlaylist_clicked(bool checked);

    void on_btnAddSongs_clicked(bool checked);

    void on_btnDeleteSongs_clicked(bool checked);

    void on_btnImportSongs_clicked(bool checked);

private:
    Ui::PlaylistManageWindow *ui;
    Sqlite3Helper *           m_sqlite3Helper{nullptr};
    PlaylistModel *           m_playlistModel{nullptr};
    SonglistModel *           m_songlistModel{nullptr};

    void createDataTables();
};

inline PlaylistManageWindow *playlistManageWindow = nullptr;

#endif // PLAYLISTMANAGEWINDOW_H
