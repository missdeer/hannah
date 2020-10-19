#include <QApplication>
#include <QCommandLineParser>
#include <QFileOpenEvent>
#include <QMessageBox>
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
    w.connect(&a, &Application::openUrl, &w, &MainWindow::onOpenUrl);
    w.connect(&a, &Application::messageReceived, &w, &MainWindow::onApplicationMessageReceived);
    if (a.isRunning())
    {
        a.sendMessage("running");
        a.exit();
    }
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
#endif

    return a.exec();
}
