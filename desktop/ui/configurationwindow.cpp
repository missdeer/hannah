#include <QClipboard>
#include <QCloseEvent>
#include <QComboBox>
#include <QCoreApplication>
#include <QDesktopServices>
#include <QEventLoop>
#include <QFileDialog>
#include <QMenu>
#include <QMessageBox>
#include <QNetworkAccessManager>
#include <QNetworkInterface>
#include <QProcess>
#include <QQmlApplicationEngine>
#include <QQuickStyle>
#include <QSettings>
#include <QStandardItem>
#include <QStandardPaths>

#if defined(Q_OS_WIN)
#    include <Windows.h>

#    include <shellapi.h>
#    include <tchar.h>

#    include "bassasio.h"
#    include "basswasapi.h"

#endif

#include "BeastServerRunner.h"
#include "bass.h"
#include "comboboxdelegate.h"
#include "configurationwindow.h"
#include "externalreverseproxyrunner.h"
#include "playlistmanagewindow.h"
#include "qmlplayer.h"
#include "ui_configurationwindow.h"

ConfigurationWindow::ConfigurationWindow(BeastServerRunner &runner, ExternalReverseProxyRunner &externalReverseProxyRunner, QWidget *parent)
    : QMainWindow(parent),
      ui(new Ui::ConfigurationWindow),
      m_builtinReverseProxyRunner(runner),
      m_externalReverseProxyRunner(externalReverseProxyRunner),
      m_nam(new QNetworkAccessManager(this))
{
    m_settings = new QSettings(QSettings::IniFormat, QSettings::UserScope, "minidump.info", "Hannah");

    ui->setupUi(this);

    ui->cbOutputDevices->setItemDelegate(new ComboBoxDelegate);

    initNetworkInterfaces();
    initOutputDevices();

    bool ok = true;
    ui->externalPlayerPath->setText(m_settings->value("externalPlayerPath").toString());
    ui->externalPlayerArguments->setText(m_settings->value("externalPlayerArguments").toString());
    ui->externalPlayerWorkingDir->setText(m_settings->value("externalPlayerWorkingDir").toString());
    ui->reverseProxyBindNetworkInterface->setCurrentText(m_settings->value("reverseProxyBindNetworkInterface", tr("-- Default --")).toString());
    ui->reverseProxyProxyType->setCurrentText(m_settings->value("reverseProxyProxyType", tr("None")).toString());
    ui->reverseProxyProxyAddress->setText(m_settings->value("reverseProxyProxyAddress").toString());
    bool bUseExternalPlayer = m_settings->value("useExternalPlayer", false).toBool();
    ui->useExternalPlayer->setChecked(bUseExternalPlayer);
    onUseExternalPlayerStateChanged(bUseExternalPlayer);
    ui->useBuiltinPlayer->setChecked(!bUseExternalPlayer);
    onUseBuiltinPlayerStateChanged(!bUseExternalPlayer);
    auto state = m_settings->value("reverseProxyAutoRedirect", 2).toInt(&ok);
    if (ok)
    {
        ui->reverseProxyAutoRedirect->setCheckState(static_cast<Qt::CheckState>(state));
    }
    state = m_settings->value("reverseProxyRedirect", 2).toInt(&ok);
    if (ok)
    {
        ui->reverseProxyRedirect->setCheckState(static_cast<Qt::CheckState>(state));
    }
    state = m_settings->value("useExternalReverseProxy", 2).toInt(&ok);
    if (ok)
    {
        ui->cbUseExternalReverseProxy->setCheckState(static_cast<Qt::CheckState>(state));
    }
    ui->edtExternalReverseProxyPath->setText(m_settings->value("externalReverseProxyPath").toString());
    auto port = m_settings->value("reverseProxyListenPort", 8090).toInt(&ok);
    if (ok)
    {
        ui->reverseProxyListenPort->setValue(port);
    }
    connect(ui->reverseProxyBindNetworkInterface,
            &QComboBox::currentTextChanged,
            this,
            &ConfigurationWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged);
    connect(ui->useBuiltinPlayer, &QRadioButton::toggled, this, &ConfigurationWindow::onUseBuiltinPlayerStateChanged);
    connect(ui->useExternalPlayer, &QRadioButton::toggled, this, &ConfigurationWindow::onUseExternalPlayerStateChanged);
    connect(ui->browseExternalPlayer, &QPushButton::clicked, this, &ConfigurationWindow::onBrowseExternalPlayerClicked);
    connect(ui->browseExternalPlayerWorkingDir, &QPushButton::clicked, this, &ConfigurationWindow::onBrowseExternalPlayerWorkingDirClicked);
    connect(ui->externalPlayerPath, &QLineEdit::textChanged, this, &ConfigurationWindow::onExternalPlayerPathTextChanged);
    connect(ui->externalPlayerArguments, &QLineEdit::textChanged, this, &ConfigurationWindow::onExternalPlayerArgumentsTextChanged);
    connect(ui->externalPlayerWorkingDir, &QLineEdit::textChanged, this, &ConfigurationWindow::onExternalPlayerWorkingDirTextChanged);
    connect(
        ui->reverseProxyListenPort, QOverload<int>::of(&QSpinBox::valueChanged), this, &ConfigurationWindow::onReverseProxyListenPortValueChanged);
    connect(ui->reverseProxyAutoRedirect, &QCheckBox::stateChanged, this, &ConfigurationWindow::onReverseProxyAutoRedirectStateChanged);
    connect(ui->reverseProxyRedirect, &QCheckBox::stateChanged, this, &ConfigurationWindow::onReverseProxyRedirectStateChanged);
    connect(ui->reverseProxyProxyType, &QComboBox::currentTextChanged, this, &ConfigurationWindow::onReverseProxyProxyTypeCurrentTextChanged);
    connect(ui->reverseProxyProxyAddress, &QLineEdit::textChanged, this, &ConfigurationWindow::onReverseProxyProxyAddressTextChanged);
    connect(ui->cbUseExternalReverseProxy, &QCheckBox::stateChanged, this, &ConfigurationWindow::onUseExternalReverseProxyStateChanged);
    connect(ui->btnBrowseExternalReverseProxy, &QPushButton::clicked, this, &ConfigurationWindow::onBrowseExternalReverseProxy);
    connect(ui->edtExternalReverseProxyPath, &QLineEdit::textChanged, this, &ConfigurationWindow::onExternalReverseProxyPathChanged);
    auto *clipboard = QGuiApplication::clipboard();
    connect(clipboard, &QClipboard::dataChanged, this, &ConfigurationWindow::onGlobalClipboardChanged);

    auto *configAction = new QAction(tr("&Configuration"), this);
    connect(configAction, &QAction::triggered, this, &ConfigurationWindow::onShowConfiguration);

    auto *showHidePlayerAction = new QAction(tr("Show/Hide &Player"), this);
    connect(showHidePlayerAction, &QAction::triggered, this, &ConfigurationWindow::onShowHideBuiltinPlayer);

    auto *playlistManageAction = new QAction(tr("Playlist Manage"), this);
    connect(playlistManageAction, &QAction::triggered, this, &ConfigurationWindow::onShowPlaylistManage);

    auto *quitAction = new QAction(tr("&Quit"), this);
    connect(quitAction, &QAction::triggered, qApp, &QCoreApplication::quit);

    m_trayIconMenu = new QMenu(this);
    m_trayIconMenu->addAction(tr("Netease"), []() { QDesktopServices::openUrl(QUrl("https://music.163.com")); });
    m_trayIconMenu->addAction(tr("QQ"), []() { QDesktopServices::openUrl(QUrl("https://y.qq.com")); });
    m_trayIconMenu->addAction(tr("Migu"), []() { QDesktopServices::openUrl(QUrl("https://music.migu.cn/v3")); });
    m_trayIconMenu->addAction(tr("Kugou"), []() { QDesktopServices::openUrl(QUrl("https://www.kugou.com")); });
    m_trayIconMenu->addAction(tr("Kuwo"), []() { QDesktopServices::openUrl(QUrl("http://kuwo.cn")); });
    m_trayIconMenu->addSeparator();
    m_trayIconMenu->addAction(configAction);
    m_trayIconMenu->addAction(showHidePlayerAction);
    m_trayIconMenu->addAction(playlistManageAction);
    m_trayIconMenu->addAction(quitAction);

    m_trayIcon = new QSystemTrayIcon(this);
    m_trayIcon->setContextMenu(m_trayIconMenu);
    m_trayIcon->setIcon(QIcon(":/hannah.png"));

    m_trayIcon->show();
    connect(m_trayIcon, &QSystemTrayIcon::activated, this, &ConfigurationWindow::onSystemTrayIconActivated);
}

