cmake_minimum_required(VERSION 3.15)
project(Hannah VERSION 1.0.0)

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
set(Boost_NO_WARN_NEW_VERSIONS ON)

find_package(Boost 1.78.0 REQUIRED system)

if(Boost_FOUND)
     include_directories(${Boost_INCLUDE_DIRS})
     link_directories(${Boost_LIBRARY_DIRS})
endif(Boost_FOUND)

if(NOT APPLE)
  find_package(ZLIB)
  find_package(sqlite3)
endif()

if(MSVC)
    # First check the command line variable VCPKG_TARGET_TRIPLET
    if(DEFINED VCPKG_TARGET_TRIPLET)
        set(TRIPLET ${VCPKG_TARGET_TRIPLET})
    else()
        # If it's not defined, check the environment variable
        set(TRIPLET $ENV{VCPKG_TARGET_TRIPLET})
    endif()

    # Now check if TRIPLET ends with "-windows-static"
    string(FIND "${TRIPLET}" "-windows-static" POSITION)
    if(${POSITION} GREATER_EQUAL 0)
        if(NOT MSVC)
            set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -static")
        else()
            add_compile_options("/MT$<$<CONFIG:Debug>:d>")
            if(POLICY CMP0091)
                cmake_policy(SET CMP0091 NEW)
            endif()
        endif()
    endif()
    add_compile_definitions(_SILENCE_CXX17_C_HEADER_DEPRECATION_WARNING)
endif()

set(SOURCES
    main.cpp 
    configurationwindow.cpp 
    player/fftdisplay.cpp 
    player/lrcbar.cpp 
    player/lyrics.cpp 
    player/osd.cpp 
    player/bassplayer.cpp 
    player/shadowlabel.cpp 
    player/spslider.cpp 
    player/qmlplayer.cpp
    player/playlist.cpp 
    player/shadowplayer.cpp 
    playlistmanagewindow.cpp 
    playlistmodel.cpp 
    qtlocalpeer.cpp 
    songlistmodel.cpp 
    sqlite/Sqlite3Helper.cpp 
    sqlite/Sqlite3DBManager.cpp
    sqlite/Sqlite3Statement.cpp
    comboboxdelegate.cpp 
    parser.cpp 
    parsercsv.cpp 
    parserm3u.cpp 
    parserpls.cpp 
    trackfile.cpp 
    qmldialog.cpp 
    network/networkreplyhelper.cpp
)

    
set(HEADERS
    configurationwindow.h 
    player/FlacPic.h 
    player/ID3v2Pic.h 
    player/fftdisplay.h 
    player/lrcbar.h 
    player/lyrics.h 
    player/osd.h 
    player/bassplayer.h 
    player/shadowlabel.h 
    player/spslider.h 
    player/playlist.h 
    player/shadowplayer.h 
    player/qmlplayer.h
    playlistmanagewindow.h 
    playlistmodel.h 
    qtlocalpeer.h 
    qtlockedfile.h 
    songlistmodel.h 
    sqlite/Sqlite3Helper.h 
    sqlite/Sqlite3DBManager.h
    sqlite/Sqlite3Statement.h
    sqlite/Sqlite3Constants.h
    comboboxdelegate.h 
    parser.h 
    parsercsv.h 
    parserm3u.h 
    parserpls.h 
    trackfile.h 
    qmldialog.h 
    network/networkreplyhelper.h
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
        qtsinglecoreapplication.cpp 
        qtsingleapplication.cpp
        )
    list(APPEND HEADERS 
        qtsinglecoreapplication.h 
        qtsingleapplication.h
        )
ENDIF()

IF(APPLE)
    list(APPEND SOURCES 
        application.cpp 
        service.mm
        )
    list(APPEND HEADERS 
        serviceslots.h 
        application.h 
        service.h
        )
ENDIF()

SET(
    TS_FILES
    translations/Hannah_zh_CN.ts
)

if(${QT_VERSION_MAJOR} GREATER_EQUAL 6)
    qt_create_translation(QM_FILES ${CMAKE_SOURCE_DIR} ${TS_FILES}
    OPTIONS -no-obsolete)
    qt_add_executable(${PROJECT_NAME}
            ${PROJECT_NAME}.rc
            ${qrcs}
            ${UIS}
            ${HEADERS}
            ${SOURCES}
            ${QM_FILES}
            )
else()
    qt5_create_translation(QM_FILES ${CMAKE_SOURCE_DIR}  ${TS_FILES}
    OPTIONS -no-obsolete)
    add_executable(${PROJECT_NAME}
            ${PROJECT_NAME}.rc
            ${qrcs}
            ${UIS}
            ${HEADERS}
            ${SOURCES}
            ${QM_FILES}
            )
endif()

target_link_libraries(${PROJECT_NAME} PRIVATE
    basslib
    reverseProxy
    Boost::system
    Qt::Widgets
    Qt::Quick
    Qt::Network
    Qt::QuickWidgets
    Qt::QuickControls2
    )

if(SQLite3_FOUND)
    target_link_libraries(${PROJECT_NAME} PRIVATE SQLite::SQLite3)
else()
    target_link_libraries(${PROJECT_NAME} PRIVATE sqlite3)
endif()

if(ZLIB_FOUND)
  target_link_libraries(${PROJECT_NAME} PRIVATE ZLIB::ZLIB)
else()
  target_link_libraries(${PROJECT_NAME} PRIVATE zlibstatic)
endif()

IF(APPLE) 
    target_link_libraries(${PROJECT_NAME} PRIVATE "-framework Security")
ENDIF()

if(WIN32)
    target_link_libraries(${PROJECT_NAME} PRIVATE ws2_32)
endif()

if(QT_VERSION_MAJOR EQUAL 6)
    target_link_libraries(${PROJECT_NAME} PRIVATE
        Qt::Core5Compat
        )
endif()
    
set_target_properties(${PROJECT_NAME} PROPERTIES
    WIN32_EXECUTABLE ON
    MACOSX_BUNDLE ON
)
