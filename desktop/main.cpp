#include <QApplication>
#include <QCommandLineParser>
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

int main(int argc, char *argv[])
{
#if defined(Q_OS_MAC)
    Application a(argc, argv);
    MainWindow  w;
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
    MainWindow w;
    w.connect(&a, &QtSingleApplication::messageReceived, &w, &MainWindow::onApplicationMessageReceived);

#    if defined(Q_OS_WIN)
    QSettings mxKey("HKEY_CLASSES_ROOT\\hannah", QSettings::NativeFormat);
    mxKey.setValue(".", "URL:hannah Protocol");
    mxKey.setValue("URL Protocol", "");

    QSettings mxOpenKey("HKEY_CLASSES_ROOT\\foo\\shell\\open\\command", QSettings::NativeFormat);
    QString   cmdLine = QString("\"%1\" \"%%1\"").arg(QDir::toNativeSeparators(QCoreApplication::applicationFilePath()));
    mxOpenKey.setValue(".", cmdLine);
#    endif
#endif

    return a.exec();
}
