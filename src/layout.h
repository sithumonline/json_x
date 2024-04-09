#ifndef LAYOUT_H
#define LAYOUT_H

#include <vector>
#include <unordered_map>
#include <string>

namespace Layout
{
    struct NodeX
    {
        int id;
        int level;
        int position;
    };

    std::unordered_map<int, std::vector<int>> buildTree(const std::vector<std::pair<int, int>> &edges);

    std::vector<int> findRoots(const std::unordered_map<int, std::vector<int>> &children);

    std::unordered_map<int, std::pair<int, int>> layoutTree(const std::vector<std::pair<int, int>> &edges, const std::string &direction, int nodeSpacing, int levelSpacing);
}

#endif // LAYOUT_H