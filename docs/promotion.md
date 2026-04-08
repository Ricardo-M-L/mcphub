# MCP Hub Promotion Materials

## One-liner

MCP Hub - The npm for MCP servers. Search, install, and manage MCP servers with one command.

---

## Twitter/X Post

```
Launched MCP Hub - the package manager for MCP servers.

Two commands to get started:

$ mcphub search filesystem
$ mcphub install io.github.xxx/server-filesystem

Auto-configures Claude Code, Claude Desktop, and Cursor.

Built in Go. Single binary. Zero dependencies.

https://github.com/Ricardo-M-L/mcphub

#MCP #ClaudeCode #AI #OpenSource
```

---

## Reddit Post (r/ClaudeAI, r/LocalLLaMA, r/programming)

**Title:** I built MCP Hub - an npm-like package manager for MCP servers

**Body:**

I've been working with MCP (Model Context Protocol) servers and found the install/config process painful. You have to:

1. Find the server on GitHub
2. Read the docs to figure out the config
3. Manually edit JSON config files for Claude Desktop / Cursor
4. Repeat for every server

So I built **MCP Hub** - a CLI that does all of this in one command:

```bash
mcphub search database
mcphub install io.github.xxx/server-database
```

It queries the official MCP Registry, installs the server, and auto-configures your MCP clients (Claude Desktop, Cursor). Config files are backed up before modification.

**It also works as an MCP server itself** - add it to Claude Code and you can say "search for database MCP servers" in conversation:

```bash
go install github.com/Ricardo-M-L/mcphub/mcp@latest
claude mcp add mcphub mcphub-mcp
```

Built in Go, single binary, zero dependencies. Available via curl, Homebrew, npm, and go install.

GitHub: https://github.com/Ricardo-M-L/mcphub

Would love feedback!

---

## Hacker News (Show HN)

**Title:** Show HN: MCP Hub – A package manager for MCP servers (like npm for AI tools)

**URL:** https://github.com/Ricardo-M-L/mcphub

---

## V2EX Post

**Title:** MCP Hub - MCP 服务的包管理器，一行命令安装和管理 MCP 服务

**Body:**

做了一个 MCP 服务的包管理器，类似 npm 但专门管理 MCP 服务。

痛点：每次给 Claude Desktop 或 Cursor 装 MCP 服务，都得手动找仓库、看文档、编辑 JSON 配置文件，很麻烦。

MCP Hub 一行命令搞定：

```bash
mcphub search database
mcphub install io.github.xxx/server-database
```

自动检测你装了 Claude Desktop 还是 Cursor，直接写入配置。

还可以作为 MCP 服务接入 Claude Code，在对话里直接说"搜索数据库 MCP"就能搜索和安装。

Go 写的，单二进制，零依赖。支持 curl / Homebrew / npm / go install 安装。

GitHub: https://github.com/Ricardo-M-L/mcphub

欢迎反馈！

---

## Discord (Claude Community, AI Discord servers)

```
Hey! I built MCP Hub - a package manager for MCP servers.

Instead of manually editing JSON configs, just:
$ mcphub search filesystem
$ mcphub install io.github.xxx/server-filesystem

It auto-configures Claude Desktop and Cursor. Also works as an MCP server inside Claude Code.

GitHub: https://github.com/Ricardo-M-L/mcphub

Feedback welcome!
```

---

## Chinese Tech Communities (掘金/CSDN/知乎)

**Title:** 开源：MCP Hub - MCP 服务的 npm，一行命令管理所有 MCP 服务

**Body:**

## 为什么做这个

MCP（Model Context Protocol）是 2026 年 AI 领域最热的协议，但安装和管理 MCP 服务的体验很差：

- 要去 GitHub 找服务
- 要读文档看怎么配
- 要手动编辑 Claude Desktop / Cursor 的 JSON 配置
- 每个服务都重复这个过程

## MCP Hub 是什么

一个 MCP 服务的包管理器，类似 npm：

```bash
# 搜索
mcphub search database

# 安装（自动配置 Claude Desktop 和 Cursor）
mcphub install io.github.xxx/server-database

# 查看已安装
mcphub list

# 卸载
mcphub remove xxx
```

## 亮点

- Go 编写，单二进制，零运行时依赖
- 自动检测并配置 Claude Desktop、Cursor
- 可作为 MCP 服务接入 Claude Code，直接在对话中搜索和安装
- 支持 curl / Homebrew / npm / go install 安装
- 完整的 Registry API 服务器 + Web UI
- 14 个单元测试，CI/CD 自动发布

## 安装

```bash
# 方式 1：curl
curl -fsSL https://raw.githubusercontent.com/Ricardo-M-L/mcphub/master/install.sh | sh

# 方式 2：Go
go install github.com/Ricardo-M-L/mcphub/cmd/mcphub@latest

# 接入 Claude Code
go install github.com/Ricardo-M-L/mcphub/mcp@latest
claude mcp add mcphub mcphub-mcp
```

GitHub: https://github.com/Ricardo-M-L/mcphub

Star 支持一下！欢迎贡献代码和反馈。
