#include <QApplication>
#include <QCommandLineParser>
#include <QDir>
#include <QFileOpenEvent>
#include <QMessageBox>
#include <QSettings>
#include <QtCore>

#include "mainwindow.h"

#if defined(Q_OS_MAC)
#    include "application.h"
#else
#    if defined(Q_OS_WIN)
#        include <Windows.h>
#        include <shellapi.h>
#        include <tchar.h>
#    endif
#    include "qtsingleapplication.h"
#endif

void i18n()
{
    QString     locale = "zh_CN";
    QTranslator translator;
    QTranslator qtTranslator;

    // main application and dynamic linked library locale
#if defined(Q_OS_MAC)
    QString localeDirPath = QApplication::applicationDirPath() + "/../Resources/translations";
#else
    QString localeDirPath = QApplication::applicationDirPath() + "/translations";
    if (!QDir(localeDirPath).exists())
    {
        localeDirPath = QApplication::applicationDirPath() + "/../translations";
    }
#endif

    if (!translator.load("Hannah_" + locale, localeDirPath))
    {
        qDebug() << "loading Hannah" << locale << " from " << localeDirPath << " failed";
    }
    else
    {
        qDebug() << "loading Hannah" << locale << " from " << localeDirPath << " success";
        if (!QApplication::installTranslator(&translator))
        {
            qDebug() << "installing translator failed ";
        }
    }

    // qt locale
    if (!qtTranslator.load("qt_" + locale, localeDirPath))
    {
        qDebug() << "loading qt" << locale << " from " << localeDirPath << " failed";
    }
    else
    {
        qDebug() << "loading qt" << locale << " from " << localeDirPath << " success";
        if (!QApplication::installTranslator(&qtTranslator))
        {
            qDebug() << "installing qt translator failed ";
        }
    }
}

int main(int argc, char *argv[])
{
#if defined(Q_OS_MAC)
    Application a(argc, argv);
    i18n();
    MainWindow  w;
    w.connect(&a, &Application::openUrl, &w, &MainWindow::onOpenUrl);
#else
    QtSingleApplication a(argc, argv);
    QCoreApplication::setApplicationName("Hannah");
    QCoreApplication::setApplicationVersion("1.0");

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

    i18n();
    MainWindow w;
    w.connect(&a, &QtSingleApplication::messageReceived, &w, &MainWindow::onApplicationMessageReceived);
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

    a.setQuitOnLastWindowClosed(false);
    return a.exec();
}
