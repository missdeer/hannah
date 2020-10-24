#include <QCoreApplication>
#include <QDir>
#include <QSettings>

int main(int argc, char *argv[])
{
    QCoreApplication a(argc, argv);

#if defined(Q_OS_WIN)
    QSettings mxKey("HKEY_CLASSES_ROOT\\hannah", QSettings::NativeFormat);
    mxKey.setValue(".", "URL:hannah Protocol");
    mxKey.setValue("URL Protocol", "");
    mxKey.sync();

    QSettings mxOpenKey("HKEY_CLASSES_ROOT\\hannah\\shell\\open\\command", QSettings::NativeFormat);
    mxOpenKey.setValue(".", QChar('"') + QDir::toNativeSeparators(QCoreApplication::applicationDirPath()) + QString("\\Hannah.exe\" \"%1\""));
    mxKey.sync();
#endif

    return 0;
}
