#include <QClipboard>
#include <QCloseEvent>
#include <QCoreApplication>
#include <QDesktopServices>
#include <QEventLoop>
#include <QFileDialog>
#include <QMenu>
#include <QMessageBox>
#include <QNetworkAccessManager>
#include <QNetworkInterface>
#include <QProcess>
#include <QSettings>

#include "mainwindow.h"

#include "librp.h"
#include "ui_mainwindow.h"

#if defined(Q_OS_WIN)
#    include <Windows.h>
#    include <shellapi.h>
#    include <tchar.h>
#endif

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    m_settings = new QSettings(QSettings::IniFormat, QSettings::UserScope, "minidump.info", "Hannah");
    m_nam      = new QNetworkAccessManager(this);
    ui->setupUi(this);

    auto interfaces = QNetworkInterface::allInterfaces();
    for (const auto &i : interfaces)
    {
        ui->reverseProxyBindNetworkInterface->addItem(i.humanReadableName());
    }

    bool ok = true;
    ui->externalPlayerPath->setText(m_settings->value("externalPlayerPath").toString());
    ui->externalPlayerArguments->setText(m_settings->value("externalPlayerArguments").toString());
    ui->externalPlayerWorkingDir->setText(m_settings->value("externalPlayerWorkingDir").toString());
    ui->reverseProxyBindNetworkInterface->setCurrentText(m_settings->value("reverseProxyBindNetworkInterface", tr("-- Default --")).toString());
    ui->reverseProxyProxyType->setCurrentText(m_settings->value("reverseProxyProxyType", tr("None")).toString());
    ui->reverseProxyProxyAddress->setText(m_settings->value("reverseProxyProxyAddress").toString());
    auto state = m_settings->value("useExternalPlayer", 2).toInt(&ok);
    if (ok)
        ui->useExternalPlayer->setCheckState(Qt::CheckState(state));
    state = m_settings->value("reverseProxyAutoRedirect", 2).toInt(&ok);
    if (ok)
        ui->reverseProxyAutoRedirect->setCheckState(Qt::CheckState(state));
    state = m_settings->value("reverseProxyRedirect", 2).toInt(&ok);
    if (ok)
        ui->reverseProxyRedirect->setCheckState(Qt::CheckState(state));
    auto port = m_settings->value("reverseProxyListenPort", 8090).toInt(&ok);
    if (ok)
        ui->reverseProxyListenPort->setValue(port);

    connect(ui->reverseProxyBindNetworkInterface,
            &QComboBox::currentTextChanged,
            this,
            &MainWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged);
    connect(ui->useExternalPlayer, &QCheckBox::stateChanged, this, &MainWindow::onUseExternalPlayerStateChanged);
    connect(ui->browseExternalPlayer, &QPushButton::clicked, this, &MainWindow::onBrowseExternalPlayerClicked);
    connect(ui->externalPlayerPath, &QLineEdit::textChanged, this, &MainWindow::onExternalPlayerPathTextChanged);
    connect(ui->externalPlayerArguments, &QLineEdit::textChanged, this, &MainWindow::onExternalPlayerArgumentsTextChanged);
    connect(ui->externalPlayerWorkingDir, &QLineEdit::textChanged, this, &MainWindow::onExternalPlayerWorkingDirTextChanged);
    connect(ui->reverseProxyListenPort, QOverload<int>::of(&QSpinBox::valueChanged), this, &MainWindow::onReverseProxyListenPortValueChanged);
    connect(ui->reverseProxyAutoRedirect, &QCheckBox::stateChanged, this, &MainWindow::onReverseProxyAutoRedirectStateChanged);
    connect(ui->reverseProxyRedirect, &QCheckBox::stateChanged, this, &MainWindow::onReverseProxyRedirectStateChanged);
    connect(ui->reverseProxyProxyType, &QComboBox::currentTextChanged, this, &MainWindow::onReverseProxyProxyTypeCurrentTextChanged);
    connect(ui->reverseProxyProxyAddress, &QLineEdit::textChanged, this, &MainWindow::onReverseProxyProxyAddressTextChanged);

    QClipboard *clipboard = QGuiApplication::clipboard();
    connect(clipboard, &QClipboard::dataChanged, this, &MainWindow::onGlobalClipboardChanged);

    auto configAction = new QAction(tr("&Configuration"), this);
    connect(configAction, &QAction::triggered, this, &MainWindow::onShowConfiguration);

    auto quitAction = new QAction(tr("&Quit"), this);
    connect(quitAction, &QAction::triggered, qApp, &QCoreApplication::quit);

    m_trayIconMenu = new QMenu(this);
    m_trayIconMenu->addAction(tr("Netease"), []() { QDesktopServices::openUrl(QUrl("https://music.163.com")); });
    m_trayIconMenu->addAction(tr("QQ"), []() { QDesktopServices::openUrl(QUrl("https://y.qq.com")); });
    m_trayIconMenu->addAction(tr("Xiami"), []() { QDesktopServices::openUrl(QUrl("https://www.xiami.com")); });
    m_trayIconMenu->addAction(tr("Migu"), []() { QDesktopServices::openUrl(QUrl("https://music.migu.cn/v3")); });
    m_trayIconMenu->addAction(tr("Kugou"), []() { QDesktopServices::openUrl(QUrl("https://www.kugou.com")); });
    m_trayIconMenu->addAction(tr("Kuwo"), []() { QDesktopServices::openUrl(QUrl("http://kuwo.cn")); });
    m_trayIconMenu->addSeparator();
    m_trayIconMenu->addAction(configAction);
    m_trayIconMenu->addAction(quitAction);

    m_trayIcon = new QSystemTrayIcon(this);
    m_trayIcon->setContextMenu(m_trayIconMenu);
    m_trayIcon->setIcon(QIcon(":/hannah.png"));

    m_trayIcon->show();
    connect(m_trayIcon, &QSystemTrayIcon::activated, this, &MainWindow::onSystemTrayIconActivated);

    m_reverseProxyAddr = QString("localhost:%1").arg(ui->reverseProxyListenPort->value()).toUtf8();
    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

