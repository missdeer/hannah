#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QUrl>

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
class QSystemTrayIcon;
class QMenu;
class QCloseEvent;
class QSettings;
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT
    
public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

protected:
    void closeEvent(QCloseEvent *event);

public slots:
    void onOpenUrl(QUrl url);

    void onApplicationMessageReceived(const QString &message);
private slots:
    void onUseExternalPlayerStateChanged(int state);

    void onExternalPlayerPathTextChanged(const QString &text);

    void onBrowseExternalPlayerClicked();

    void onExternalPlayerArgumentsTextChanged(const QString &text);

    void onExternalPlayerWorkingDirTextChanged(const QString &text);

    void onBrowseExternalPlayerWorkingDirClicked();

    void onReverseProxyListenPortValueChanged(int port);

    void onReverseProxyBindNetworkInterfaceCurrentTextChanged(const QString &text);

    void onReverseProxyAutoRedirectStateChanged(int state);

    void onReverseProxyRedirectStateChanged(int state);

    void onReverseProxyProxyTypeCurrentTextChanged(const QString &text);

    void onReverseProxyProxyAddressTextChanged(const QString &text);

    void onGlobalClipboardChanged();

private:
    Ui::MainWindow * ui;
    QMenu *          trayIconMenu;
    QSystemTrayIcon *trayIcon;
    QSettings *      settings;
    QByteArray       reverseProxyAddr;

    void handle(const QString &url);
};
#endif // MAINWINDOW_H
