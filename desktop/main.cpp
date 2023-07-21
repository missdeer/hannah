#include <QApplication>
#include <QCommandLineParser>
#include <QDir>
#include <QFileOpenEvent>
#include <QMessageBox>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QSettings>
#include <QtCore>

#include "bass.h"
#include "bassplayer.h"
#include "configurationwindow.h"
#include "playlistmanagewindow.h"
#include "qmlplayer.h"
#if defined(Q_OS_MACOS)
#    include "application.h"
#    include "serviceslots.h"
#else
#    if defined(Q_OS_WIN)
#        include <Windows.h>
#        include <shellapi.h>
#        include <tchar.h>
#        if (QT_VERSION < QT_VERSION_CHECK(6, 0, 0))
#            include <QtPlatformHeaders/QWindowsWindowFunctions>
#        endif
#    endif
#    include "qtsingleapplication.h"
#endif

#if defined(Q_OS_MACOS)

void serviceSearch(const QString &s)
{
    if (configurationWindow)
    {
        configurationWindow->onSearch(s);
    }
}
void serviceOpenUrl(const QString &s)
{
    if (configurationWindow)
    {
        configurationWindow->onSearch(s);
    }
}

void serviceOpenLink(const QString &s)
{
    if (configurationWindow)
    {
        configurationWindow->onOpenLink(s);
    }
}

void serviceAppendToPlaylist(const QStringList &s)
{
    if (playlistManageWindow)
    {
        playlistManageWindow->onAppendToPlaylist(s);
    }
}

void serviceClearAndAddToPlaylist(const QStringList &s)
{
    if (playlistManageWindow)
    {
        playlistManageWindow->onClearAndAddToPlaylist(s);
    }
}

void serviceAppendToPlaylistFile(const QStringList &s)
{
    if (playlistManageWindow)
    {
        playlistManageWindow->onAppendToPlaylistFile(s);
    }
}

void serviceClearAndAddToPlaylistFile(const QStringList &s)
{
    if (playlistManageWindow)
    {
        playlistManageWindow->onClearAndAddToPlaylistFile(s);
    }
}

#endif

void i18n(QTranslator &translator, QTranslator &qtTranslator)
{
    QString locale = "zh_CN";

    // main application and dynamic linked library locale
    QString localeDirPath = QCoreApplication::applicationDirPath() +
#if defined(Q_OS_MACOS)
                            "/../Resources/translations";
#else
                            "/translations";
#endif

    if (translator.load("Hannah_" + locale, localeDirPath))
    {
        qDebug() << "loading Hannah" << locale << " from " << localeDirPath << " success";
        if (QCoreApplication::installTranslator(&translator))
        {
            qDebug() << "installing translator success ";
        }
    }

    if (qtTranslator.load("qt_" + locale, localeDirPath))
    {
        qDebug() << "loading qt" << locale << " from " << localeDirPath << " success";
        if (QCoreApplication::installTranslator(&qtTranslator))
        {
            qDebug() << "installing qt translator success ";
        }
    }
}

int main(int argc, char *argv[])
{
#if (QT_VERSION < QT_VERSION_CHECK(6, 0, 0))
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QCoreApplication::setAttribute(Qt::AA_UseHighDpiPixmaps);
#endif
    QCoreApplication::setApplicationName("Hannah");
    QCoreApplication::setApplicationVersion("2.0");
    QCoreApplication::setOrganizationName("Minidump.Info");
    QCoreApplication::setOrganizationDomain("Minidump.Info");

    QTranslator translator;
    QTranslator qtTranslator;

    // check the correct BASS was loaded
    if (HIWORD(BASS_GetVersion()) != BASSVERSION)
    {
        QMessageBox::critical(nullptr, QObject::tr("Critical Error"), QObject::tr("An incorrect version of BASS.DLL was loaded"));
        return -1;
    }
    BASS_SetConfig(BASS_CONFIG_UNICODE, TRUE);

    QIcon::setThemeName("musicplayer");

#if defined(Q_OS_MACOS)
    Application app(argc, argv);
    i18n(translator, qtTranslator);
    ConfigurationWindow  configWin;
    configWin.connect(&app, &Application::openUrl, &configWin, qOverload<const QUrl &>(&ConfigurationWindow::onOpenUrl));

    configurationWindow = &configWin;

    void registerHannahService();
    registerHannahService();
#else
    QtSingleApplication app(argc, argv);

    QCommandLineParser parser;
    parser.setApplicationDescription("Hannah");
    parser.addHelpOption();
    parser.addVersionOption();

    parser.process(app);

    const QStringList args = parser.positionalArguments();

    if (app.isRunning())
    {
        if (args.length() > 0)
        {
            app.sendMessage(args.join("~"));
        }
        return 0;
    }

#    if defined(Q_OS_WIN) && (QT_VERSION < QT_VERSION_CHECK(6, 0, 0))
    QWindowsWindowFunctions::setWindowActivationBehavior(QWindowsWindowFunctions::AlwaysActivateWindow);
#    endif

    i18n(translator, qtTranslator);
    ConfigurationWindow configWin;
    QObject::connect(&app, &QtSingleApplication::messageReceived, &configWin, &ConfigurationWindow::onApplicationMessageReceived);
    configurationWindow = &configWin;
    if (args.length() > 0)
    {
        configWin.onApplicationMessageReceived(args.join("~"));
    }
#    if defined(Q_OS_WIN)
    else
    {
        QSettings mxKey("HKEY_CLASSES_ROOT\\hannah", QSettings::NativeFormat);
        QString   v1 = mxKey.value(".").toString();
        QSettings mxOpenKey(R"(HKEY_CLASSES_ROOT\hannah\shell\open\command)", QSettings::NativeFormat);
        QString   v2 = mxOpenKey.value(".").toString();

        if (v1 != "URL:hannah Protocol" ||
            v2 != QChar('"') + QDir::toNativeSeparators(QCoreApplication::applicationFilePath()) + QStringLiteral(R"(" "%1")"))
        {
            auto cmd = QDir::toNativeSeparators(QCoreApplication::applicationDirPath() + "/registerProtocolHandler.exe");
            auto workingDir = QDir::toNativeSeparators(QCoreApplication::applicationDirPath());
            SHELLEXECUTEINFO execInfo;
            ZeroMemory(&execInfo, sizeof(execInfo));
            execInfo.lpFile = (const wchar_t *)cmd.utf16();
            execInfo.lpDirectory = (const wchar_t *)workingDir.utf16();
            execInfo.cbSize = sizeof(execInfo);
            execInfo.lpVerb = L"runas";
            execInfo.fMask = SEE_MASK_NOCLOSEPROCESS | SEE_MASK_FLAG_NO_UI;
            execInfo.nShow = SW_HIDE;
            ShellExecuteEx(&execInfo);
        }
    }
#    endif
#endif

    gQmlApplicationEngine = new QQmlApplicationEngine;
    QQmlContext *context  = gQmlApplicationEngine->rootContext();

    gBassPlayer = new BassPlayer;
    gQmlPlayer  = new QmlPlayer;
    context->setContextProperty("playerCore", gQmlPlayer);
    gQmlApplicationEngine->load(QUrl("qrc:/rc/qml/musicplayer.qml"));

    gQmlPlayer->setTaskbarButtonWindow();

    PlaylistManageWindow pmw;
    playlistManageWindow = &pmw;

    QtSingleApplication::setQuitOnLastWindowClosed(false);
    return QtSingleApplication::exec();
}
