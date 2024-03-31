#include "json_x.h"
#include <imnodes.h>
#include <imgui.h>
#include "json_node.h"
#include <json.hpp>
#include <iostream>
#include <time.h>

using json = nlohmann::json;

namespace jsonX
{
    namespace
    {
        class IDGenerator
        {
        public:
            static int GenerateNodeID()
            {
                static int nodeId = 0;
                return ++nodeId;
            }

            static int GenerateLinkID()
            {
                static int linkId = 10000; // Start link IDs from 10000 to differentiate from node IDs
                return ++linkId;
            }
        };

        struct Node
        {
            int id;
            std::string text;

            Node(int i, std::string text) : id(i), text(text)
            {
            }
            Node(std::string text)
                : text(text)
            {
                id = IDGenerator::GenerateNodeID();
            }

            int getId()
            {
                return id;
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
            Link(int sa, int ea) : start_attr(sa), end_attr(ea)
            {
                id = IDGenerator::GenerateLinkID();
            }

            int getId()
            {
                return id;
            }
        };

        struct Editor
        {
            ImNodesEditorContext *context = nullptr;
            std::vector<Node> nodes;
            std::vector<Link> links;
            int current_id = 0;
            std::size_t size;
        };

        void createNode(Node node)
        {
            ImNodes::BeginNode(node.id);

            ImNodes::BeginNodeTitleBar();
            ImGui::TextUnformatted("node");
            ImNodes::EndNodeTitleBar();

            ImNodes::BeginInputAttribute(node.id << 8);
            ImGui::TextUnformatted(node.text.c_str()); // Convert std::string to const char*
            ImNodes::EndInputAttribute();

            ImNodes::BeginOutputAttribute(node.id << 16);
            ImNodes::EndOutputAttribute();

            ImNodes::EndNode();
        }

        void createNodes(JsonNode node, Editor &editor, int parent = 0)
        {
            std::cout << "Create Node " << node.id << std::endl;
            int nodeId;
            if (!node.key.empty())
            {
                Node no = Node(node.key);
                nodeId = no.getId();
                editor.nodes.push_back(no);
            }
            else if (!node.values.empty())
            {
                std::stringstream txt;
                for (const auto &[k, v] : node.values)
                {
                    txt << k << ": " << v << "\n";
                }

                Node no = Node(txt.str());
                nodeId = no.getId();
                editor.nodes.push_back(no);
            }
            else
            {
                Node no = Node("root");
                nodeId = no.getId();
                editor.nodes.push_back(no);
            }

            if (parent != 0)
            {
                Link lin = Link(parent, nodeId << 8);
                editor.links.push_back(lin);
            }

            for (const auto &[k, v] : node.lists)
            {
                Node no = Node(k);
                editor.nodes.push_back(no);

                Link lin = Link(nodeId << 16, no.getId() << 8);
                editor.links.push_back(lin);

                for (long unsigned int j = 0; j < v.size(); ++j)
                {
                    std::cout << j << " showEditor" << std::endl;
                    for (auto &[key, val] : v[j].items())
                    {
                        Node no1 = Node(val);
                        editor.nodes.push_back(no1);

                        Link lin = Link(no.getId() << 16, no1.getId() << 8);
                        editor.links.push_back(lin);
                    }
                }
            }

            if (!node.key.empty() && !node.values.empty())
            {
                std::stringstream txt;
                for (const auto &[k, v] : node.values)
                {
                    txt << k << ": " << v << "\n";
                }

                Node no = Node(txt.str());
                editor.nodes.push_back(no);

                Link lin = Link(nodeId << 16, no.getId() << 8);
                editor.links.push_back(lin);
            }

            for (JsonNode &child : node.children)
            {
                createNodes(child, editor, nodeId << 16);
            }
        }

        void show_editor(const char *editor_name, Editor &editor)
        {
            std::cout << time(0) << " showEditor" << std::endl;

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
                {"list", {"x", "y", "z"}},
                {"object", {{"currency", "USD"}, {"value", 42.99}}}};

            JsonNode jNodes = parse_json(ex1);

            std::size_t size = ex1.size();
            if (editor.size != size)
            {
                editor.size = size;
                createNodes(jNodes, editor);
            }

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

        static Editor editor1;
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
