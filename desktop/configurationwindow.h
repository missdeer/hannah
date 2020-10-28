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

class ConfigurationWindow : public QMainWindow
{
    Q_OBJECT
    
public:
    ConfigurationWindow(QWidget *parent = nullptr);
    ~ConfigurationWindow();

    void onSearch(const QString &s);
    void onOpenUrl(const QString &s);
    void onOpenLink(const QString &s);

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

    void onReplyError(QNetworkReply::NetworkError code);

    void onReplyFinished();

    void onReplySslErrors(const QList<QSslError> &errors);

    void onReplyReadyRead();

    void onSystemTrayIconActivated(QSystemTrayIcon::ActivationReason reason);

    void onShowConfiguration();

    void onShowPlaylistManage();

private:
    Ui::ConfigurationWindow *ui;
    QMenu *                m_trayIconMenu;
    QSystemTrayIcon *      m_trayIcon;
    QSettings *            m_settings;
    QNetworkAccessManager *m_nam;
    QByteArray             m_reverseProxyAddr;
    QByteArray             m_playlistContent;

    void handle(const QString &url, bool needConfirm);
    void openLink(const QString &text);
    void restartReverseProxy();
    void startReverseProxy();
};

inline ConfigurationWindow *configurationWindow = nullptr;

#endif // MAINWINDOW_H
