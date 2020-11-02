QT       += core gui network widgets

CONFIG += c++17

TARGET = Hannah

include($$PWD/3rdparty/sqlite3/sqlite3.pri)

INCLUDEPATH += $$PWD/../lib/reverseProxy
# You can make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    main.cpp \
    configurationwindow.cpp \
    playlistmanagewindow.cpp \
    playlistmodel.cpp \
    qtlocalpeer.cpp \
    qtlockedfile.cpp \
    qtlockedfile_unix.cpp \
    qtlockedfile_win.cpp \
    songlistmodel.cpp \
    sqlite3helper.cpp

HEADERS += \
    configurationwindow.h \
    playlistmanagewindow.h \
    playlistmodel.h \
    qtlocalpeer.h \
    qtlockedfile.h \
    songlistmodel.h \
    sqlite3helper.h

FORMS += \
    configurationwindow.ui \
    playlistmanagewindow.ui
    
RC_FILE = Hannah.rc

CODECFORTR      = UTF-8
CODECFORSRC     = UTF-8
TRANSLATIONS    = $$PWD/translations/Hannah_zh_CN.ts

isEmpty(QMAKE_LUPDATE) {
    QMAKE_LUPDATE = $$shell_path($$[QT_INSTALL_BINS]\lupdate)
}

isEmpty(QMAKE_LRELEASE) {
    QMAKE_LRELEASE = $$shell_path($$[QT_INSTALL_BINS]\lrelease)
}

lupdate.commands = $$QMAKE_LUPDATE -no-obsolete $$shell_path($$PWD/desktop.pro)
lupdates.depends = $$SOURCES $$HEADERS $$FORMS $$TRANSLATIONS
lrelease.commands = $$QMAKE_LRELEASE $$shell_path($$PWD/desktop.pro)
lrelease.depends = lupdate
translate.depends = lrelease
QMAKE_EXTRA_TARGETS += lupdate lrelease translate qti18n 
POST_TARGETDEPS += translate qti18n

win32: {
    CONFIG(release, debug|release) : {
        WINDEPLOYQT = $$shell_path($$[QT_INSTALL_BINS]/windeployqt.exe)
        QMAKE_EXTRA_TARGETS += mkdir
        
        qti18n.depends = translate
        win32-*g++*: {
            translate.commands = '$(COPY_FILE) $$shell_path($$PWD/translations/*.qm) $$shell_path($$OUT_PWD/release/translations/)'
            qti18n.commands = '$(COPY_FILE) $$shell_path($$[QT_INSTALL_BINS]/../share/qt5/translations/qt_zh_CN.qm) $$shell_path($$OUT_PWD/release/translations/qt_zh_CN.qm)'
        } else: {
            mkdir.commands = '$(CHK_DIR_EXISTS) $$shell_path($$OUT_PWD/release/translations/) $(MKDIR) $$shell_path($$OUT_PWD/release/translations/)'
            translate.depends += mkdir
            translate.commands = '$(CHK_DIR_EXISTS) $$shell_path($$PWD/translations/Hannah_zh_CN.qm) $(COPY_FILE) $$shell_path($$PWD/translations/*.qm) $$shell_path($$OUT_PWD/release/translations/)'
            qti18n.commands = '$(COPY_FILE) $$shell_path($$[QT_INSTALL_BINS]/../translations/qt_zh_CN.qm) $$shell_path($$OUT_PWD/release/translations/)'
        }
    }
}

LIBS += -L$$PWD/../lib/reverseProxy -lrp

macx : {
    HEADERS += \
        serviceslots.h \
        application.h
    SOURCES += \
        application.cpp 
    OBJECTIVE_HEADERS += \
        service.h
    OBJECTIVE_SOURCES += \
        service.mm
    DESTDIR = $$OUT_PWD
    QMAKE_INFO_PLIST = macInfo.plist
    ICON = hannah.icns
    icon.path = $$PWD
    INSTALLS += icon
    LIBS += -framework Security 
    
    CONFIG(release, debug|release) : {
        MACDEPLOYQT = $$[QT_INSTALL_BINS]/macdeployqt
    
        translate.depends = lrelease
        translate.files = $$system("find $${PWD}/translations -name '*.qm' ")
        translate.path = Contents/Resources/translations/
        translate.commands = '$(COPY_DIR) $$shell_path($${PWD}/translations) $$shell_path($${DESTDIR}/$${TARGET}.app/Contents/Resources/)'
    
        qti18n.depends = translate
        qti18n.commands = '$(COPY_FILE) $$shell_path($$[QT_INSTALL_BINS]/../translations/qt_zh_CN.qm) $$shell_path($${DESTDIR}/$${TARGET}.app/Contents/Resources/translations/qt_zh_CN.qm)'
    
        QMAKE_BUNDLE_DATA += translate qti18n 
    
        deploy.commands += $$MACDEPLOYQT \"$${DESTDIR}/$${TARGET}.app\"
    
        deploy_appstore.depends += deploy
        deploy_appstore.commands += $$MACDEPLOYQT \"$${DESTDIR}/$${TARGET}.app\" -appstore-compliant
    
        deploy_webengine.depends += deploy_appstore
        deploy_webengine.commands += $$MACDEPLOYQT \"$${DESTDIR}/$${TARGET}.app/Contents/Frameworks/QtWebEngineCore.framework/Helpers/QtWebEngineProcess.app\"
    
        fixdeploy.depends += deploy_webengine
        fixdeploy.commands += $$PWD/../macdeploy/macdeploy \"$${DESTDIR}/$${TARGET}.app\"
    
        APPCERT = Developer ID Application: Fan Yang (Y73SBCN2CG)
        INSTALLERCERT = 3rd Party Mac Developer Installer: Fan Yang (Y73SBCN2CG)
        BUNDLEID = info.minidump.hannah
    
        codesign.depends += fixdeploy
        codesign.commands = codesign -s \"$${APPCERT}\" -v -f --timestamp=none --deep \"$${DESTDIR}/$${TARGET}.app\"
    
        makedmg.depends += codesign
        makedmg.commands = hdiutil create -srcfolder \"$${DESTDIR}/$${TARGET}.app\" -volname \"$${TARGET}\" -format UDBZ \"$${DESTDIR}/$${TARGET}.dmg\" -ov -scrub -stretch 2g
    
        QMAKE_EXTRA_TARGETS += deploy deploy_webengine deploy_appstore fixdeploy codesign makedmg 
    }
    
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
