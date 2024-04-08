#include <iostream>
#include <queue>
#include <algorithm>
#include <unordered_map>

namespace Layout
{
    struct NodeX
    {
        int id;
        int level;
        int position;
    };

    std::unordered_map<int, std::vector<int>> buildTree(const std::vector<std::pair<int, int>> &edges)
    {
        std::unordered_map<int, std::vector<int>> children;
        for (const auto &edge : edges)
        {
            children[edge.first].push_back(edge.second);
        }
        return children;
    }

    std::vector<int> findRoots(const std::unordered_map<int, std::vector<int>> &children)
    {
        std::vector<int> roots;
        for (const auto &node : children)
        {
            if (std::find_if(children.begin(), children.end(), [&](const std::pair<int, std::vector<int>> &p)
                             { return std::find(p.second.begin(), p.second.end(), node.first) != p.second.end(); }) == children.end())
            {
                roots.push_back(node.first);
            }
        }
        return roots;
    }

    std::unordered_map<int, std::pair<int, int>> layoutTree(const std::vector<std::pair<int, int>> &edges, const std::string &direction, int nodeSpacing, int levelSpacing)
    {
        std::unordered_map<int, std::vector<int>> children = buildTree(edges);
        std::vector<int> roots = findRoots(children);
        std::unordered_map<int, std::pair<int, int>> positions;
        std::queue<NodeX> q;
        for (const auto &root : roots)
        {
            q.push(NodeX{root, 0, 0});
        }
        std::unordered_map<int, int> levelWidth;
        while (!q.empty())
        {
            NodeX node = q.front();
            q.pop();
            int x, y;
            if (direction == "top-to-bottom")
            {
                x = node.position * nodeSpacing;
                y = node.level * levelSpacing;
            }
            else
            { // "left-to-right"
                x = node.level * levelSpacing;
                y = node.position * nodeSpacing;
            }
            positions[node.id] = {x, y};
            for (const auto &child : children[node.id])
            {
                levelWidth[node.level + 1]++;
                q.push(NodeX{child, node.level + 1, levelWidth[node.level + 1]});
            }
        }
        return positions;
    }
}
