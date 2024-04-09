#ifndef JSON_NODE_H
#define JSON_NODE_H

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

    JsonNode(std::string key = "", std::string parent_id = "", int level = 0, int index = 0);

    void add_value(std::string key, nlohmann::json value);

    void add_list(std::string key, nlohmann::json value);

    void add_child(JsonNode child);
};

JsonNode parse_json(nlohmann::json data, std::string key = "", std::string parent_id = "", int level = 0, int index = 0);

void process_list(nlohmann::json lst, JsonNode &parent_node, std::string key, int level);

#endif // JSON_NODE_H
