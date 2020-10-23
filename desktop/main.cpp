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
        qDebug() << "loading " << locale << " from " << localeDirPath << " failed";
    }
    else
    {
        qDebug() << "loading " << locale << " from " << localeDirPath << " success";
        if (!QApplication::installTranslator(&translator))
        {
            qDebug() << "installing translator failed ";
        }
    }

    // qt locale
    if (!qtTranslator.load("qt_" + locale, localeDirPath))
    {
        qDebug() << "loading " << locale << " from " << localeDirPath << " failed";
    }
    else
    {
        qDebug() << "loading " << locale << " from " << localeDirPath << " success";
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

    if (args.length() > 0)
    {
        QMessageBox::information(nullptr, "arguments", args.join(" "));
        if (a.isRunning())
        {
            a.sendMessage(args.join("~"));
            a.exit();
        }
    }
    i18n();
    MainWindow w;
    w.connect(&a, &QtSingleApplication::messageReceived, &w, &MainWindow::onApplicationMessageReceived);

#    if defined(Q_OS_WIN)
    QSettings mxKey("HKEY_CLASSES_ROOT\\hannah", QSettings::NativeFormat);
    mxKey.setValue(".", "URL:hannah Protocol");
    mxKey.setValue("URL Protocol", "");

    QSettings mxOpenKey("HKEY_CLASSES_ROOT\\foo\\shell\\open\\command", QSettings::NativeFormat);
    mxOpenKey.setValue(".", QChar('"') + QDir::toNativeSeparators(QCoreApplication::applicationFilePath()) + QString("\" \"%1\""));
#    endif
#endif

    return a.exec();
}
