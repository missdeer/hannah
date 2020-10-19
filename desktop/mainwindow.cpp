#include <QClipboard>
#include <QCloseEvent>
#include <QCoreApplication>
#include <QFileDialog>
#include <QMenu>
#include <QMessageBox>
#include <QSystemTrayIcon>

#include "mainwindow.h"

#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    QClipboard *clipboard = QGuiApplication::clipboard();
    connect(clipboard, &QClipboard::dataChanged, this, &MainWindow::onGlobalClipboardChanged);

    auto configAction = new QAction(tr("&Configuration"), this);
    connect(configAction, &QAction::triggered, this, &MainWindow::showNormal);

    auto quitAction = new QAction(tr("&Quit"), this);
    connect(quitAction, &QAction::triggered, qApp, &QCoreApplication::quit);

    trayIconMenu = new QMenu(this);
    trayIconMenu->addAction(configAction);
    trayIconMenu->addSeparator();
    trayIconMenu->addAction(quitAction);

    trayIcon = new QSystemTrayIcon(this);
    trayIcon->setContextMenu(trayIconMenu);
    trayIcon->setIcon(QIcon(":/hannah.png"));

    trayIcon->show();
}

MainWindow::~MainWindow()
{
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
    QMessageBox::information(this, tr("Open URL"), url.toString());
}

void MainWindow::on_useExternalPlayer_stateChanged(int arg1) {}

void MainWindow::on_externalPlayerPath_textChanged(const QString &arg1) {}

void MainWindow::on_browseExternalPlayer_clicked()
{
    QString fn = QFileDialog::getOpenFileName(this, tr("External Player"));
    ui->externalPlayerPath->setText(fn);
}

void MainWindow::on_externalPlayerArguments_textChanged(const QString &arg1) {}

void MainWindow::on_externalPlayerWorkingDir_textChanged(const QString &arg1) {}

void MainWindow::on_browseExternalPlayerWorkingDir_clicked()
{
    QString dir = QFileDialog::getExistingDirectory(this, tr("Working Directory"));
    ui->externalPlayerWorkingDir->setText(dir);
}

void MainWindow::on_reverseProxyListenPort_valueChanged(int arg1) {}

void MainWindow::on_reverseProxyBindNetworkInterface_currentTextChanged(const QString &arg1) {}

void MainWindow::on_reverseProxyAutoRedirect_stateChanged(int arg1) {}

void MainWindow::on_reverseProxyRedirect_stateChanged(int arg1) {}

void MainWindow::on_reverseProxyProxyType_currentTextChanged(const QString &arg1) {}

void MainWindow::on_reverseProxyProxyAddress_textChanged(const QString &arg1) {}

void MainWindow::onGlobalClipboardChanged()
{
    QClipboard *clipboard = QGuiApplication::clipboard();
    QString     text      = clipboard->text();
    // check text pattern
}
