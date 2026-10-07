# sdd-tdd-loop-demo

老刀Cheney 视频《AI 写代码总跑偏？用 OpenSpec 和 Superpowers 搭一个最小的开发循环》的演示仓。

一个读温度的小程序。视频里从 tag `start` 出发，用 OpenSpec 写规格，用 Superpowers 先写测试再实现，给它加了一个超温告警。

## 运行说明

验证于 2026-10-06，版本见下表。这一节是视频里命令的唯一来源。

### 前置条件

- macOS。录制在 macOS 上做，其他系统没试过。
- Go 1.27 或以上（录制用 1.27.1），只用标准库。
- Node.js 20.19.0 或以上，OpenSpec 要求（录制用 24.15.0）。
- git。
- Claude Code，和一个能登录它的 Claude 账号。录制用的是 Claude Pro 订阅，模型选 Sonnet，权限模式是客户端缺省的 auto mode。

### 录制时的版本与安装来源

| 工具 | 版本 | 安装来源 |
|---|---|---|
| Claude Code | 2.1.289 | 官方安装脚本。装指定版本：`curl -fsSL https://claude.ai/install.sh \| bash -s 2.1.289`（[官方文档](https://code.claude.com/docs/en/setup#install-a-specific-version)） |
| OpenSpec | 1.14.0 | npm 包 `@fission-ai/openspec`（[官方安装文档](https://github.com/Fission-AI/OpenSpec/blob/main/docs/installation.md)） |
| Superpowers | 6.4.2（commit `8ca22dba9a94`） | Claude Code 里输入 `/plugin` 安装，来源 `anthropic-plugin-directory`（[仓库](https://github.com/obra/superpowers)） |

想装到和录制一样的版本：

```bash
npm install -g @fission-ai/openspec@1.14.0
openspec --version        # 1.14.0
```

视频画面上的两条是官方文档的写法，装到的是当时的最新版：

```bash
npm install -g @fission-ai/openspec@latest                # 2026-10-07 是 1.14.1
```

```
/plugin install superpowers@claude-plugins-official
```

第二条在 Claude Code 对话框里输入。官方市场 2026-10-07 指向的是 Superpowers 6.4.1，不是录制用的 6.4.2。6.4.1 的 `writing-plans` 会把代码整段写进计划，6.4.2 只写函数签名、测试断言和规格里的值（[RELEASE-NOTES v6.4.2](https://github.com/obra/superpowers/blob/main/RELEASE-NOTES.md)）。其余步骤两个版本一样。

### 步骤

一共 5 步。耗时是录制时的实际时长，你的会因模型、网络和回答多少而不同，按估算看。

**0. 准备**（约 15 分钟，估算，只做一次）

装好上面的工具，然后：

```bash
git clone https://github.com/laodao-ai/sdd-tdd-loop-demo.git
cd sdd-tdd-loop-demo
git switch -c my-run start
go test ./...             # 起点的测试能跑通
claude
```

`start` 里已经跑过 `openspec init`，也打开了 `/opsx:ff` 和 `/opsx:verify` 两个扩展命令，命令文件在 `.claude/commands/opsx/`，不用再开。在你自己的项目里，要先在终端跑 `openspec init`；扩展命令要用 `openspec config profile` 勾上，再跑 `openspec update`。

**1. OpenSpec：先聊清楚，再写规格**（录制约 5 分钟）

```
/opsx:explore Add an over-temperature alert.
```

AI 读完代码会问你：告警发到哪、阈值写死还是做成参数、正好 80 度算不算。录制时的回答是「B, hardcode 80.0, alert when over 80」。它问要不要记成一个 change 时，回答要。然后：

```
/opsx:ff add-over-temperature-alert
```

生成提案、规格、设计、任务四个文档，在 `openspec/changes/add-over-temperature-alert/`。不想开扩展命令的话，用默认的 `/opsx:propose` 一条命令出齐四个文档。

**2. Superpowers：写计划**（录制约 3 分钟）

```
/superpowers:writing-plans Use openspec/changes/add-over-temperature-alert/ (proposal, specs, design, tasks) as the spec, and break its tasks.md into implementation tasks.
```

计划写在 `docs/superpowers/plans/`。写完它会问用哪种方式执行，录制时没在这里回答，下一步用命令指定。

**3. Superpowers：按计划做**（录制约 12 分钟）

```
/superpowers:subagent-driven-development Execute the plan in docs/superpowers/plans/.
```

它会问要不要建 git worktree，按它的推荐答。每个任务派一个子代理先写测试、再写实现，做完派另一个子代理审。最后问怎么收尾，录制时选在本地合并回 main。

**4. 跑一遍，再用 OpenSpec 核对和归档**（录制约 4 分钟）

在 Claude Code 里用 `!` 直接跑命令：

```
! go test ./...
! printf '23.5\n81.2\n' | go run .
```

应该看到两个包 `ok`，输出 `23.5°C`、`81.2°C`，以及一行 `ALERT: 81.2°C > 80.0°C`。然后：

```
/opsx:verify
```

录制时它报了任务清单六项一项都没勾，还报了一条写错的测试。让它勾上任务、改掉那条，再：

```
/opsx:archive
```

问要不要把规格同步进主规格时，选 Sync now。

### 结果不会和视频一字不差

AI 每次的回答都不同。它问的问题、计划怎么拆、子代理用哪个模型、测试怎么写，都可能和视频不一样。录制时 OpenSpec 在任务清单里写错了一条测试，你复走时不一定会出现。

对照录制结果看思路，不要逐字比：

| tag | 到哪一步 |
|---|---|
| `start` | 起点：最小的读温度程序，OpenSpec 已初始化 |
| `seg1-done` | 第 1 步：四个规划文档 |
| `seg2-done` | 第 2 步：计划 |
| `seg3-done` | 第 3 步：代码、测试、README |
| `seg4-done` | 第 4 步：勾完任务，归档，规格并进 `openspec/specs/` |

```bash
git diff start seg4-done --stat
```

## Run (English)

```bash
go test ./...
printf '23.5\n81.2\n' | go run .
```

The program reads temperatures from stdin (one per line) and prints them to stdout with the unit. Readings strictly above 80.0°C also emit an alert to stderr: `ALERT: <reading>°C > 80.0°C`.

Requires Go 1.27 or later. Standard library only. The tag `start` marks the starting point; `seg1-done` to `seg4-done` mark each step of the recorded run.