ConfigurationWindow::~ConfigurationWindow()
{
    delete m_nam;

    m_settings->sync();
    delete m_settings;

    delete ui;
}

void ConfigurationWindow::onSearch(const QString &s)
{
    Q_UNUSED(s);
}

void ConfigurationWindow::onOpenUrl(const QString &s)
{
    openLink(s);
}

void ConfigurationWindow::onOpenLink(const QString &s)
{
    openLink(s);
}

void ConfigurationWindow::closeEvent(QCloseEvent *event)
{
#if defined(Q_OS_MACOS)
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

void ConfigurationWindow::onOpenUrl(const QUrl &url)
{
    onApplicationMessageReceived(url.toString());
}

void ConfigurationWindow::onApplicationMessageReceived(const QString &message)
{
    const QString &u       = message;
    QString        pattern = "hannah://play";
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

void ConfigurationWindow::onUseBuiltinPlayerStateChanged(bool checked)
{
    ui->cbOutputDevices->setEnabled(checked);
}

void ConfigurationWindow::onUseExternalPlayerStateChanged(bool checked)
{
    ui->externalPlayerArguments->setEnabled(checked);
    ui->externalPlayerPath->setEnabled(checked);
    ui->externalPlayerWorkingDir->setEnabled(checked);
    ui->browseExternalPlayer->setEnabled(checked);
    ui->browseExternalPlayerWorkingDir->setEnabled(checked);

    Q_ASSERT(m_settings);
    m_settings->setValue("useExternalPlayer", checked);
    m_settings->sync();
}

void ConfigurationWindow::onExternalPlayerPathTextChanged(const QString &text)
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

void ConfigurationWindow::onBrowseExternalPlayerClicked()
{
    QString fn = QFileDialog::getOpenFileName(this, tr("External Player"));
    ui->externalPlayerPath->setText(fn);
}

void ConfigurationWindow::onExternalPlayerArgumentsTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("externalPlayerArguments", text);
    m_settings->sync();
}

void ConfigurationWindow::onExternalPlayerWorkingDirTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("externalPlayerWorkingDir", text);
    m_settings->sync();
}

