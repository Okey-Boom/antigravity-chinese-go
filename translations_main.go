package main

import (
	"strings"
)

// TranslateStats 翻译统计
type TranslateStats struct {
	NormalCount   int
	TemplateCount int
	VariableCount int
}

// normalTranslationsMain main.js 的普通翻译规则
var normalTranslationsMain = map[string]string{
	`"Request Review"`: `"请求确认"`,
	`"Always Proceed"`: `"始终继续"`,
	`"Agent Decides"`: `"由代理决定"`,
	`"Enable Telemetry"`: `"启用遥测"`,
	`"Allow List Terminal Commands"`: `"终端命令允许列表"`,
	`"Deny List Terminal Commands"`: `"终端命令拒绝列表"`,
	`"Agent Auto-Fix Lints"`: `"代理自动修复代码问题"`,
	`"Explain and Fix in Current Conversation"`: `"在当前对话中解释和修复"`,
	`"Secure Mode"`: `"安全模式"`,
	`"Agent Gitignore Access"`: `"代理访问 Gitignore"`,
	`"Agent Non-Workspace File Access"`: `"代理访问工作区外文件"`,
	`"Terminal Command Auto Execution"`: `"终端命令自动执行"`,
	`"Review Policy"`: `"审核策略"`,
	`"Auto-Continue"`: `"自动继续"`,
	`"Conversation History"`: `"会话历史"`,
	`"Knowledge"`: `"知识库"`,
	`"Auto-Open Edited Files"`: `"自动打开编辑的文件"`,
	`"Open Agent on Reload"`: `"重新加载时打开代理"`,
	`"Enable Sounds for Agent"`: `"启用 Agent 提示音"`,
	`"Suggestions in Editor"`: `"编辑器中的建议"`,
	`"Tab to Jump"`: `"Tab 跳转"`,
	`"Tab to Import"`: `"Tab 导入"`,
	`"Tab Speed"`: `"Tab 速度"`,
	`"Clipboard Context"`: `"剪贴板上下文"`,
	`"Highlight After Accept"`: `"接受后高亮"`,
	`"Tab Gitignore Access"`: `"Tab 访问 Gitignore"`,
	`"Enable Browser Tools"`: `"启用浏览器工具"`,
	`"Browser Javascript Execution Policy"`: `"浏览器 Javascript 执行策略"`,
	`"Chrome Binary Path"`: `"Chrome 二进制文件路径"`,
	`"Browser User Profile Path"`: `"浏览器用户配置文件路径"`,
	`"Browser CDP Port"`: `"浏览器 CDP 端口"`,
	`"Enable Demo Mode (Beta)"`: `"启用演示模式 (Beta)"`,
	`"Marketplace Item URL"`: `"扩展市场项目 URL"`,
	`"Marketplace Gallery URL"`: `"扩展市场库 URL"`,
	`"Browser URL Allowlist"`: `"浏览器 URL 允许列表"`,
	`"[Dev] GCP Project ID"`: `"[开发] GCP 项目 ID"`,
	`"General"`: `"常规"`,
	`"Security"`: `"安全性"`,
	`"Artifact"`: `"工件"`,
	`"Terminal"`: `"终端"`,
	`"File Access"`: `"文件访问"`,
	`"Automation"`: `"自动化"`,
	`"History"`: `"历史"`,
	`"Suggestions"`: `"建议"`,
	`"Navigation"`: `"导航"`,
	`"Context"`: `"上下文"`,
	`"Advanced"`: `"高级"`,
	`"Account"`: `"账户"`,
	`"Allowlist"`: `"允许列表"`,
	`"Marketplace"`: `"市场"`,
	`"Email"`: `"邮箱"`,
	`"Sign out"`: `"退出登录"`,
	`"Not signed in"`: `"未登录"`,
	`"Sign in"`: `"登录"`,
	`"Terms of Service"`: `"服务条款"`,
	`"Advanced settings"`: `"高级设置"`,
	`"Add"`: `"添加"`,
	`"Editor Settings"`: `"编辑器设置"`,
	`"Open Editor Settings"`: `"打开编辑器设置"`,
	`"Notification Settings"`: `"通知设置"`,
	`"Open System Preferences"`: `"打开系统偏好设置"`,
	`"Your Plan: "`: `"当前套餐: "`,
	`return"Always Proceed";default:return"Request Review"`: `return"始终继续";default:return"请求确认"`,
	`"Fast"`: `"快速"`,
	`"Slow"`: `"慢速"`,
	`"Disabled"`: `"已禁用"`,
}

