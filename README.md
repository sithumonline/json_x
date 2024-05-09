# Json X

[![image.png](https://i.postimg.cc/PrnQV06W/image.png)](https://postimg.cc/5jpz6KWy)

Json X is a cross-platform desktop application for JSON manipulation written in pure C++.

## Features

- **JSON Viewer**: View _large_ JSON data in a tree structure.
- **Fast**: Loads JSON data in milliseconds.
- **Simple and Clean UI**: Easy to use and navigate.
- **Single Executable**: No dependencies, no installation required.
- **Cross-Platform**: Works on Windows, Linux, and macOS.
- **Lightweight**: Less than 5MB in size.

## Installation

Download the latest release from the [releases page](https://github.com/sithumonline/json_x/releases) or 
[json_x](https://sithum.online/json_x/#download) and run the executable.

## Usage

1. Open the application.
2. Click on the `Open File Dialog` button.
3. Select a JSON file.
4. View the JSON data in the tree structure.

## Build from Source

1. Clone the repository.

```bash
git clone git@github.com:sithumonline/json_x.git
```

2. Build the project.

```bash
cd json_x
mkdir build
cd build
cmake ..
make
```

3. Run the executable.

```bash
./json_x_glfw_opengl3
```