void ConfigurationWindow::onBrowseExternalPlayerWorkingDirClicked()
{
    QString dir = QFileDialog::getExistingDirectory(this, tr("Working Directory"));
    ui->externalPlayerWorkingDir->setText(dir);
}

void ConfigurationWindow::initOutputDevices()
{
    auto *model = new QStandardItemModel;

    auto *item = new QStandardItem(tr("Default Driver"));
    item->setFlags(item->flags() & ~(Qt::ItemIsEnabled | Qt::ItemIsSelectable));
    item->setData("parent", Qt::AccessibleDescriptionRole);
    QFont font = item->font();
    font.setBold(true);
    item->setFont(font);
    model->appendRow(item);

    BASS_DEVICEINFO info;
    for (int deviceIndex = 1; BASS_GetDeviceInfo(deviceIndex, &info); deviceIndex++)
    {
        if (info.flags & BASS_DEVICE_ENABLED)
        {
            auto *item = new QStandardItem(QString::fromUtf8(info.name) + QString(4, QChar(' ')));
            item->setData("child", Qt::AccessibleDescriptionRole);
            model->appendRow(item);
        }
    }
#if defined(Q_OS_WIN)

    BASS_ASIO_DEVICEINFO asioinfo;
    for (int deviceIndex = 0; BASS_ASIO_GetDeviceInfo(deviceIndex, &asioinfo); deviceIndex++)
    {
        if (deviceIndex == 0)
        {
            item = new QStandardItem("ASIO");
            item->setFlags(item->flags() & ~(Qt::ItemIsEnabled | Qt::ItemIsSelectable));
            item->setData("parent", Qt::AccessibleDescriptionRole);
            font.setBold(true);
            item->setFont(font);
            model->appendRow(item);
        }
        auto *item = new QStandardItem(QString::fromUtf8(asioinfo.name) + QString(4, QChar(' ')));
        item->setData("child", Qt::AccessibleDescriptionRole);
        model->appendRow(item);
    }

    BASS_WASAPI_DEVICEINFO wasapiinfo;
    for (int deviceIndex = 0; BASS_WASAPI_GetDeviceInfo(deviceIndex, &wasapiinfo); deviceIndex++)
    {
        if (deviceIndex == 0)
        {
            item = new QStandardItem("WASAPI");
            item->setFlags(item->flags() & ~(Qt::ItemIsEnabled | Qt::ItemIsSelectable));
            item->setData("parent", Qt::AccessibleDescriptionRole);
            font.setBold(true);
            item->setFont(font);
            model->appendRow(item);
        }
        if (!(wasapiinfo.flags & BASS_DEVICE_INPUT)      // device is an output device (not input)
            && (wasapiinfo.flags & BASS_DEVICE_ENABLED)) // and it is enabled
        {
            auto *item = new QStandardItem(QString::fromUtf8(wasapiinfo.name) + QString(4, QChar(' ')));
            item->setData("child", Qt::AccessibleDescriptionRole);
            model->appendRow(item);
        }
    }

#endif
    ui->cbOutputDevices->setModel(model);
}

