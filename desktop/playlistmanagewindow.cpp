#include <QCloseEvent>

#include "playlistmanagewindow.h"

#include "ui_playlistmanagewindow.h"

PlaylistManageWindow::PlaylistManageWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::PlaylistManageWindow)
{
    ui->setupUi(this);
}

PlaylistManageWindow::~PlaylistManageWindow()
{
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

void PlaylistManageWindow::onAppendToPlaylist(const QStringList &s) {}

void PlaylistManageWindow::onClearAndAddToPlaylist(const QStringList &s) {}

void PlaylistManageWindow::onAppendToPlaylistFile(const QStringList &s) {}

void PlaylistManageWindow::onClearAndAddToPlaylistFile(const QStringList &s) {}

void PlaylistManageWindow::on_edtPlaylistFilter_textChanged(const QString &s) {}

void PlaylistManageWindow::on_tblSongs_activated(const QModelIndex &index) {}

void PlaylistManageWindow::on_btnAddPlaylist_triggered(QAction *) {}

void PlaylistManageWindow::on_btnDeletePlaylist_triggered(QAction *) {}

void PlaylistManageWindow::on_btnImportPlaylist_triggered(QAction *) {}

void PlaylistManageWindow::on_btnSavePlaylist_triggered(QAction *) {}

void PlaylistManageWindow::on_btnAddSongs_triggered(QAction *) {}

void PlaylistManageWindow::on_btnDeleteSongs_triggered(QAction *) {}

void PlaylistManageWindow::on_btnImportSongs_triggered(QAction *) {}
