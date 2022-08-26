cmake_minimum_required(VERSION 3.15)
project(Hannah VERSION 1.0.0)

set(CMAKE_EXPORT_COMPILECOMMANDS ON)

IF(CMAKE_BUILD_TYPE MATCHES Debug)
    set(BUILD_TYPE Debug)
    add_compile_definitions(LOGS_ENABLED)
ELSE(CMAKE_BUILD_TYPE MATCHES Debug)
    set(BUILD_TYPE Release)
ENDIF(CMAKE_BUILD_TYPE MATCHES Debug)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}/bin)

add_subdirectory(lib)
add_subdirectory(desktop)
add_subdirectory(registerProtocolHandler)

if(QT_VERSION_MAJOR EQUAL 6)
qt_finalize_project()
endif()