void ConfigurationWindow::initNetworkInterfaces()
{
    auto interfaces = QNetworkInterface::allInterfaces();
    for (const auto &networkInterface : interfaces)
    {
        if (networkInterface.type() == QNetworkInterface::Ethernet || networkInterface.type() == QNetworkInterface::Wifi ||
            networkInterface.type() == QNetworkInterface::Ppp)
        {
            ui->reverseProxyBindNetworkInterface->addItem(networkInterface.humanReadableName());
        }
    }
}
void ConfigurationWindow::restartReverseProxy()
{
    if (m_settings->value("useExternalReverseProxy", 2).toInt() == Qt::Checked)
    {
        m_externalReverseProxyRunner.stop();
        m_externalReverseProxyRunner.start();
    }
    else
    {
        m_builtinReverseProxyRunner.stop();
        m_builtinReverseProxyRunner.wait();
        m_builtinReverseProxyRunner.start();
    }
}

void ConfigurationWindow::onReverseProxyListenPortValueChanged(int port)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyListenPort", port);
    m_settings->sync();

    restartReverseProxy();
}

void ConfigurationWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyBindNetworkInterface", text);
    m_settings->sync();

    restartReverseProxy();
}

void ConfigurationWindow::onReverseProxyAutoRedirectStateChanged(int state)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyAutoRedirect", state);
    m_settings->sync();

    restartReverseProxy();
}

void ConfigurationWindow::onReverseProxyRedirectStateChanged(int state)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyRedirect", state);
    m_settings->sync();

    restartReverseProxy();
}

void ConfigurationWindow::onReverseProxyProxyTypeCurrentTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyProxyType", text);
    m_settings->sync();
    restartReverseProxy();
}

void ConfigurationWindow::onReverseProxyProxyAddressTextChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("reverseProxyProxyAddress", text);
    m_settings->sync();

    restartReverseProxy();
}

