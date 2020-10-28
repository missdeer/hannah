#ifndef PLAYLISTMANAGEWINDOW_H
#define PLAYLISTMANAGEWINDOW_H

#include <QMainWindow>

namespace Ui {
    class PlaylistManageWindow;
}
class QCloseEvent;

class PlaylistManageWindow : public QMainWindow
{
    Q_OBJECT
    
public:
    explicit PlaylistManageWindow(QWidget *parent = nullptr);
    ~PlaylistManageWindow();

protected:
    void closeEvent(QCloseEvent *event);

private:
    Ui::PlaylistManageWindow *ui;
};

inline PlaylistManageWindow *playlistManageWindow = nullptr;

#endif // PLAYLISTMANAGEWINDOW_H
