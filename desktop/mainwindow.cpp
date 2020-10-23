#include <QClipboard>
#include <QCloseEvent>
#include <QCoreApplication>
#include <QDesktopServices>
#include <QFileDialog>
#include <QMenu>
#include <QMessageBox>
#include <QNetworkInterface>
#include <QProcess>
#include <QSettings>
#include <QSystemTrayIcon>

#include "mainwindow.h"

#include "librp.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    settings = new QSettings(QSettings::IniFormat, QSettings::UserScope, "minidump.info", "Hannah");
    ui->setupUi(this);

    auto interfaces = QNetworkInterface::allInterfaces();
    for (const auto &i : interfaces)
    {
        ui->reverseProxyBindNetworkInterface->addItem(i.humanReadableName());
    }

    bool ok = true;
    ui->externalPlayerPath->setText(settings->value("externalPlayerPath").toString());
    ui->externalPlayerArguments->setText(settings->value("externalPlayerArguments").toString());
    ui->externalPlayerWorkingDir->setText(settings->value("externalPlayerWorkingDir").toString());
    ui->reverseProxyBindNetworkInterface->setCurrentText(settings->value("reverseProxyBindNetworkInterface", tr("-- Default --")).toString());
    ui->reverseProxyProxyType->setCurrentText(settings->value("reverseProxyProxyType", tr("None")).toString());
    ui->reverseProxyProxyAddress->setText(settings->value("reverseProxyProxyAddress").toString());
    auto state = settings->value("useExternalPlayer", true).toInt(&ok);
    if (ok)
        ui->useExternalPlayer->setCheckState(Qt::CheckState(state));
    state = settings->value("reverseProxyAutoRedirect", true).toInt(&ok);
    if (ok)
        ui->reverseProxyAutoRedirect->setCheckState(Qt::CheckState(state));
    state = settings->value("reverseProxyRedirect", true).toInt(&ok);
    if (ok)
        ui->reverseProxyRedirect->setCheckState(Qt::CheckState(state));
    auto port = settings->value("reverseProxyListenPort", 8090).toInt(&ok);
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
    connect(configAction, &QAction::triggered, this, [this]() {
        if (isHidden())
        {
            showNormal();
        }
        activateWindow();
        raise();
    });

    auto quitAction = new QAction(tr("&Quit"), this);
    connect(quitAction, &QAction::triggered, qApp, &QCoreApplication::quit);

    trayIconMenu = new QMenu(this);
    trayIconMenu->addAction(tr("Netease"), []() { QDesktopServices::openUrl(QUrl("https://music.163.com")); });
    trayIconMenu->addAction(tr("QQ"), []() { QDesktopServices::openUrl(QUrl("https://y.qq.com")); });
    trayIconMenu->addAction(tr("Xiami"), []() { QDesktopServices::openUrl(QUrl("https://www.xiami.com")); });
    trayIconMenu->addAction(tr("Migu"), []() { QDesktopServices::openUrl(QUrl("https://music.migu.cn/v3")); });
    trayIconMenu->addAction(tr("Kugou"), []() { QDesktopServices::openUrl(QUrl("https://www.kugou.com")); });
    trayIconMenu->addAction(tr("Kuwo"), []() { QDesktopServices::openUrl(QUrl("http://kuwo.cn")); });
    trayIconMenu->addSeparator();
    trayIconMenu->addAction(configAction);
    trayIconMenu->addAction(quitAction);

    trayIcon = new QSystemTrayIcon(this);
    trayIcon->setContextMenu(trayIconMenu);
    trayIcon->setIcon(QIcon(":/hannah.png"));

    trayIcon->show();

    reverseProxyAddr = QString(":%1").arg(ui->reverseProxyListenPort->value()).toUtf8();
    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

MainWindow::~MainWindow()
{
    StopReverseProxy();

    settings->sync();
    delete settings;
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
    if (trayIcon->isVisible())
    {
        hide();
        event->ignore();
    }
}

void MainWindow::onOpenUrl(QUrl url)
{
    QString u = url.toString();
    if (u.startsWith("hannah://play?url="))
    {
        u = u.replace("hannah://play?url=", QString("http://localhost:%1/m3u/generate?u=").arg(ui->reverseProxyListenPort->value()));
        handle(u);
    }
}

void MainWindow::onApplicationMessageReceived(const QString &message)
{
    QString u = message;
    if (u.startsWith("hannah://play?url="))
    {
        u = u.replace("hannah://play?url=", QString("http://localhost:%1/m3u/generate?u=").arg(ui->reverseProxyListenPort->value()));
        handle(u);
    }
}

void MainWindow::onUseExternalPlayerStateChanged(int state)
{
    ui->externalPlayerArguments->setEnabled(state == Qt::Checked);
    ui->externalPlayerPath->setEnabled(state == Qt::Checked);
    ui->externalPlayerWorkingDir->setEnabled(state == Qt::Checked);
    ui->browseExternalPlayer->setEnabled(state == Qt::Checked);
    ui->browseExternalPlayerWorkingDir->setEnabled(state == Qt::Checked);

    Q_ASSERT(settings);
    settings->setValue("useExternalPlayer", state);
    settings->sync();
}

void MainWindow::onExternalPlayerPathTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("externalPlayerPath", text);
    settings->sync();
}

void MainWindow::onBrowseExternalPlayerClicked()
{
    QString fn = QFileDialog::getOpenFileName(this, tr("External Player"));
    ui->externalPlayerPath->setText(fn);
    settings->sync();
}

void MainWindow::onExternalPlayerArgumentsTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("externalPlayerArguments", text);
    settings->sync();
}

void MainWindow::onExternalPlayerWorkingDirTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("externalPlayerWorkingDir", text);
    settings->sync();
}