// templateTranslationsMain main.js 的模板翻译规则
var templateTranslationsMain = [][2]string{
	{`"Always ask for permission"`, `"始终请求权限"`},
	{`"Always run terminal commands"`, `"始终运行终端命令"`},
	{`"Agent never asks for review. This maximizes the autonomy of the Agent, but also has the highest risk of the Agent operating over unsafe or injected Artifact content."`, `"代理从不请求确认。这最大化了代理的自主性，但也存在代理操作不安全或被注入的工件内容的最高风险。"`},
	{`"Agent will decide when to ask for review based on task complexity and user preference."`, `"代理将根据任务复杂性和用户偏好决定何时请求确认。"`},
	{`"Agent always asks for review."`, `"代理始终请求确认。"`},
	{`"Agent can plan before executing tasks. Use for deep research, complex tasks, or collaborative work"`, `"代理可在执行任务前进行规划。适用于深度研究、复杂任务或协作工作"`},
	{`"Agent will execute tasks directly. Use for simple tasks that can be completed faster"`, `"代理将直接执行任务。适用于可快速完成的简单任务"`},
	{`"When enabled, Agent is given awareness of lint errors created by its edits and may fix them without explicit user prompting."`, `"启用后，代理会自动感知其编辑产生的代码检查错误，并可在无需用户明确提示的情况下自动修复。"`},
	{`"When enabled, enforces settings that prevent the agent from autonomously running targeted exploits and requires human review for all agent actions. Visit antigravity.google/docs/secure-mode for details."`, `"启用后，将强制执行防止代理自动运行针对性漏洞利用的设置，并要求对所有代理操作进行人工审核。访问 antigravity.google/docs/secure-mode 了解详情。"`},
	{`"Allow Agent to view and edit the files in .gitignore automatically. Use with caution if your .gitignore lists files containing credentials, secrets, or other sensitive information."`, `"允许代理自动查看和编辑 .gitignore 中列出的文件。请谨慎使用：如果您的 .gitignore 包含凭据、密钥或其他敏感信息，请务必小心。"`},
	{`"Allow Agent to view and edit files outside of the current workspace automatically. Use with caution: this provides the Agent access to additional potentially-relevant information, but also allows the Agent to access credential files, secrets, and other files outside of the workspace that could be targeted in prompt injection attacks or other exploits by malicious actors."`, `"允许 Agent 自动查看和编辑当前工作区之外的文件。请谨慎使用：此选项可让 Agent 获取更多可能相关的信息，但也意味着 Agent 能够访问工作区外的凭据文件、密钥及其他敏感文件，这些文件可能被恶意攻击者通过提示词注入或其他手段加以利用。"`},
	{`"When enabled, Agent will automatically continue its response when it reaches its per-response invocation limit. If this setting is off, Agent will instead prompt you to continue upon reaching the limit."`, `"启用后，代理在达到单次响应调用限制时会自动继续。禁用此设置时，代理会在达到限制时提示您是否继续。"`},
	{`"GCP Project ID for enterprise features."`, `"企业功能的 GCP 项目 ID。"`},
	{`"When enabled, the agent will be able to access past conversations to inform its responses."`, `"启用后，Agent 可以访问过往对话记录，以此作为回答的参考依据。"`},
	{`"When enabled, the agent will be able to access its knowledge base to inform its responses and automatically generate knowledge items in the background. Disabling this will prevent the agent from accessing existing knowledge items, but will not delete them."`, `"启用后，代理可以访问其知识库以辅助回答，并在后台自动生成知识条目。禁用此选项会阻止代理访问现有知识条目，但不会删除它们。"`},
	{`"Open files in the background if Agent creates or edits them"`, `"代理创建或编辑文件时在后台自动打开"`},
	{`"Open Agent panel on window reload"`, `"窗口重新加载时打开代理面板"`},
	{`"When enabled, Antigravity will play a sound when Agent finishes generating a response."`, `"启用后，当 Agent 生成完回复时，Antigravity 将播放提示音。"`},
	{`"Show suggestions when typing in the editor"`, `"在编辑器中输入时显示建议"`},
	{`"Predict the location of your next edit and navigates you there with a tab keypress."`, `"预测您下一个编辑位置并通过 Tab 键跳转到那里。"`},
	{`"Quickly add and update imports with a tab keypress."`, `"通过 Tab 键快速添加和更新导入语句。"`},
	{`"Set the speed of tab suggestions"`, `"设置 Tab 建议的速度"`},
	{`"Highlight newly inserted text after accepting a Tab completion."`, `"接受 Tab 补全后高亮显示新插入的文本。"`},
	{`"Allow Tab to view and edit the files in .gitignore. Use with caution if your .gitignore lists files containing credentials, secrets, or other sensitive information."`, `"允许 Tab 查看和编辑 .gitignore 中列出的文件。请谨慎使用：如果您的 .gitignore 包含凭据、密钥或其他敏感信息，请务必小心。"`},
	{`"When enabled, Agent can use browser tools to open URLs, read web pages, and interact with browser content. This allows the Agent access to important (and often critical) knowledge and methods of validation, but any browser integration does increase exposure to external malicious parties for security exploits."`, `"启用后，代理可以使用浏览器工具打开 URL、读取网页并与浏览器内容交互。这能让代理访问重要（甚至关键）的知识和验证方法，但任何浏览器集成都会增加遭受外部恶意攻击的安全风险。"`},
	{`"Control which URLs the browser can access. Add domains or full URLs to the allowlist."`, `"控制浏览器可访问的 URL。可将域名或完整 URL 添加到白名单中。"`},
	{`"To modify editor settings, open Settings within the editor window."`, `"要修改编辑器设置，请在编辑器窗口中打开设置。"`},
	{`"To modify notification settings, open your operating system's system preferences."`, `"要修改通知设置，请打开操作系统的系统偏好设置。"`},
	{`"By using this app, you agree to its"`, `"使用本应用即表示你同意其"`},
	{`"Agent auto-executes commands matched by an allow list entry. For Unix shells, an allow list entry matches a command if its space-separated tokens form a prefix of the command's tokens. For PowerShell, the entry tokens may match any contiguous subsequence of the command tokens."`, `"代理会自动执行与允许列表匹配的命令。对于 Unix Shell，允许列表条目通过空格分隔的命令前缀进行匹配；对于 PowerShell，允许列表条目可以匹配命令中任意连续的子序列。"`},
	{`"Agent asks for permission before executing commands matched by a deny list entry. The deny list follows the same matching rules as the allow list and takes precedence over the allow list."`, `"代理在执行与拒绝列表匹配的命令前会请求您的授权。拒绝列表使用与允许列表相同的匹配规则，且优先级高于允许列表。"`},
	{`upgradeButtonText||"Upgrade"`, `upgradeButtonText||"升级"`},
	{"\\u2022 Disabled - Agent will never run Javascript code in the browser.", "\\u2022 已禁用 - 代理永远不会在浏览器中运行 Javascript 代码。"},
	{"\\u2022 Request Review - Agent will always stop to ask for permission to run Javascript code in the browser.", "\\u2022 请求确认 - 代理在浏览器中运行 Javascript 代码前会始终请求您的许可。"},
	{"\\u2022 Always Proceed - Agent will not stop to ask for permission to run Javascript in the browser. This provides the Agent with maximum autonomy to perform complex actions and validation in the browser, but also has the highest exposure to security exploits.", "\\u2022 始终继续 - 代理在浏览器中运行 Javascript 时不会请求许可。这为代理提供了在浏览器中执行复杂操作和验证的最大自主权，但也面临最高的安全漏洞风险。"},
	{`"Path to the Chrome/Chromium executable. Leave empty for auto-detection."`, `"Chrome/Chromium 可执行文件的路径。留空以自动检测。"`},
	{`"Custom path for the browser user profile directory. Leave empty for default (~/.gemini/antigravity-browser-profile)."`, `"浏览器用户配置文件目录的自定义路径。留空以使用默认值 (~/.gemini/antigravity-browser-profile)。"`},
	{`"Port number for Chrome DevTools Protocol remote debugging. Leave empty for default (9222)."`, `"Chrome DevTools 协议远程调试的端口号。留空以使用默认值 (9222)。"`},
}

