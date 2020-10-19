#include <QApplication>
#include <QCommandLineParser>
#include <QFileOpenEvent>
#include <QMessageBox>
#include <QtCore>

#include "mainwindow.h"

#if defined(Q_OS_MAC)
#    include "application.h"
#endif

int main(int argc, char *argv[])
{
#if defined(Q_OS_MAC)
    Application a(argc, argv);
    MainWindow  w;
    a.connect(&a, &Application::openUrl, &w, &MainWindow::onOpenUrl);
#else
    QApplication a(argc, argv);
    QCoreApplication::setApplicationName("Hannah");
    QCoreApplication::setApplicationVersion("1.0");

    QCommandLineParser parser;
    parser.setApplicationDescription("Hannah");
    parser.addHelpOption();
    parser.addVersionOption();

    parser.process(a);

    const QStringList args = parser.positionalArguments();

    if (args.length() > 0)
        QMessageBox::information(nullptr, "arguments", args.join(" "));
    MainWindow w;
#endif

    return a.exec();
}
