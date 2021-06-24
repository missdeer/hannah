#include <QApplication>
#include <QCommandLineParser>
#include <QDir>
#include <QFileOpenEvent>
#include <QMessageBox>
#include <QSettings>
#include <QtCore>

#include "bass.h"
#include "configurationwindow.h"
#include "playlistmanagewindow.h"
#include "shadowplayer.h"
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
    QCoreApplication::setApplicationVersion("1.0");

    QTranslator translator;
    QTranslator qtTranslator;

#if defined(Q_OS_WIN) && (QT_VERSION < QT_VERSION_CHECK(6, 0, 0))
    QWindowsWindowFunctions::setWindowActivationBehavior(QWindowsWindowFunctions::AlwaysActivateWindow);
#endif

    // check the correct BASS was loaded
    if (HIWORD(BASS_GetVersion()) != BASSVERSION)
    {
        QMessageBox::critical(0, QObject::tr("Critical Error"), QObject::tr("An incorrect version of BASS.DLL was loaded"));
        return -1;
    }
    BASS_SetConfig(BASS_CONFIG_UNICODE, TRUE);

#if defined(Q_OS_MACOS)
    Application a(argc, argv);
    i18n(translator, qtTranslator);
    ConfigurationWindow  w;
    w.connect(&a, &Application::openUrl, &w, qOverload<QUrl>(&ConfigurationWindow::onOpenUrl));

    configurationWindow = &w;

    void registerHannahService();
    registerHannahService();
#else
    QtSingleApplication a(argc, argv);

    QCommandLineParser parser;
    parser.setApplicationDescription("Hannah");
    parser.addHelpOption();
    parser.addVersionOption();

    parser.process(a);

    const QStringList args = parser.positionalArguments();

    if (a.isRunning())
    {
        if (args.length() > 0)
        {
            a.sendMessage(args.join("~"));
        }
        return 0;
    }

    i18n(translator, qtTranslator);
    ConfigurationWindow w;
    w.connect(&a, &QtSingleApplication::messageReceived, &w, &ConfigurationWindow::onApplicationMessageReceived);
    configurationWindow = &w;
    if (args.length() > 0)
    {
        w.onApplicationMessageReceived(args.join("~"));
    }
#    if defined(Q_OS_WIN)
    else
    {
        QSettings mxKey("HKEY_CLASSES_ROOT\\hannah", QSettings::NativeFormat);
        QString   v1 = mxKey.value(".").toString();
        QSettings mxOpenKey("HKEY_CLASSES_ROOT\\hannah\\shell\\open\\command", QSettings::NativeFormat);
        QString   v2 = mxOpenKey.value(".").toString();

        if (v1 != "URL:hannah Protocol" ||
            v2 != QChar('"') + QDir::toNativeSeparators(QCoreApplication::applicationFilePath()) + QString("\" \"%1\""))
        {
            QString cmd = QDir::toNativeSeparators(QCoreApplication::applicationDirPath() + "/registerProtocolHandler.exe");
            ::ShellExecuteW(nullptr,
                            L"open",
                            cmd.toStdWString().c_str(),
                            nullptr,
                            QDir::toNativeSeparators(QCoreApplication::applicationDirPath()).toStdWString().c_str(),
                            SW_SHOWNORMAL);
        }
    }
#    endif
#endif

    ShadowPlayer player;
    shadowPlayer = &player;

    PlaylistManageWindow pmw;
    playlistManageWindow = &pmw;

    a.setQuitOnLastWindowClosed(false);
    return a.exec();
}