void MainWindow::onBrowseExternalPlayerWorkingDirClicked()
{
    QString dir = QFileDialog::getExistingDirectory(this, tr("Working Directory"));
    ui->externalPlayerWorkingDir->setText(dir);
    settings->sync();
}

void MainWindow::onReverseProxyListenPortValueChanged(int port)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyListenPort", port);
    settings->sync();
    StopReverseProxy();

    reverseProxyAddr = QString(":%1").arg(ui->reverseProxyListenPort->value()).toUtf8();
    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyBindNetworkInterface", text);
    settings->sync();
    StopReverseProxy();

    QByteArray ba = ui->reverseProxyBindNetworkInterface->currentText().toUtf8();
    SetNetworkInterface(GoString {(const char *)ba.data(), (ptrdiff_t)ba.length()});

    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyAutoRedirectStateChanged(int state)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyAutoRedirect", state);
    settings->sync();
    StopReverseProxy();
    SetAutoRedirect(state);
    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyRedirectStateChanged(int state)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyRedirect", state);
    settings->sync();
    StopReverseProxy();
    SetRedirect(state);
    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyProxyTypeCurrentTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyProxyType", text);
    settings->sync();
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
    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onReverseProxyProxyAddressTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyProxyAddress", text);
    settings->sync();
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

    StartReverseProxy(GoString {(const char *)reverseProxyAddr.data(), (ptrdiff_t)reverseProxyAddr.length()}, GoString {nullptr, 0});
}

void MainWindow::onGlobalClipboardChanged()
{
    QClipboard *clipboard = QGuiApplication::clipboard();
    QString     text      = clipboard->text();
    QStringList patterns  = {"^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)discover\\/toplist\?id=(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)playlist\?id=(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)my\\/m\\/music\\/playlist\?id=(\\d+)",
                            "^https?:\\/\\/www\\.xiami\\.com\\/collect\\/(\\d+)",
                            "^https?:\\/\\/y\\.qq\\.com\\/n\\/yqq\\/playlist\\/(\\d+)\\.html",
                            "^https?:\\/\\/www\\.kugou\\.com\\/yy\\/special\\/single\\/(\\d+)\\.html",
                            "^https?:\\/\\/(?:www\\.)?kuwo\\.cn\\/playlist_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/playlist\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)song\?id=(\\d+)",
                            "^https?:\\/\\/www\\.xiami\\.com\\/song\\/(\\w+)",
                            "^https?:\\/\\/y\\.qq\\.com/n\\/yqq\\/song\\/(\\w+)\\.html",
                            "^https?:\\/\\/www\\.kugou\\.com\\/song\\/#hash=([0-9A-F]+)",
                            "^https?:\\/\\/(?:www\\.)kuwo.cn\\/play_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/song\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/weapi\\/v1\\/artist\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)artist\?id=(\\d+)",
                            "^https?:\\/\\/y\\.qq\\.com\\/n\\/yqq\\/singer\\/(\\w+)\\.html",
                            "^https?:\\/\\/www\\.xiami\\.com\\/artist\\/(\\w+)",
                            "^https?:\\/\\/www\\.xiami\\.com\\/list\?scene=artist&type=\\w+&query={%22artistId%22:%22(\\d+)%22}",
                            "^https?:\\/\\/www\\.xiami\\.com\\/list\?scene=artist&type=\\w+&query={\"artistId\":\"(\\d+)\"}",
                            "^https?:\\/\\/(?:www\\.)?kuwo\\.cn\\/singer_detail\\/(\\d+)",
                            "^https?:\\/\\/music\\.migu\\.cn\\/v3\\/music\\/artist\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/weapi\\/v1\\/album\\/(\\d+)",
                            "^https?:\\/\\/music\\.163\\.com\\/(?:#\\/)album\?id=(\\d+)",
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
            handle(text);
            break;
        }
    }
}

void MainWindow::handle(const QString &url)
{
    auto player     = ui->externalPlayerPath->text();
    auto arguments  = ui->externalPlayerArguments->text();
    auto workingDir = ui->externalPlayerWorkingDir->text();

    QFileInfo fi(player);

    if (!fi.exists())
    {
        QMessageBox::critical(this, tr("Erorr"), tr("External player path not configured properly"));
        return;
    }
#if defined(Q_OS_MAC)
    if (fi.isBundle() && player.endsWith(".app"))
    {
        QStringList args = {player, "--args"};
        args << arguments.split(" ") << url;
        args.removeAll("");
        QProcess::startDetached("/usr/bin/open", args, workingDir);
        return;
    }

    QFile f(":/rc/runInTerminal.app.scpt");
    if (f.open(QIODevice::ReadOnly))
    {
        auto data = f.readAll();
        f.close();

        auto  path = QStandardPaths::writableLocation(QStandardPaths::TempLocation) + "/runInTerminal.app.scpt";
        QFile tf(path);
        if (tf.open(QIODevice::WriteOnly))
        {
            tf.write(data);
            tf.close();
            QStringList args = {QDir::toNativeSeparators(path), QString("%1 %2 %3").arg(player, arguments, url)};
            QProcess::startDetached("/usr/bin/osascript", args, workingDir);
            return;
        }
    }
#elif defined(Q_OS_WIN)
    auto args = arguments.split(" ");
    args << url;
    args.removeAll("");
    ::ShellExecuteW(
        nullptr, L"open", player.toStdWString().c_str(), args.join(" ").toStdWString().c_str(), workingDir.toStdWString().c_str(), SW_SHOWNORMAL);
    return;
#else
#endif
    auto args = arguments.split(" ");
    args << url;
    args.removeAll("");
    QProcess::startDetached(player, args, workingDir);
}
