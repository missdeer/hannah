#include <QQmlEngine>
#include <QQmlContext>
#include "qmldialog.h"
#include "ui_qmldialog.h"

QmlDialog::QmlDialog(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::QmlDialog)
{
    ui->setupUi(this);
    ui->quickWidget->setAttribute(Qt::WA_TranslucentBackground, true);
    ui->quickWidget->setAttribute(Qt::WA_AlwaysStackOnTop, true);
    ui->quickWidget->setClearColor(Qt::transparent);
    ui->quickWidget->setResizeMode(QQuickWidget::SizeRootObjectToView);
    ui->quickWidget->engine()->addImportPath("qrc:/rc/qml");
}

QmlDialog::~QmlDialog()
{
    delete ui;
}

void QmlDialog::doClose()
{
    QDialog::accept();
}

void QmlDialog::loadQml(const QUrl &u)
{
    ui->quickWidget->setSource(u);
}

QQmlEngine *QmlDialog::engine()
{
    return ui->quickWidget->engine();
}

QQmlContext *QmlDialog::context()
{
    return ui->quickWidget->engine()->rootContext();
}
