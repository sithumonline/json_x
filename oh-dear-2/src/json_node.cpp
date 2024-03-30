#include <vector>
#include <string>
#include <map>
#include <json.hpp>

class JsonNode
{
public:
    std::string id;
    std::string key;
    std::string parent_id;
    int level;
    int index;
    std::map<std::string, nlohmann::json> values;
    std::map<std::string, std::vector<nlohmann::json>> lists;
    std::vector<JsonNode> children;

    JsonNode(std::string key = "", std::string parent_id = "", int level = 0, int index = 0)
        : key(key), parent_id(parent_id), level(level), index(index)
    {
        // Generate a random UUID for the id
        // This requires a separate library or system call, so it's omitted here
    }

    void add_value(std::string key, nlohmann::json value)
    {
        values[key] = value;
    }

    void add_list(std::string key, nlohmann::json value)
    {
        lists[key].push_back(value);
    }

    void add_child(JsonNode child)
    {
        children.push_back(child);
    }
};

void process_list(nlohmann::json lst, JsonNode &parent_node, std::string key, int level);

JsonNode parse_json(nlohmann::json data, std::string key = "", std::string parent_id = "", int level = 0, int index = 0)
{
    JsonNode node(key, parent_id, level, index);

    if (data.is_object())
    {
        for (auto &[k, v] : data.items())
        {
            if (v.is_object())
            {
                JsonNode child = parse_json(v, k, node.id, level + 1, 0);
                node.add_child(child);
            }
            else if (v.is_array())
            {
                process_list(v, node, k, level + 1);
            }
            else
            {
                node.add_value(k, v);
            }
        }
    }

    return node;
}

void process_list(nlohmann::json lst, JsonNode &parent_node, std::string key, int level)
{
    int i = 0;
    for (auto &item : lst)
    {
        if (item.is_object())
        {
            JsonNode child = parse_json(item, key, parent_node.id, level, i++);
            parent_node.add_child(child);
        }
        else
        {
            parent_node.add_list(key, item);
        }
    }
}