void ConfigurationWindow::openLink(const QString &text)
{
    static const QVector<QRegularExpression> patterns = {
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/(?:#\/)?discover\/toplist\?id=(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/(?:#\/)?playlist\?id=(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/(?:#\/)?my\/m\/music\/playlist\?id=(\d+))"),
        QRegularExpression(R"(^https?:\/\/y\.qq\.com\/n\/yqq\/playlist\/(\d+)\.html)"),
        QRegularExpression(R"(^https?:\/\/www\.kugou\.com\/yy\/special\/single\/(\d+)\.html)"),
        QRegularExpression(R"(^https?:\/\/(?:www\.)?kuwo\.cn\/playlist_detail\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.migu\.cn\/v3\/music\/playlist\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/(?:#\/)?song\?id=(\d+))"),
        QRegularExpression(R"(^https?:\/\/y\.qq\.com/n\/yqq\/song\/(\w+)\.html)"),
        QRegularExpression(R"(^https?:\/\/www\.kugou\.com\/song\/#hash=([0-9A-F]+))"),
        QRegularExpression(R"(^https?:\/\/(?:www\.)kuwo.cn\/play_detail\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.migu\.cn\/v3\/music\/song\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/weapi\/v1\/artist\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/(?:#\/)?artist\?id=(\d+))"),
        QRegularExpression(R"(^https?:\/\/y\.qq\.com\/n\/yqq\/singer\/(\w+)\.html)"),
        QRegularExpression(R"(^https?:\/\/(?:www\.)?kuwo\.cn\/singer_detail\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.migu\.cn\/v3\/music\/artist\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/weapi\/v1\/album\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.163\.com\/(?:#\/)?album\?id=(\d+))"),
        QRegularExpression(R"(^https?:\/\/y\.qq\.com\/n\/yqq\/album\/(\w+)\.html)"),
        QRegularExpression(R"(^https?:\/\/(?:www\.)?kuwo\.cn\/album_detail\/(\d+))"),
        QRegularExpression(R"(^https?:\/\/music\.migu\.cn\/v3\/music\/album\/(\d+))")};
    auto iter = std::find_if(patterns.begin(), patterns.end(), [&text](const auto &r) { return r.match(text).hasMatch(); });
    if (patterns.end() != iter)
    {
        handle(QString(QUrl::toPercentEncoding(text)), true);
    }
}

void ConfigurationWindow::onGlobalClipboardChanged()
{
    auto *clipboard = QGuiApplication::clipboard();
    auto  text      = clipboard->text();
    openLink(text);
}

void ConfigurationWindow::onReplyError(QNetworkReply::NetworkError code)
{
    Q_UNUSED(code);
#if !defined(QT_NO_DEBUG)
    auto *reply = qobject_cast<QNetworkReply *>(sender());
    Q_ASSERT(reply);
    qDebug() << reply->errorString();
#endif
}

void ConfigurationWindow::onReplyFinished()
{
    auto *reply = qobject_cast<QNetworkReply *>(sender());
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

void ConfigurationWindow::onReplySslErrors(const QList<QSslError> &
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

void ConfigurationWindow::onReplyReadyRead()
{
    auto *reply      = qobject_cast<QNetworkReply *>(sender());
    int   statusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
    if (statusCode >= 200 && statusCode < 300)
    {
        m_playlistContent.append(reply->readAll());
    }
}

void ConfigurationWindow::onSystemTrayIconActivated(QSystemTrayIcon::ActivationReason reason)
{
    if (reason == QSystemTrayIcon::DoubleClick)
    {
        onShowConfiguration();
    }
}

void ConfigurationWindow::onShowConfiguration()
{
    if (isHidden())
    {
        showNormal();
    }
    activateWindow();
    raise();
}

void ConfigurationWindow::onShowPlaylistManage()
{
    if (playlistManageWindow->isHidden())
    {
        playlistManageWindow->showNormal();
    }
    playlistManageWindow->activateWindow();
    playlistManageWindow->raise();
}

void ConfigurationWindow::onShowHideBuiltinPlayer()
{
    gQmlPlayer->showNormal();
}

void ConfigurationWindow::handle(const QString &url, bool needConfirm)
{
    auto player     = QDir::toNativeSeparators(ui->externalPlayerPath->text());
    auto arguments  = ui->externalPlayerArguments->text();
    auto workingDir = QDir::toNativeSeparators(ui->externalPlayerWorkingDir->text());

    QFileInfo fi(player);

    if (!fi.exists())
    {
        if (!needConfirm)
        {
            QMessageBox::critical(this, tr("Error"), tr("External player path not configured properly"));
        }
        return;
    }

    m_playlistContent.clear();

    QNetworkRequest req(QUrl::fromUserInput(QString("http://localhost:%1/m3u/generate?u=").arg(ui->reverseProxyListenPort->value()) + url));
    auto           *reply = m_nam->get(req);
    connect(reply, &QNetworkReply::finished, this, &ConfigurationWindow::onReplyFinished);
    connect(reply, &QNetworkReply::readyRead, this, &ConfigurationWindow::onReplyReadyRead);
    connect(reply, &QNetworkReply::errorOccurred, this, &ConfigurationWindow::onReplyError);
    connect(reply, &QNetworkReply::sslErrors, this, &ConfigurationWindow::onReplySslErrors);
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
    {
        return;
    }
#if defined(Q_OS_MACOS)
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
                    (const wchar_t *)player.utf16(),
                    (const wchar_t *)args.join(" ").utf16(),
                    (const wchar_t *)workingDir.utf16(),
                    SW_SHOWNORMAL);
    return;
#else
#endif
    auto argv = arguments.split(" ");
    argv << localTempPlaylist;
    argv.removeAll("");
    QProcess::startDetached(player, argv, workingDir);
}

void ConfigurationWindow::onBrowseExternalReverseProxy()
{
    QString filePath = QFileDialog::getOpenFileName(this, tr("Find Reverse Proxy"));
    if (filePath.isEmpty() || !QFile::exists(filePath))
    {
        return;
    }

    ui->edtExternalReverseProxyPath->setText(filePath);
}

void ConfigurationWindow::onUseExternalReverseProxyStateChanged(int state)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("useExternalReverseProxy", state);
    m_settings->sync();
}

void ConfigurationWindow::onExternalReverseProxyPathChanged(const QString &text)
{
    Q_ASSERT(m_settings);
    m_settings->setValue("externalReverseProxyPath", text);
    m_settings->sync();
}