MainWindow::~MainWindow()
{
    StopReverseProxy();

    delete m_nam;

    m_settings->sync();
    delete m_settings;

    delete ui;
}

void MainWindow::closeEvent(QCloseEvent *event)
{
#ifdef Q_OS_MACOS
    if (!event->spontaneous() || !isVisible())
    {
        return;
    }
#endif
    if (m_trayIcon->isVisible())
    {
        hide();
        event->ignore();
    }
}

void MainWindow::onOpenUrl(QUrl url)
{
    onApplicationMessageReceived(url.toString());
}

void MainWindow::onApplicationMessageReceived(const QString &message)
{
    QString u = message;
    QString pattern = "hannah://play";
    if (u.startsWith(pattern))
    {
        auto index = u.indexOf("url=");
        if (index > pattern.length())
        {
            auto url = u.mid(index + 4);
            handle(url, false);
        }
    }
}

void MainWindow::onUseExternalPlayerStateChanged(int state)
{
    ui->externalPlayerArguments->setEnabled(state == Qt::Checked);
    ui->externalPlayerPath->setEnabled(state == Qt::Checked);
    ui->externalPlayerWorkingDir->setEnabled(state == Qt::Checked);
    ui->browseExternalPlayer->setEnabled(state == Qt::Checked);
    ui->browseExternalPlayerWorkingDir->setEnabled(state == Qt::Checked);

    Q_ASSERT(m_settings);
    m_settings->setValue("useExternalPlayer", state);
    m_settings->sync();
}

void MainWindow::onExternalPlayerPathTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("externalPlayerPath", text);
    if (ui->externalPlayerWorkingDir->text().isEmpty())
    {
        QFileInfo fi(text);
        if (fi.isAbsolute() && fi.isFile() && fi.exists())
        {
            ui->externalPlayerWorkingDir->setText(fi.absolutePath());
        }
    }
    m_settings->sync();
}

void MainWindow::onBrowseExternalPlayerClicked()
{
    QString fn = QFileDialog::getOpenFileName(this, tr("External Player"));
    ui->externalPlayerPath->setText(fn);
    m_settings->sync();
}

void MainWindow::onExternalPlayerArgumentsTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("externalPlayerArguments", text);
    m_settings->sync();
}

void MainWindow::onExternalPlayerWorkingDirTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("externalPlayerWorkingDir", text);
    m_settings->sync();
}

void MainWindow::onBrowseExternalPlayerWorkingDirClicked()
{
    QString dir = QFileDialog::getExistingDirectory(this, tr("Working Directory"));
    ui->externalPlayerWorkingDir->setText(dir);
    m_settings->sync();
}

