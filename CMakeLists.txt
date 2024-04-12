cmake_minimum_required(VERSION 3.10)

project(json_x_glfw_opengl3)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(OpenGL REQUIRED)

include_directories(${OPENGL_INCLUDE_DIRS})

set(GLFW_BUILD_DOCS OFF CACHE BOOL "GLFW lib only")
set(GLFW_INSTALL OFF CACHE BOOL "GLFW lib only")

add_subdirectory(glfw)

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
    set(PLATFORM "-linux-amd64")
    set(LIBS ${OPENGL_LIBRARIES})
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Darwin")
    set(PLATFORM "-macos-arm64")
    set(LIBS ${OPENGL_LIBRARIES})
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Windows")
    set(PLATFORM "-windows-amd64")
    set(LIBS gdi32 opengl32 imm32)
endif()

set(SHA_SHORT "$ENV{SHA_SHORT}")
if(SHA_SHORT)
    set(EXE "json_x_glfw_opengl3${PLATFORM}_${SHA_SHORT}")
else()
    set(EXE "json_x_glfw_opengl3${PLATFORM}")
endif()

add_executable(${EXE} ${SOURCES})
target_link_libraries(${EXE} glfw ${LIBS})

if(MINGW)
    set_target_properties(${EXE} PROPERTIES LINK_FLAGS "-static-libgcc -static-libstdc++")
endif()
