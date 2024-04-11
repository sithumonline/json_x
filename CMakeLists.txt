cmake_minimum_required(VERSION 3.10)

project(json_x_glfw_opengl3)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

if(${CMAKE_SYSTEM_NAME} STREQUAL "Linux")
    set(PLATFORM "-linux-amd64")
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Darwin")
    set(PLATFORM "-macos-arm64")
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Windows")
    set(PLATFORM "-windows-amd64")
endif()

set(EXE "json_x_glfw_opengl3${PLATFORM}")

if(DEFINED SHA_SHORT)
  set(EXE "${EXE}_${SHA_SHORT}")
endif()

include_directories(
    ${CMAKE_CURRENT_SOURCE_DIR}/imgui
    ${CMAKE_CURRENT_SOURCE_DIR}/imgui/backends
    ${CMAKE_CURRENT_SOURCE_DIR}/imnodes
    ${CMAKE_CURRENT_SOURCE_DIR}/json
    ${CMAKE_CURRENT_SOURCE_DIR}/json/nlohmann
    ${CMAKE_CURRENT_SOURCE_DIR}/json/nlohmann/thirdparty
    ${CMAKE_CURRENT_SOURCE_DIR}/ImGuiFileDialog
)

file(GLOB_RECURSE SOURCES 
    src/*.cpp
    imgui/*.cpp
    imgui/backends/*.cpp
    ImGuiFileDialog/*.cpp
    imnodes/*.cpp
)
    
if(${CMAKE_SYSTEM_NAME} STREQUAL "Linux")
    find_package(PkgConfig REQUIRED)
    pkg_search_module(GLFW REQUIRED glfw3)
    include_directories(${GLFW_INCLUDE_DIRS})
    link_directories(${GLFW_LIBRARY_DIRS})
    add_definitions(${GLFW_CFLAGS_OTHER})
    set(LIBS ${GLFW_LIBRARIES} -lGL)
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Darwin")
    find_library(GLFW_LIBRARY NAMES glfw)
    find_library(OPENGL_LIBRARY OpenGL)
    find_library(COCOA_LIBRARY Cocoa)
    find_library(IOKIT_LIBRARY IOKit)
    find_library(COREVIDEO_LIBRARY CoreVideo)
    set(LIBS ${GLFW_LIBRARY} ${OPENGL_LIBRARY} ${COCOA_LIBRARY} ${IOKIT_LIBRARY} ${COREVIDEO_LIBRARY})
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Windows")
    find_package(glfw3 REQUIRED)
    set(LIBS glfw3 gdi32 opengl32 imm32)
endif()

add_executable(${EXE} ${SOURCES})
target_link_libraries(${EXE} ${LIBS})