// variableTranslationsMain main.js 的变量翻译规则 (包含模板字面量)
var variableTranslationsMain = [][2]string{
	{"=`Specifies Agent's behavior when asking for review on artifacts, which are documents it creates to enable a richer conversation experience.\n${", "=`指定 Agent 在请求用户审阅工件时的行为。工件是 Agent 创建的文档，用于提供更丰富的对话体验。\n${"},
	{"`When toggled on, ${e.product.nameShort} collects usage data to help Google enhance performance and features.`", "`启用后，${e.product.nameShort} 会收集使用数据以帮助 Google 提升性能和功能。`"},
	{"`When enabled, ${Rje.name} will use the clipboard as context for completions. May increase exposure to security exploits based on unintentional contents in clipboard.`", "`启用后，${Rje.name} 将使用剪贴板内容作为自动补全的上下文。若剪贴板中无意间包含敏感内容，可能会增加遭遇安全利用的风险。`"},
	{"`When enabled, ${QQe.name} will use the clipboard as context for completions. May increase exposure to security exploits based on unintentional contents in clipboard.`", "`启用后，${QQe.name} 将使用剪贴板内容作为自动补全的上下文。若剪贴板中无意间包含敏感内容，可能会增加遭遇安全利用的风险。`"},
	{"`Changes the base URL on each extension page. You must restart ${e.nameShort} to use the new marketplace after changing this value.`", "`更改每个扩展页面的基础 URL。修改此值后，您必须重启 ${e.nameShort} 才能使用新的扩展市场。`"},
	{"`Changes the base URL for marketplace search results. You must restart ${e.nameShort} to use the new marketplace after changing this value.`", "`更改扩展市场搜索结果的基础 URL。修改此值后，您必须重启 ${e.nameShort} 才能使用新的扩展市场。`"},
	{"`\\u2022 Always Proceed - Agent never asks for confirmation before executing terminal commands (except those in the Deny list). This provides the Agent with the maximum ability to operate over long periods without intervention, but also has the highest risk of an Agent executing an unsafe terminal command.\n        \\u2022 Request Review - Agent always asks for confirmation before executing terminal commands (except those in the Allow list).\n\n        Note: A change to this setting will only apply to new messages sent to Agent. In-progress responses will use the previous setting value.\n        `", "`\\u2022 始终继续 - 代理在执行终端命令之前从不请求确认（拒绝列表中的除外）。这为代理提供了长时间无干预运作的最大能力，但也存在代理执行不安全终端命令的最高风险。\n        \\u2022 请求确认 - 代理在执行终端命令之前始终请求确认（允许列表中的除外）。\n\n        注意：此设置的更改仅适用于发送给代理的新消息。正在进行的响应将使用之前的设置值。\n        `"},
	{"'When enabled, \"Explain and Fix\" actions will continue in the current conversation instead of starting a new one.'", "'启用后，\"解释并修复\"操作将在当前对话中继续进行，而不会另起新对话。'"},
	{"'When enabled, your UI will be slightly modified to ensure more consistent demos. This is only recommended for demo purposes. In most cases, you can run \"Antigravity: Start Demo Mode\" and \"Antigravity: Stop Demo Mode\" to control this switch and update your ~/.gemini/antigravity data directory.'", "'启用后，界面将进行微调以确保演示效果更加一致。此选项仅建议在演示场景下使用。通常情况下，你可以运行 \"Antigravity: Start Demo Mode\" 和 \"Antigravity: Stop Demo Mode\" 来控制此开关并更新你的 ~/.gemini/antigravity 数据目录。'"},
}

// applyMainTranslations 应用 main.js 的翻译规则
func applyMainTranslations(content string) (string, TranslateStats) {
	stats := TranslateStats{}

	// 1. 应用普通翻译
	for from, to := range normalTranslationsMain {
		if strings.Contains(content, from) {
			content = strings.ReplaceAll(content, from, to)
			stats.NormalCount++
		}
	}

	// 2. 应用模板翻译
	for _, pair := range templateTranslationsMain {
		from, to := pair[0], pair[1]
		if strings.Contains(content, from) {
			content = strings.ReplaceAll(content, from, to)
			stats.TemplateCount++
		}
	}

	// 3. 应用变量翻译
	for _, pair := range variableTranslationsMain {
		from, to := pair[0], pair[1]
		if strings.Contains(content, from) {
			content = strings.ReplaceAll(content, from, to)
			stats.VariableCount++
		}
	}

	return content, stats
}
