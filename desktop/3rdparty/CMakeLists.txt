cmake_minimum_required(VERSION 3.15)
project(3rdparty VERSION 1.0.0)

add_subdirectory(bass)

IF(APPLE)
    add_subdirectory (sqlite3)
    add_subdirectory (zlib)
ELSE()
    find_package(sqlite3)
    if(NOT SQLite3_FOUND)
        add_subdirectory(sqlite3)
    endif()

    find_package(ZLIB)
    IF(NOT ZLIB_FOUND)
        add_subdirectory (zlib)
    ENDIF()
ENDIF()