void MainWindow::onReverseProxyListenPortValueChanged(int port)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyListenPort", port);
    m_settings->sync();
    StopReverseProxy();

    m_reverseProxyAddr = QString("localhost:%1").arg(ui->reverseProxyListenPort->value()).toUtf8();
    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyBindNetworkInterface", text);
    m_settings->sync();
    StopReverseProxy();

    QByteArray ba = ui->reverseProxyBindNetworkInterface->currentText().toUtf8();
    SetNetworkInterface(GoString {(const char *)ba.data(), (ptrdiff_t)ba.length()});

    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyAutoRedirectStateChanged(int state)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyAutoRedirect", state);
    m_settings->sync();
    StopReverseProxy();
    SetAutoRedirect(state);
    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyRedirectStateChanged(int state)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyRedirect", state);
    m_settings->sync();
    StopReverseProxy();
    SetRedirect(state);
    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyProxyTypeCurrentTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyProxyType", text);
    m_settings->sync();
    StopReverseProxy();
    if (text == "Http")
    {
        QByteArray ba = ui->reverseProxyProxyAddress->text().toUtf8();
        SetHttpProxy(GoString {(const char *)ba.data(), (ptrdiff_t)ba.length()});
    }
    else if (text == "Socks5")
    {
        QByteArray ba = ui->reverseProxyProxyAddress->text().toUtf8();
        SetSocks5Proxy(GoString {(const char *)ba.data(), (ptrdiff_t)ba.length()});
    }
    else
    {
        SetHttpProxy(GoString {nullptr, 0});
        SetSocks5Proxy(GoString {nullptr, 0});
    }
    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyProxyAddressTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyProxyAddress", text);
    m_settings->sync();
    StopReverseProxy();

    if (text == "Http")
    {
        QByteArray ba = ui->reverseProxyProxyAddress->text().toUtf8();
        SetHttpProxy(GoString {(const char *)ba.data(), (ptrdiff_t)ba.length()});
    }
    else if (text == "Socks5")
    {
        QByteArray ba = ui->reverseProxyProxyAddress->text().toUtf8();
        SetSocks5Proxy(GoString {(const char *)ba.data(), (ptrdiff_t)ba.length()});
    }
    else
    {
        SetHttpProxy(GoString {nullptr, 0});
        SetSocks5Proxy(GoString {nullptr, 0});
    }

    StartReverseProxy(GoString {(const char *)m_reverseProxyAddr.data(), (ptrdiff_t)m_reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onGlobalClipboardChanged()
{
    QClipboard *clipboard = QGuiApplication::clipboard();
    QString     text      = clipboard->text();
    QStringList patterns  = {"^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)discover\\/toplist\\?id=(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)playlist\\?id=(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)my\\/m\\/music\\/playlist\\?id=(\\d+)",
                            "^https?:\\/\\/www\\.xiami\\.com\\/collect\\/(\\d+)",
                            "^https?:\\/\\/y\\.qq\\.com\\/n\\/yqq\\/playlist\\/(\\d+)\\.html",
                            "^https?:\\/\\/www\\.kugou\\.com\\/yy\\/special\\/single\\/(\\d+)\\.html",
                            "^https?:\\/\\/(?:www\\.)?kuwo\\.cn\\/playlist_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/playlist\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)song\\?id=(\\d+)",
                            "^https?:\\/\\/www\\.xiami\\.com\\/song\\/(\\w+)",
                            "^https?:\\/\\/y\\.qq\\.com/n\\/yqq\\/song\\/(\\w+)\\.html",
                            "^https?:\\/\\/www\\.kugou\\.com\\/song\\/#hash=([0-9A-F]+)",
                            "^https?:\\/\\/(?:www\\.)kuwo.cn\\/play_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/song\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/weapi\\/v1\\/artist\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)artist\\?id=(\\d+)",
                            "^https?:\\/\\/y\\.qq\\.com\\/n\\/yqq\\/singer\\/(\\w+)\\.html",
                            "^https?:\\/\\/www\\.xiami\\.com\\/artist\\/(\\w+)",
                            "^https?:\\/\\/www\\.xiami\\.com\\/list\\?scene=artist&type=\\w+&query={%22artistId%22:%22(\\d+)%22}",
                            "^https?:\\/\\/www\\.xiami\\.com\\/list\\?scene=artist&type=\\w+&query={\"artistId\":\"(\\d+)\"}",
                            "^https?:\\/\\/(?:www\\.)?kuwo\\.cn\\/singer_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/artist\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/weapi\\/v1\\/album\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)album\\?id=(\\d+)",
                            "^https?:\\/\\/y\\.qq\\.com\\/n\\/yqq\\/album\\/(\\w+)\\.html",
                            "^https?:\\/\\/www\\.xiami\\.com\\/album\\/(\\w+)",
                            "^https?:\\/\\/(?:www\\.)?kuwo\\.cn\\/album_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/album\\/(\\d+)"};
    for (const auto &p : patterns)
    {
        QRegularExpression r(p);
        auto               match = r.match(text);
        if (match.hasMatch())
        {
            handle(QString(QUrl::toPercentEncoding(text)), true);
            break;
        }
    }
}

void MainWindow::onReplyError(QNetworkReply::NetworkError code)
{
    Q_UNUSED(code);
#if !defined(QT_NO_DEBUG)
    auto reply = qobject_cast<QNetworkReply *>(sender());
    Q_ASSERT(reply);
    qDebug() << reply->errorString();
#endif
}

void MainWindow::onReplyFinished()
{
    auto reply = qobject_cast<QNetworkReply *>(sender());
    reply->deleteLater();

#if !defined(QT_NO_DEBUG)
    qDebug() << " finished: " << QString(m_playlistContent).left(256) << "\n";
#endif
    auto fn = QStandardPaths::writableLocation(QStandardPaths::TempLocation) + "/hannah.m3u";
    if (m_playlistContent.isEmpty())
    {
        QFile::remove(fn);
        return;
    }
    QFile f(fn);
    if (f.open(QIODevice::WriteOnly | QIODevice::Truncate))
    {
        f.write(m_playlistContent);
        f.close();
    }
}

void MainWindow::onReplySslErrors(const QList<QSslError> &
#if !defined(QT_NO_DEBUG)
                                      errors
#endif
)
{
#if !defined(QT_NO_DEBUG)
    for (const auto &e : errors)
    {
        qDebug() << "ssl error:" << e.errorString();
    }
#endif
}

void MainWindow::onReplyReadyRead()
{
    auto *reply      = qobject_cast<QNetworkReply *>(sender());
    int   statusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
    if (statusCode >= 200 && statusCode < 300)
    {
        m_playlistContent.append(reply->readAll());
    }
}

void MainWindow::onSystemTrayIconActivated(QSystemTrayIcon::ActivationReason reason)
{
    if (reason == QSystemTrayIcon::DoubleClick)
    {
        onShowConfiguration();
    }
}

void MainWindow::onShowConfiguration()
{
    if (isHidden())
    {
        showNormal();
    }
    activateWindow();
    raise();
}

void MainWindow::handle(const QString &url, bool needConfirm)
{
    auto player     = QDir::toNativeSeparators(ui->externalPlayerPath->text());
    auto arguments  = ui->externalPlayerArguments->text();
    auto workingDir = QDir::toNativeSeparators(ui->externalPlayerWorkingDir->text());

    QFileInfo fi(player);

    if (!fi.exists())
    {
        if (!needConfirm)
            QMessageBox::critical(this, tr("Erorr"), tr("External player path not configured properly"));
        return;
    }

    m_playlistContent.clear();

    QNetworkRequest req(QUrl::fromUserInput(QString("http://localhost:%1/m3u/generate?u=").arg(ui->reverseProxyListenPort->value()) + url));
    auto            reply = m_nam->get(req);
    connect(reply, &QNetworkReply::finished, this, &MainWindow::onReplyFinished);
    connect(reply, &QNetworkReply::readyRead, this, &MainWindow::onReplyReadyRead);
    connect(reply, &QNetworkReply::errorOccurred, this, &MainWindow::onReplyError);
    connect(reply, &QNetworkReply::sslErrors, this, &MainWindow::onReplySslErrors);
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    auto localTempPlaylist = QDir::toNativeSeparators(QStandardPaths::writableLocation(QStandardPaths::TempLocation) + "/hannah.m3u");
    if (!QFile::exists(localTempPlaylist))
    {
        QMessageBox::critical(this, tr("Error"), tr("Can't get song(s), maybe VIP is requested."));
        return;
    }

    if (needConfirm &&
        QMessageBox::question(this, tr("Confirm"), tr("Play song(s) by %1?").arg(player), QMessageBox::Ok | QMessageBox::Cancel) != QMessageBox::Ok)
        return;
#if defined(Q_OS_MAC)
    if (fi.isBundle() && player.endsWith(".app"))
    {
        auto script = QString("tell application \"%1\" to open \"%2\"").arg(player, localTempPlaylist);
        QProcess::startDetached("/usr/bin/osascript", {"-e", script}, workingDir);
        return;
    }
    else
    {
        QFile f(":/rc/runInTerminal.app.scpt");
        if (f.open(QIODevice::ReadOnly))
        {
            auto data = f.readAll();
            f.close();

            auto  runInTerminalScriptPath = QStandardPaths::writableLocation(QStandardPaths::TempLocation) + "/runInTerminal.app.scpt";
            QFile tf(runInTerminalScriptPath);
            if (tf.open(QIODevice::WriteOnly))
            {
                tf.write(data);
                tf.close();
                QStringList args = {QDir::toNativeSeparators(runInTerminalScriptPath),
                                    QString("\"%1\" %2 \"%3\"").arg(player, arguments, localTempPlaylist)};
                QProcess::startDetached("/usr/bin/osascript", args, workingDir);
                return;
            }
        }
    }
#elif defined(Q_OS_WIN)
    auto args = arguments.split(" ");
    args << localTempPlaylist;
    args.removeAll("");
    ::ShellExecuteW((HWND)winId(),
                    L"open",
                    player.toStdWString().c_str(),
                    args.join(" ").toStdWString().c_str(),
                    workingDir.toStdWString().c_str(),
                    SW_SHOWNORMAL);
    return;
#else
#endif
    auto argv = arguments.split(" ");
    argv << localTempPlaylist;
    argv.removeAll("");
    QProcess::startDetached(player, argv, workingDir);
}
