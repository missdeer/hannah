#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QNetworkReply>
#include <QSystemTrayIcon>
#include <QUrl>

QT_BEGIN_NAMESPACE
namespace Ui
{
    class ConfigurationWindow;
}
class QMenu;
class QCloseEvent;
class QSettings;
class QNetworkAccessManager;
QT_END_NAMESPACE

class BeastServerRunner;

class ConfigurationWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit ConfigurationWindow(BeastServerRunner &runner, QWidget *parent = nullptr);
    ~ConfigurationWindow();

    void onSearch(const QString &s);
    void onOpenUrl(const QString &s);
    void onOpenLink(const QString &s);

protected:
    void closeEvent(QCloseEvent *event);

public slots:
    void onOpenUrl(const QUrl &url);

    void onApplicationMessageReceived(const QString &message);

    void onShowConfiguration();
private slots:
    void onUseBuiltinPlayerStateChanged(bool checked);

    void onUseExternalPlayerStateChanged(bool checked);

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

    void onReplyError(QNetworkReply::NetworkError code);

    void onReplyFinished();

    void onReplySslErrors(const QList<QSslError> &errors);

    void onReplyReadyRead();

    void onSystemTrayIconActivated(QSystemTrayIcon::ActivationReason reason);

    void onShowPlaylistManage();

    void onShowHideBuiltinPlayer();

private:
    Ui::ConfigurationWindow *ui;
    BeastServerRunner       &m_runner;
    QMenu                   *m_trayIconMenu;
    QSystemTrayIcon         *m_trayIcon;
    QSettings               *m_settings;
    QNetworkAccessManager   *m_nam;
    QByteArray               m_reverseProxyAddr;
    QByteArray               m_playlistContent;

    void handle(const QString &url, bool needConfirm);
    void openLink(const QString &text);
    void restartReverseProxy();
    void startReverseProxy();
    void initOutputDevices();
    void initNetworkInterfaces();
};

inline ConfigurationWindow *configurationWindow = nullptr;

#endif // MAINWINDOW_H
