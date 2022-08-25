cmake_minimum_required(VERSION 3.15)
project(reverseProxy VERSION 1.0.0)

add_library(reverseProxy INTERFACE)

target_include_directories(reverseProxy INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})

target_link_directories(reverseProxy INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(reverseProxy INTERFACE rp)