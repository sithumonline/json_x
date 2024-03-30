#include "json_x.h"
#include <imnodes.h>
#include <imgui.h>
#include "json_node.h"
#include <json.hpp>
using json = nlohmann::json;

namespace jsonX
{
    namespace
    {
        class HelloWorldNodeEditor
        {
        public:
            void show()
            {
                ImGui::Begin("simple node editor");

                ImNodes::BeginNodeEditor();

                json ex1 = {
                    {"pi", 3.141},
                    {"happy", true},
                    {"name", "Niels"},
                    {"nothing", nullptr},
                    {"answer", {{"everything", 42}}},
                    {"list", {1, 0, 2}},
                    {"object", {{"currency", "USD"}, {"value", 42.99}}}};

                JsonNode jNodes = parse_json(ex1);
                createNodes(jNodes);

                ImNodes::EndNodeEditor();

                ImGui::End();
            }

            void createNodes(JsonNode node, int parent = 0)
            {
                uint64_t nodeId = std::stoull(node.id);
                if (!node.key.empty())
                {
                    ImNodes::BeginNode(nodeId);
                    ImNodes::BeginInputAttribute(nodeId + 1);
                    ImGui::Text(node.key.c_str());
                    ImNodes::EndInputAttribute();
                    ImNodes::BeginOutputAttribute(nodeId + 2);
                    ImNodes::EndOutputAttribute();
                    ImNodes::EndNode();
                }
                else if (!node.values.empty())
                {
                    ImNodes::BeginNode(nodeId);
                    ImNodes::BeginInputAttribute(nodeId + 1);

                    std::stringstream txt;
                    for (const auto &[k, v] : node.values)
                    {
                        txt << k << ": " << v << "\n";
                    }

                    ImGui::Text(txt.str().c_str());
                    ImNodes::EndInputAttribute();
                    ImNodes::BeginOutputAttribute(nodeId + 2);
                    ImNodes::EndOutputAttribute();
                    ImNodes::EndNode();
                }
                else
                {
                    ImNodes::BeginNode(nodeId);
                    ImNodes::BeginInputAttribute(nodeId + 1);
                    ImNodes::EndInputAttribute();
                    ImNodes::BeginOutputAttribute(nodeId + 2);
                    ImNodes::EndOutputAttribute();
                    ImNodes::EndNode();
                }

                if (parent != 0)
                {
                    ImNodes::Link(nodeId + 3, parent, nodeId + 1);
                }

                for (JsonNode &child : node.children)
                {
                    createNodes(child, nodeId + 2);
                }
            }
        };

        static HelloWorldNodeEditor editor;
    } // namespace

    void NodeEditorInitialize() { ImNodes::SetNodeGridSpacePos(1, ImVec2(200.0f, 200.0f)); }

    void NodeEditorShow() { editor.show(); }

    void NodeEditorShutdown() {}

} // namespace example