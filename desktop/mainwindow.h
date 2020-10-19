#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT
    
public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:
    void on_useExternalPlayer_stateChanged(int arg1);

    void on_externalPlayerPath_textChanged(const QString &arg1);

    void on_browseExternalPlayer_clicked();

    void on_externalPlayerArguments_textChanged(const QString &arg1);

    void on_externalPlayerWorkingDir_textChanged(const QString &arg1);

    void on_browseExternalPlayerWorkingDir_clicked();

    void on_reverseProxyListenPort_valueChanged(int arg1);

    void on_reverseProxyBindNetworkInterface_currentTextChanged(const QString &arg1);

    void on_reverseProxyAutoRedirect_stateChanged(int arg1);

    void on_reverseProxyRedirect_stateChanged(int arg1);

    void on_reverseProxyProxyType_currentTextChanged(const QString &arg1);

    void on_reverseProxyProxyAddress_textChanged(const QString &arg1);

private:
    Ui::MainWindow *ui;
};
#endif // MAINWINDOW_H
