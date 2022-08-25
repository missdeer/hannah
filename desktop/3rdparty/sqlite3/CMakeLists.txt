cmake_minimum_required(VERSION 3.15)
project(sqlite3 VERSION 1.0.0)

include_directories(${CMAKE_CURRENT_SOURCE_DIR})

add_library(sqlite3 STATIC sqlite3.c sqlite3.h sqlite3ext.h)

target_include_directories(sqlite3 PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})