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
        struct Node
        {
            int id;
            std::string text;

            Node(int i, std::string text)
                : id(i), text(text)
            {
            }
        };

        struct Link
        {
            int id;
            int start_attr, end_attr;

            Link() {}
            Link(int id, int sa, int ea) : id(id), start_attr(sa), end_attr(ea)
            {
            }
        };

        struct Editor
        {
            ImNodesEditorContext *context = nullptr;
            std::vector<Node> nodes;
            std::vector<Link> links;
            int current_id = 0;
        };

        void createNode(Node node)
        {
            ImNodes::BeginNode(node.id);

            ImNodes::BeginNodeTitleBar();
            ImGui::TextUnformatted("node");
            ImNodes::EndNodeTitleBar();

            ImNodes::BeginInputAttribute(node.id + 1);
            ImGui::TextUnformatted(node.text.c_str()); // Convert std::string to const char*
            ImNodes::EndInputAttribute();

            ImNodes::BeginOutputAttribute(node.id + 2);
            ImNodes::EndOutputAttribute();

            ImNodes::EndNode();
        }

        void createNodes(JsonNode node, Editor &editor, int parent = 0)
        {
            uint64_t nodeId = std::stoull(node.id);
            if (!node.key.empty())
            {
                editor.nodes.push_back(Node(nodeId, node.key));
            }
            else if (!node.values.empty())
            {
                std::stringstream txt;
                for (const auto &[k, v] : node.values)
                {
                    txt << k << ": " << v << "\n";
                }

                editor.nodes.push_back(Node(nodeId, txt.str()));
            }
            else
            {
                editor.nodes.push_back(Node(nodeId, ""));
            }

            if (parent != 0)
            {
                editor.links.push_back(Link(nodeId + 3, parent, nodeId + 1));
            }

            for (JsonNode &child : node.children)
            {
                createNodes(child, editor, nodeId + 2);
            }
        }

        void show_editor(const char *editor_name, Editor &editor)
        {
            ImNodes::EditorContextSet(editor.context);

            ImGui::Begin(editor_name);
            ImGui::TextUnformatted("Json X");

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
            createNodes(jNodes, editor);

            for (Node &node : editor.nodes)
            {
                createNode(node);
            }

            for (const Link &link : editor.links)
            {
                ImNodes::Link(link.id, link.start_attr, link.end_attr);
            }

            ImNodes::EndNodeEditor();

            {
                Link link;
                if (ImNodes::IsLinkCreated(&link.start_attr, &link.end_attr))
                {
                    link.id = ++editor.current_id;
                    editor.links.push_back(link);
                }
            }

            {
                int link_id;
                if (ImNodes::IsLinkDestroyed(&link_id))
                {
                    auto iter = std::find_if(
                        editor.links.begin(), editor.links.end(), [link_id](const Link &link) -> bool
                        { return link.id == link_id; });
                    assert(iter != editor.links.end());
                    editor.links.erase(iter);
                }
            }

            ImGui::End();
        }

        Editor editor1;
    } // namespace

    void NodeEditorInitialize()
    {
        editor1.context = ImNodes::EditorContextCreate();
        ImNodes::PushAttributeFlag(ImNodesAttributeFlags_EnableLinkDetachWithDragClick);

        ImNodesIO &io = ImNodes::GetIO();
        io.LinkDetachWithModifierClick.Modifier = &ImGui::GetIO().KeyCtrl;
        io.MultipleSelectModifier.Modifier = &ImGui::GetIO().KeyCtrl;

        ImNodesStyle &style = ImNodes::GetStyle();
        style.Flags |= ImNodesStyleFlags_GridLinesPrimary | ImNodesStyleFlags_GridSnapping;
    }

    void NodeEditorShow()
    {
        show_editor("Json Editor", editor1);
    }

    void NodeEditorShutdown()
    {
        ImNodes::PopAttributeFlag();
        ImNodes::EditorContextFree(editor1.context);
    }
} // namespace example
