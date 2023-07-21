cmake_minimum_required(VERSION 3.15)
project(desktop VERSION 1.0.0)

set(ProjectName Hannah)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_subdirectory(3rdparty)

find_package(QT NAMES Qt6 Qt5 REQUIRED COMPONENTS Widgets Network Quick QuickWidgets QuickControls2 LinguistTools)
find_package(Qt${QT_VERSION_MAJOR} REQUIRED COMPONENTS Widgets Network Quick QuickWidgets QuickControls2 LinguistTools)
if(${QT_VERSION_MAJOR} GREATER_EQUAL 6)
    find_package(Qt6 COMPONENTS Core5Compat REQUIRED)
endif()

# Instruct CMake to run moc automatically when needed.
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTOUIC ON) # if you have any .ui files
set(CMAKE_AUTORCC ON) # if you have any .qrc files

include_directories(${PROJECT_SOURCE_DIR} ${CMAKE_CURRENT_SOURCE_DIR}/player)
include_directories(${PROJECT_BINARY_DIR})

if(MSVC)
    add_compile_definitions(_SILENCE_CXX17_C_HEADER_DEPRECATION_WARNING)
endif()

set(SOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/main.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/configurationwindow.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/fftdisplay.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/lrcbar.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/lyrics.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/osd.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/bassplayer.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/shadowlabel.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/spslider.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/playlistmanagewindow.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/playlistmodel.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/qtlocalpeer.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/songlistmodel.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/sqlite3helper.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/comboboxdelegate.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/playlist.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/shadowplayer.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/parser.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/parsercsv.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/parserm3u.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/parserpls.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/trackfile.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/qmldialog.cpp 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/qmlplayer.cpp
)

    
set(HEADERS
    ${CMAKE_CURRENT_SOURCE_DIR}/configurationwindow.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/FlacPic.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/ID3v2Pic.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/fftdisplay.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/lrcbar.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/lyrics.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/osd.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/bassplayer.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/shadowlabel.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/spslider.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/playlistmanagewindow.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/playlistmodel.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/qtlocalpeer.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/qtlockedfile.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/songlistmodel.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/sqlite3helper.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/comboboxdelegate.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/playlist.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/shadowplayer.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/parser.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/parsercsv.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/parserm3u.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/parserpls.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/trackfile.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/qmldialog.h 
    ${CMAKE_CURRENT_SOURCE_DIR}/player/qmlplayer.h
    )

file(GLOB_RECURSE
    UIS
    ${CMAKE_CURRENT_SOURCE_DIR}/*.ui)
    
file(GLOB_RECURSE
    RESOURCES
    ${CMAKE_CURRENT_SOURCE_DIR}/*.qrc)

qt_add_resources(qrcs ${RESOURCES})

IF(WIN32)
    list(APPEND SOURCES 
        ${CMAKE_CURRENT_SOURCE_DIR}/qtsinglecoreapplication.cpp 
        ${CMAKE_CURRENT_SOURCE_DIR}/qtsingleapplication.cpp
        )
    list(APPEND HEADERS 
        ${CMAKE_CURRENT_SOURCE_DIR}/qtsinglecoreapplication.h 
        ${CMAKE_CURRENT_SOURCE_DIR}/qtsingleapplication.h
        )
ENDIF()

IF(APPLE)
    list(APPEND SOURCES 
        ${CMAKE_CURRENT_SOURCE_DIR}/application.cpp 
        ${CMAKE_CURRENT_SOURCE_DIR}/service.mm
        )
    list(APPEND HEADERS 
        ${CMAKE_CURRENT_SOURCE_DIR}/serviceslots.h 
        ${CMAKE_CURRENT_SOURCE_DIR}/application.h 
        ${CMAKE_CURRENT_SOURCE_DIR}/service.h
        )
ENDIF()

SET(
    TS_FILES
    ${CMAKE_CURRENT_SOURCE_DIR}/translations/Hannah_zh_CN.ts
)

if(${QT_VERSION_MAJOR} GREATER_EQUAL 6)
    qt_create_translation(QM_FILES ${CMAKE_SOURCE_DIR} ${TS_FILES}
    OPTIONS -no-obsolete)
    qt_add_executable(${ProjectName}
            ${CMAKE_CURRENT_SOURCE_DIR}/${ProjectName}.rc
            ${qrcs}
            ${UIS}
            ${HEADERS}
            ${SOURCES}
            ${QM_FILES}
            )
else()
    qt5_create_translation(QM_FILES ${CMAKE_SOURCE_DIR}  ${TS_FILES}
    OPTIONS -no-obsolete)
    add_executable(${ProjectName}
            ${CMAKE_CURRENT_SOURCE_DIR}/${ProjectName}.rc
            ${qrcs}
            ${UIS}
            ${HEADERS}
            ${SOURCES}
            ${QM_FILES}
            )
endif()

IF(Boost_FOUND)
    include_directories(${Boost_INCLUDE_DIRS})
    IF(MSVC)
        add_compile_definitions(
            BOOST_AUTO_LINK_NOMANGLE
            )
        target_link_directories(${ProjectName} PRIVATE "${Boost_INCLUDE_DIRS}/lib/${BUILD_TYPE}")
    ENDIF(MSVC)
ENDIF(Boost_FOUND)

target_link_libraries(${ProjectName} PRIVATE
    sqlite3
    basslib
    reverseProxy
    Qt::Widgets
    Qt::Quick
    Qt::Network
    Qt::QuickWidgets
    Qt::QuickControls2
    )

IF(APPLE) 
    target_link_libraries(${ProjectName} PRIVATE "-framework Security")
ENDIF()

set_target_properties(${ProjectName} PROPERTIES
    WIN32_EXECUTABLE ON
    MACOSX_BUNDLE ON
    )


if(QT_VERSION_MAJOR EQUAL 6)
    target_link_libraries(${ProjectName} PRIVATE
        Qt::Core5Compat
        )
    qt_finalize_project()
endif()
    