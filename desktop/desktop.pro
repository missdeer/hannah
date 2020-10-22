QT       += core gui network widgets

CONFIG += c++17

TARGET = Hannah

INCLUDEPATH += $$PWD/../lib/reverseProxy
LIBS += -L$$PWD/../lib/reverseProxy -lrp
# You can make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    main.cpp \
    mainwindow.cpp \
    qtlocalpeer.cpp \
    qtlockedfile.cpp \
    qtlockedfile_unix.cpp \
    qtlockedfile_win.cpp

HEADERS += \
    mainwindow.h \
    qtlocalpeer.h \
    qtlockedfile.h

FORMS += \
    mainwindow.ui
    
macx : {
    HEADERS += \
        application.h 
    SOURCES += \
        application.cpp 
    QMAKE_INFO_PLIST = macInfo.plist
    ICON = hannah.icns
    icon.path = $$PWD
    INSTALLS += icon
} else : {
    HEADERS += \
        qtsinglecoreapplication.h  \
        qtsingleapplication.h 
    SOURCES += \
        qtsinglecoreapplication.cpp  \
        qtsingleapplication.cpp 
}

# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target

RESOURCES += \
    hannah.qrc
