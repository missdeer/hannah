#ifndef QMLDIALOG_H
#define QMLDIALOG_H

#include <QDialog>

namespace Ui {
class QmlDialog;
}
class QQmlEngine;
class QQmlContext;

class QmlDialog : public QDialog
{
    Q_OBJECT

public:
    explicit QmlDialog(QWidget *parent = nullptr);
    void loadQml(const QUrl& u);
    QQmlEngine* engine();
    QQmlContext* context();
    ~QmlDialog();

public slots:
    void doClose();

private:
    Ui::QmlDialog *ui;
};

#endif // QMLDIALOG_H
