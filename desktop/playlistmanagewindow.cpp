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
