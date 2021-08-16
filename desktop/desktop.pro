QT       += core gui network widgets quick quickwidgets quickcontrols2
greaterThan(QT_MAJOR_VERSION, 5): QT += core5compat

CONFIG += c++17 sdk_no_version_check

TARGET = Hannah

include($$PWD/3rdparty/sqlite3/sqlite3.pri)

DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x050F00 _CRT_SECURE_NO_WARNINGS=

INCLUDEPATH += $$PWD/../lib/reverseProxy $$PWD $$PWD/player $$PWD/3rdparty/bass/include
# You can make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    $$PWD/main.cpp \
    $$PWD/configurationwindow.cpp \
    $$PWD/player/fftdisplay.cpp \
    $$PWD/player/lrcbar.cpp \
    $$PWD/player/lyrics.cpp \
    $$PWD/player/osd.cpp \
    $$PWD/player/player.cpp \
    $$PWD/player/shadowlabel.cpp \
    $$PWD/player/spslider.cpp \
    $$PWD/playlistmanagewindow.cpp \
    $$PWD/playlistmodel.cpp \
    $$PWD/qtlocalpeer.cpp \
    $$PWD/qtlockedfile.cpp \
    $$PWD/qtlockedfile_unix.cpp \
    $$PWD/qtlockedfile_win.cpp \
    $$PWD/songlistmodel.cpp \
    $$PWD/sqlite3helper.cpp \
    $$PWD/comboboxdelegate.cpp \
    $$PWD/player/playlist.cpp \
    $$PWD/player/shadowplayer.cpp \
    $$PWD/parser.cpp \
    $$PWD/parsercsv.cpp \
    $$PWD/parserm3u.cpp \
    $$PWD/parserpls.cpp \
    $$PWD/trackfile.cpp \
    $$PWD/qmldialog.cpp \
    $$PWD/player/qmlplayer.cpp

HEADERS += \
    $$PWD/configurationwindow.h \
    $$PWD/player/FlacPic.h \
    $$PWD/player/ID3v2Pic.h \
    $$PWD/player/fftdisplay.h \
    $$PWD/player/lrcbar.h \
    $$PWD/player/lyrics.h \
    $$PWD/player/osd.h \
    $$PWD/player/player.h \
    $$PWD/player/shadowlabel.h \
    $$PWD/player/spslider.h \
    $$PWD/playlistmanagewindow.h \
    $$PWD/playlistmodel.h \
    $$PWD/qtlocalpeer.h \
    $$PWD/qtlockedfile.h \
    $$PWD/songlistmodel.h \
    $$PWD/sqlite3helper.h \
    $$PWD/comboboxdelegate.h \
    $$PWD/player/playlist.h \
    $$PWD/player/shadowplayer.h \
    $$PWD/parser.h \
    $$PWD/parsercsv.h \
    $$PWD/parserm3u.h \
    $$PWD/parserpls.h \
    $$PWD/trackfile.h \
    $$PWD/qmldialog.h \
    $$PWD/player/qmlplayer.h

FORMS +=  \
    $$PWD/qmldialog.ui \
    $$PWD/configurationwindow.ui \
    $$PWD/playlistmanagewindow.ui \
    $$PWD/player/lrcbar.ui \
    $$PWD/player/osd.ui \
    $$PWD/player/playlist.ui \
    $$PWD/player/shadowplayer.ui
    
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
#translate.depends = lrelease
#QMAKE_EXTRA_TARGETS += lupdate lrelease translate qti18n
#POST_TARGETDEPS += translate qti18n

win32: {
    lessThan(QT_MAJOR_VERSION, 6): QT += winextras
    CONFIG(release, debug|release) : {
        WINDEPLOYQT = $$shell_path($$[QT_INSTALL_BINS]/windeployqt.exe)
        QMAKE_EXTRA_TARGETS += mkdir
        
        #qti18n.depends = translate
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
    contains(QMAKE_HOST.arch, x86_64): {
        LIBS += -L$$PWD/3rdparty/bass/lib/windows/amd64
    } else: {
        LIBS += -L$$PWD/3rdparty/bass/lib/windows/386
    }
    LIBS += -lbasswasapi -lbassasio -lbassmix
}

LIBS += -L$$PWD/../lib/reverseProxy -lrp -lbass -lbass_fx -ltags

macx : {
    HEADERS += \
        $$PWD/serviceslots.h \
        $$PWD/application.h
    SOURCES += \
        $$PWD/application.cpp
    OBJECTIVE_HEADERS += \
        $$PWD/service.h
    OBJECTIVE_SOURCES += \
        $$PWD/service.mm
    DESTDIR = $$OUT_PWD
    QMAKE_INFO_PLIST = $$PWD/macInfo.plist
    ICON = $$PWD/hannah.icns
    icon.path = $$PWD
    INSTALLS += icon
    LIBS += -framework Security
    LIBS += -L$$PWD/3rdparty/bass/lib/mac
    
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
    hannah.qrc \
    rc/icons/icons.qrc \
    rc/imagine-assets/imagine-assets.qrc
