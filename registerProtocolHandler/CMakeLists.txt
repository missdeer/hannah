cmake_minimum_required(VERSION 3.15)
project(registerProtocolHandler VERSION 1.0.0)

set(ProjectName registerProtocolHandler)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(QT NAMES Qt6 Qt5 REQUIRED COMPONENTS Core)
find_package(Qt${QT_VERSION_MAJOR} REQUIRED COMPONENTS Core)

# Instruct CMake to run moc automatically when needed.
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTOUIC ON) # if you have any .ui files
set(CMAKE_AUTORCC ON) # if you have any .qrc files

include_directories(${PROJECT_SOURCE_DIR})
include_directories(${PROJECT_BINARY_DIR})

if (MSVC)
    add_compile_options(/utf-8)

    SET(CMAKE_EXE_LINKER_FLAGS  "${CMAKE_EXE_LINKER_FLAGS} /MANIFESTUAC:\"level=\'requireAdministrator\' uiAccess=\'false\'\"")
endif()


if(${QT_VERSION_MAJOR} GREATER_EQUAL 6)
    qt_add_executable(${ProjectName}
            ${CMAKE_CURRENT_SOURCE_DIR}/${ProjectName}.rc
            main.cpp
            )
else()
    add_executable(${ProjectName}
            ${CMAKE_CURRENT_SOURCE_DIR}/${ProjectName}.rc
            main.cpp
            )
endif()

if(QT_VERSION_MAJOR EQUAL 6)
    qt_finalize_project()
endif()