#include "mainwindow.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_useExternalPlayer_stateChanged(int arg1) {}

void MainWindow::on_externalPlayerPath_textChanged(const QString &arg1) {}

void MainWindow::on_browseExternalPlayer_clicked() {}

void MainWindow::on_externalPlayerArguments_textChanged(const QString &arg1) {}

void MainWindow::on_externalPlayerWorkingDir_textChanged(const QString &arg1) {}

void MainWindow::on_browseExternalPlayerWorkingDir_clicked() {}

void MainWindow::on_reverseProxyListenPort_valueChanged(int arg1) {}

void MainWindow::on_reverseProxyBindNetworkInterface_currentTextChanged(const QString &arg1) {}

void MainWindow::on_reverseProxyAutoRedirect_stateChanged(int arg1) {}

void MainWindow::on_reverseProxyRedirect_stateChanged(int arg1) {}

void MainWindow::on_reverseProxyProxyType_currentTextChanged(const QString &arg1) {}

void MainWindow::on_reverseProxyProxyAddress_textChanged(const QString &arg1) {}
