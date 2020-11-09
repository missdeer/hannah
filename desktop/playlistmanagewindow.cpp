#include <QCloseEvent>
#include <QDir>
#include <QFile>
#include <QFileDialog>
#include <QInputDialog>
#include <QStandardPaths>
#include <map>

#include "playlistmanagewindow.h"

#include "playlistmodel.h"
#include "songlistmodel.h"
#include "sqlite3helper.h"
#include "ui_playlistmanagewindow.h"

PlaylistManageWindow::PlaylistManageWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::PlaylistManageWindow)
    , m_sqlite3Helper(new Sqlite3Helper)
    , m_playlistModel(new PlaylistModel(this))
    , m_songlistModel(new SonglistModel(this))
{
    ui->setupUi(this);
    ui->lstPlaylist->setModel(m_playlistModel);
    ui->tblSongs->setModel(m_songlistModel);

    QString fn = QDir::toNativeSeparators(QStandardPaths::writableLocation(QStandardPaths::AppLocalDataLocation) + "/default.hpls");
    if (QFile::exists(fn))
    {
        m_sqlite3Helper->openDatabase(fn);
    }
    else
    {
        m_sqlite3Helper->createDatabase(fn);
    }
    createDataTables();
}

PlaylistManageWindow::~PlaylistManageWindow()
{
    delete m_sqlite3Helper;
    delete ui;
}

void PlaylistManageWindow::closeEvent(QCloseEvent *event)
{
#if defined(Q_OS_MACOS)
    if (!event->spontaneous() || !isVisible())
    {
        return;
    }
#endif

    hide();
    event->ignore();
}

void PlaylistManageWindow::onAppendToPlaylist(const QStringList &s)
{
    Q_ASSERT(m_playlistModel);
    Q_ASSERT(m_songlistModel);
    m_songlistModel->appendToSonglist(s);
}

void PlaylistManageWindow::onClearAndAddToPlaylist(const QStringList &s)
{
    Q_ASSERT(m_playlistModel);
    Q_ASSERT(m_songlistModel);
    m_songlistModel->clearAndAddToSonglist(s);
}

void PlaylistManageWindow::onAppendToPlaylistFile(const QStringList &s)
{
    Q_ASSERT(m_playlistModel);
    Q_ASSERT(m_songlistModel);
    m_songlistModel->appendToSonglistFile(s);
}

void PlaylistManageWindow::onClearAndAddToPlaylistFile(const QStringList &s)
{
    Q_ASSERT(m_playlistModel);
    Q_ASSERT(m_songlistModel);
    m_songlistModel->clearAndAddToSonglistFile(s);
}

void PlaylistManageWindow::on_edtPlaylistFilter_textChanged(const QString &s)
{
    Q_ASSERT(m_playlistModel);
    m_playlistModel->filterPlaylist(s);
    bool isFiltered = m_playlistModel->isFilteredPlaylists();
    ui->btnAddPlaylist->setEnabled(!isFiltered);
    ui->btnDeletePlaylist->setEnabled(!isFiltered);
    ui->btnImportPlaylist->setEnabled(!isFiltered);
    ui->btnExportPlaylist->setEnabled(!isFiltered);
}

void PlaylistManageWindow::on_tblSongs_activated(const QModelIndex &index)
{
    Q_ASSERT(m_songlistModel);
}

void PlaylistManageWindow::on_btnAddPlaylist_clicked(bool)
{
    Q_ASSERT(m_playlistModel);
    m_playlistModel->addPlaylist();
}

void PlaylistManageWindow::on_btnDeletePlaylist_clicked(bool)
{
    auto model = ui->lstPlaylist->selectionModel();
    if (model->hasSelection())
    {
        Q_ASSERT(m_playlistModel);
        m_playlistModel->deletePlaylist(model->currentIndex().row());
    }
}

void PlaylistManageWindow::on_btnImportPlaylist_clicked(bool)
{
    QString fn = QFileDialog::getOpenFileName(this, tr("Import playlist"), "", tr("Playlist (*.m3u *.m3u8)"));
    Q_ASSERT(m_playlistModel);
}

void PlaylistManageWindow::on_btnExportPlaylist_clicked(bool)
{
    QString fn = QFileDialog::getSaveFileName(this, tr("Export playlist"), "", tr("Playlist (*.m3u)"));
    Q_ASSERT(m_playlistModel);
}

void PlaylistManageWindow::on_btnAddSongs_clicked(bool)
{
    QString lines = QInputDialog::getMultiLineText(this, tr("Add song(s)"), tr("Input song url, one url per line:"));
    Q_ASSERT(m_songlistModel);
}

void PlaylistManageWindow::on_btnDeleteSongs_clicked(bool)
{
    auto model = ui->tblSongs->selectionModel();
    if (model->hasSelection())
    {
        Q_ASSERT(m_songlistModel);
    }
}

void PlaylistManageWindow::on_btnImportSongs_clicked(bool)
{
    QStringList songs = QFileDialog::getOpenFileNames(
        this, tr("Import song(s)"), "", tr("Songs (*.m3u *.m3u8 *.mp1 *.mp2 *.mp3 *.wav *.ogg *.ape *.flac *.m4a *.aac *.caf *.wma *.opus)"));
    Q_ASSERT(m_songlistModel);
}

void PlaylistManageWindow::createDataTables()
{
    std::map<QString, QString> tablesCreationMap = {
        {"playlist", "CREATE TABLE IF NOT EXISTS playlist(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);"},
        {"song", "CREATE TABLE IF NOT EXISTS song(id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT);"},
        {"playlist_song_map", "CREATE TABLE IF NOT EXISTS playlist_song_map(id INTEGER PRIMARY KEY AUTOINCREMENT, playlist INTEGER, song INTEGER);"}};
    std::map<QString, QString> tablesMap = {{"table", "playlist"}, {"table", "song"}, {"table", "playlist_song_map"}};
    for (const auto &[type, name] : tablesMap)
    {
        if (!m_sqlite3Helper->checkTableIndexExists(type, name))
        {
            m_sqlite3Helper->execDML(tablesCreationMap[name]);
        }
    }
}
