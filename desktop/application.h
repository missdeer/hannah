#ifndef APPLICATION_H
#define APPLICATION_H

#include <QApplication>
#include <QUrl>

QT_BEGIN_NAMESPACE
class QEvent;
QT_END_NAMESPACE

class Application : public QApplication
{
    Q_OBJECT

public:
    Application(int &argc, char **argv) : QApplication(argc, argv) {}

    bool event(QEvent *event) override;

signals:
    void openUrl(QUrl);
};

#endif // APPLICATION_H
