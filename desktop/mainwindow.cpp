#include <QClipboard>
#include <QCloseEvent>
#include <QCoreApplication>
#include <QFileDialog>
#include <QMenu>
#include <QMessageBox>
#include <QNetworkInterface>
#include <QSettings>
#include <QSystemTrayIcon>

#include "mainwindow.h"

#include "librp.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    settings = new QSettings(QSettings::IniFormat, QSettings::UserScope, "Minidump.info");
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
    ui->reverseProxyBindNetworkInterface->setCurrentText(settings->value("reverseProxyBindNetworkInterface").toString());
    ui->reverseProxyProxyType->setCurrentText(settings->value("reverseProxyProxyType").toString());
    ui->reverseProxyProxyAddress->setText(settings->value("reverseProxyProxyAddress").toString());
    auto state = settings->value("useExternalPlayer").toInt(&ok);
    if (ok)
        ui->useExternalPlayer->setCheckState(Qt::CheckState(state));
    state = settings->value("reverseProxyAutoRedirect").toInt(&ok);
    if (ok)
        ui->reverseProxyAutoRedirect->setCheckState(Qt::CheckState(state));
    state = settings->value("reverseProxyRedirect").toInt(&ok);
    if (ok)
        ui->reverseProxyRedirect->setCheckState(Qt::CheckState(state));
    auto port = settings->value("reverseProxyListenPort").toInt(&ok);
    if (ok)
        ui->reverseProxyListenPort->setValue(port);

    connect(ui->reverseProxyBindNetworkInterface,
            &QComboBox::currentTextChanged,
            this,
            &MainWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged);

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
    QMessageBox::information(this, tr("Open URL"), url.toString());
}

void MainWindow::onApplicationMessageReceived(const QString &message)
{
    QMessageBox::information(this, tr("application message received"), message);
}

void MainWindow::on_useExternalPlayer_stateChanged(int state)
{
    ui->externalPlayerArguments->setEnabled(state == Qt::Checked);
    ui->externalPlayerPath->setEnabled(state == Qt::Checked);
    ui->browseExternalPlayer->setEnabled(state == Qt::Checked);
    ui->browseExternalPlayerWorkingDir->setEnabled(state == Qt::Checked);
    ui->externalPlayerWorkingDir->setEnabled(state == Qt::Checked);

    Q_ASSERT(settings);
    settings->setValue("useExternalPlayer", state);
}

void MainWindow::on_externalPlayerPath_textChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("externalPlayerPath", text);
}

void MainWindow::on_browseExternalPlayer_clicked()
{
    QString fn = QFileDialog::getOpenFileName(this, tr("External Player"));
    ui->externalPlayerPath->setText(fn);
}

void MainWindow::on_externalPlayerArguments_textChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("externalPlayerArguments", text);
}

void MainWindow::on_externalPlayerWorkingDir_textChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("externalPlayerWorkingDir", text);
}

void MainWindow::on_browseExternalPlayerWorkingDir_clicked()
{
    QString dir = QFileDialog::getExistingDirectory(this, tr("Working Directory"));
    ui->externalPlayerWorkingDir->setText(dir);
}

void MainWindow::on_reverseProxyListenPort_valueChanged(int port)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyListenPort", port);
}

void MainWindow::onReverseProxyBindNetworkInterfaceCurrentTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyBindNetworkInterface", text);
}

void MainWindow::on_reverseProxyAutoRedirect_stateChanged(int state)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyAutoRedirect", state);
}

void MainWindow::on_reverseProxyRedirect_stateChanged(int state)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyRedirect", state);
}

void MainWindow::on_reverseProxyProxyType_currentTextChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyProxyType", text);
}

void MainWindow::on_reverseProxyProxyAddress_textChanged(const QString &text)
{
    Q_ASSERT(settings);
    settings->setValue("reverseProxyProxyAddress", text);
}

void MainWindow::onGlobalClipboardChanged()
{
    QClipboard *clipboard = QGuiApplication::clipboard();
    QString     text      = clipboard->text();
    // check text pattern
}
