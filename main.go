package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	version       = "3.3"
	backupDirName = "antigravity_backup"
)

// 文件信息
type FileInfo struct {
	RelPath     string // 相对于安装目录的路径
	Description string // 文件描述
	Type        string // 文件类型 ("main", "chat", "continue")
}

// 备份记录
type BackupRecord struct {
	Timestamp   string            `json:"timestamp"`
	InstallPath string            `json:"install_path"`
	BackupType  string            `json:"backup_type"` // "antigravity" 或 "continue"
	Files       map[string]string `json:"files"`       // 原始路径 -> 备份文件名
}

// 需要汉化的文件列表 - Antigravity
var targetFilesAntigravity = []FileInfo{
	{
		RelPath:     `resources\app\out\jetskiAgent\main.js`,
		Description: "设置页 (主文件)",
		Type:        "main",
	},
	{
		RelPath:     `resources\app\out\vs\workbench\workbench.desktop.main.js`,
		Description: "设置页 (工作台)",
		Type:        "main",
	},
	{
		RelPath:     `resources\app\extensions\antigravity\out\media\chat.js`,
		Description: "聊天页",
		Type:        "chat",
	},
}

func main() {
	printBanner()

	// 显示主菜单
	for {
		choice := showMainMenu()
		switch choice {
		case "1":
			runAntigravityTranslation()
		case "2":
			runContinueTranslation()
		case "3":
			runRestore()
		case "4":
			showBackupList()
		case "0", "q", "Q":
			fmt.Println("\n👋 再见！")
			return
		default:
			fmt.Println("\n❌ 无效的选择，请重新输入")
		}
	}
}

func printBanner() {
	fmt.Println("╔═══════════════════════════════════════════════════╗")
	fmt.Printf("║   Antigravity 汉化工具 v%s (Go 语言版)           ║\n", version)
	fmt.Println("║   支持 Antigravity + Continue 扩展                ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝")
}

func showMainMenu() string {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📋 主菜单:")
	fmt.Println("   1. 🌐 汉化 Antigravity (主程序)")
	fmt.Println("   2. 🔧 汉化 Continue 扩展")
	fmt.Println("   3. ♻️  一键还原")
	fmt.Println("   4. 📂 查看备份列表")
	fmt.Println("   0. 🚪 退出")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Print("请选择 (1/2/3/4/0): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ========================================
// Antigravity 汉化功能
// ========================================

func runAntigravityTranslation() {
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("🌐 Antigravity 汉化模式")
	fmt.Println(strings.Repeat("═", 50))

	fmt.Println("\n🎯 本工具将自动汉化以下文件:")
	for _, f := range targetFilesAntigravity {
		fmt.Printf("   • %s\n", f.Description)
		fmt.Printf("     %s\n", f.RelPath)
	}

	// 自动检测 Antigravity 安装路径
	var installPath string
	detectedPath := findAntigravityInstallPath()

	if detectedPath != "" {
		fmt.Printf("\n✓ 自动检测到 Antigravity 安装路径:\n")
		fmt.Printf("   %s\n", detectedPath)
		fmt.Print("\n使用此路径？(Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "" || input == "y" || input == "yes" {
			installPath = detectedPath
		}
	}

	// 如果自动检测失败或用户拒绝，手动输入
	if installPath == "" {
		installPath = getInstallPath("Antigravity")
	}

	// 验证路径
	if !validateAntigravityPath(installPath) {
		fmt.Println("\n❌ 无效的 Antigravity 安装路径！")
		fmt.Println("   请确保路径中包含 resources\\app 目录")
		waitForKeypress()
		return
	}

	fmt.Printf("\n✓ 确认安装路径: %s\n", installPath)

	// 检测文件并显示状态
	foundFiles := detectAntigravityFiles(installPath)

	if len(foundFiles) == 0 {
		fmt.Println("\n❌ 未找到任何可汉化的文件！")
		fmt.Println("   请检查 Antigravity 是否正确安装")
		waitForKeypress()
		return
	}

	fmt.Printf("\n📋 找到 %d 个可汉化的文件:\n", len(foundFiles))
	for i, f := range foundFiles {
		fmt.Printf("   %d. %s (%s)\n", i+1, f.Description, f.RelPath)
	}

	// 询问是否继续
	fmt.Print("\n是否开始汉化？(Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		fmt.Println("已取消操作")
		return
	}

	// 创建备份目录
	backupDir, err := createBackupDir("antigravity")
	if err != nil {
		fmt.Printf("\n❌ 创建备份目录失败: %v\n", err)
		waitForKeypress()
		return
	}
	fmt.Printf("\n📁 备份目录: %s\n", backupDir)

	// 创建备份记录
	record := BackupRecord{
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		InstallPath: installPath,
		BackupType:  "antigravity",
		Files:       make(map[string]string),
	}

	// 开始汉化
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("🚀 开始汉化...")
	fmt.Println(strings.Repeat("─", 50))

	successCount := 0
	for _, f := range foundFiles {
		fullPath := filepath.Join(installPath, f.RelPath)
		fmt.Printf("\n📁 处理文件: %s\n", f.Description)
		fmt.Printf("   路径: %s\n", fullPath)

		// 备份文件
		backupFileName, err := createBackup(fullPath, backupDir)
		if err != nil {
			fmt.Printf("   ❌ 备份失败: %v\n", err)
			continue
		}
		record.Files[fullPath] = backupFileName
		fmt.Printf("   ✓ 备份已创建: %s\n", backupFileName)

		// 读取文件
		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("   ❌ 读取失败: %v\n", err)
			continue
		}
		originalSize := len(content)
		fmt.Printf("   📊 文件大小: %.2f MB\n", float64(originalSize)/1024/1024)

		// 应用翻译
		var translated string
		var stats TranslateStats
		if f.Type == "main" {
			translated, stats = applyMainTranslations(string(content))
		} else {
			translated, stats = applyChatTranslations(string(content))
		}

		// 保存文件
		err = os.WriteFile(fullPath, []byte(translated), 0644)
		if err != nil {
			fmt.Printf("   ❌ 保存失败: %v\n", err)
			continue
		}

		sizeDiff := len(translated) - originalSize
		diffSign := "+"
		if sizeDiff < 0 {
			diffSign = ""
		}

		fmt.Printf("   ✓ 翻译完成！\n")
		fmt.Printf("     - 普通翻译: %d 条\n", stats.NormalCount)
		fmt.Printf("     - 模板翻译: %d 条\n", stats.TemplateCount)
		if stats.VariableCount > 0 {
			fmt.Printf("     - 变量翻译: %d 条\n", stats.VariableCount)
		}
		fmt.Printf("     - 文件大小变化: %s%d 字节\n", diffSign, sizeDiff)

		successCount++
	}

	// 备份 product.json
	productJsonPath := filepath.Join(installPath, "resources", "app", "product.json")
	if _, err := os.Stat(productJsonPath); err == nil {
		backupFileName, err := createBackup(productJsonPath, backupDir)
		if err == nil {
			record.Files[productJsonPath] = backupFileName
		}
	}

	// 保存备份记录
	saveBackupRecord(backupDir, record)

	// 处理 product.json 校验和
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("🔧 移除 product.json 校验和...")
	removeProductJsonChecksums(installPath)

	// 显示结果
	fmt.Println("\n" + strings.Repeat("═", 50))
	if successCount == len(foundFiles) {
		fmt.Println("║         ✅ 全部汉化完成！                        ║")
	} else {
		fmt.Printf("║  ⚠️ 汉化完成 (%d/%d 成功)                         ║\n", successCount, len(foundFiles))
	}
	fmt.Println(strings.Repeat("═", 50))

	fmt.Println("\n💡 提示:")
	fmt.Println("   1. 请完全关闭并重新打开 Antigravity 以应用汉化")
	fmt.Println("   2. 备份已保存，可随时使用 [3] 一键还原")

	waitForKeypress()
}

// ========================================
// Continue 扩展汉化功能
// ========================================

func runContinueTranslation() {
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("🔧 Continue 扩展汉化模式")
	fmt.Println(strings.Repeat("═", 50))

	fmt.Println("\n📍 Continue 扩展路径格式:")
	fmt.Println("   C:\\Users\\{用户名}\\.antigravity\\extensions\\")
	fmt.Println("   continue.continue-{版本号}-win32-x64\\gui\\assets\\index.js")

	// 自动查找 Continue 扩展
	continueDir, indexPath := findContinueExtension()

	if indexPath != "" {
		fmt.Printf("\n✓ 自动检测到 Continue 扩展:\n")
		fmt.Printf("   目录: %s\n", filepath.Base(continueDir))
		fmt.Printf("   文件: %s\n", indexPath)

		fmt.Print("\n使用检测到的路径？(Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "" && input != "y" && input != "yes" {
			indexPath = ""
		}
	}

	if indexPath == "" {
		fmt.Print("\n请输入 index.js 的完整路径: ")
		reader := bufio.NewReader(os.Stdin)
		indexPath, _ = reader.ReadString('\n')
		indexPath = strings.TrimSpace(indexPath)
		indexPath = strings.Trim(indexPath, "\"'")
	}

	// 验证文件存在
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		fmt.Printf("\n❌ 文件不存在: %s\n", indexPath)
		waitForKeypress()
		return
	}

	fmt.Printf("\n✓ 确认文件路径: %s\n", indexPath)

	// 询问是否继续
	fmt.Print("\n是否开始汉化？(Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		fmt.Println("已取消操作")
		return
	}

	// 创建备份目录
	backupDir, err := createBackupDir("continue")
	if err != nil {
		fmt.Printf("\n❌ 创建备份目录失败: %v\n", err)
		waitForKeypress()
		return
	}
	fmt.Printf("\n📁 备份目录: %s\n", backupDir)

	// 创建备份记录
	record := BackupRecord{
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		InstallPath: filepath.Dir(filepath.Dir(filepath.Dir(indexPath))), // 保存扩展根目录
		BackupType:  "continue",
		Files:       make(map[string]string),
	}

	// 备份文件
	backupFileName, err := createBackup(indexPath, backupDir)
	if err != nil {
		fmt.Printf("\n❌ 备份失败: %v\n", err)
		waitForKeypress()
		return
	}
	record.Files[indexPath] = backupFileName
	fmt.Printf("   ✓ 备份已创建: %s\n", backupFileName)

	// 开始汉化
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("🚀 开始汉化...")
	fmt.Println(strings.Repeat("─", 50))

	// 读取文件
	content, err := os.ReadFile(indexPath)
	if err != nil {
		fmt.Printf("\n❌ 读取失败: %v\n", err)
		waitForKeypress()
		return
	}
	originalSize := len(content)
	fmt.Printf("   📊 文件大小: %.2f MB\n", float64(originalSize)/1024/1024)

	// 应用翻译
	translated, stats := applyContinueTranslations(string(content))

	// 保存文件
	err = os.WriteFile(indexPath, []byte(translated), 0644)
	if err != nil {
		fmt.Printf("\n❌ 保存失败: %v\n", err)
		waitForKeypress()
		return
	}

	// 保存备份记录
	saveBackupRecord(backupDir, record)

	sizeDiff := len(translated) - originalSize
	diffSign := "+"
	if sizeDiff < 0 {
		diffSign = ""
	}

	fmt.Printf("\n   ✓ 翻译完成！\n")
	fmt.Printf("     - 引号翻译: %d 条\n", stats.NormalCount)
	fmt.Printf("     - 全局替换: %d 条\n", stats.TemplateCount)
	fmt.Printf("     - 文件大小变化: %s%d 字节\n", diffSign, sizeDiff)

	// 显示结果
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("║         ✅ Continue 扩展汉化完成！               ║")
	fmt.Println(strings.Repeat("═", 50))

	fmt.Println("\n💡 提示:")
	fmt.Println("   1. 请完全关闭并重新打开 Antigravity 以应用汉化")
	fmt.Println("   2. 备份已保存，可随时使用 [3] 一键还原")

	waitForKeypress()
}

// findContinueExtension 自动查找 Continue 扩展
func findContinueExtension() (string, string) {
	// 获取用户目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", ""
	}

	// 查找扩展目录
	extensionsDir := filepath.Join(homeDir, ".antigravity", "extensions")
	if _, err := os.Stat(extensionsDir); os.IsNotExist(err) {
		return "", ""
	}

	// 查找 continue.continue-* 目录
	entries, err := os.ReadDir(extensionsDir)
	if err != nil {
		return "", ""
	}

	var latestVersion string
	var latestDir string

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "continue.continue-") {
			// 选择最新版本（按字符串排序）
			if entry.Name() > latestVersion {
				latestVersion = entry.Name()
				latestDir = filepath.Join(extensionsDir, entry.Name())
			}
		}
	}

	if latestDir == "" {
		return "", ""
	}

	// 检查 index.js 是否存在
	indexPath := filepath.Join(latestDir, "gui", "assets", "index.js")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return latestDir, ""
	}

	return latestDir, indexPath
}

// ========================================
// 还原功能
// ========================================

func runRestore() {
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("♻️  一键还原模式")
	fmt.Println(strings.Repeat("═", 50))

	// 获取备份目录
	programDir, err := os.Executable()
	if err != nil {
		fmt.Printf("\n❌ 获取程序目录失败: %v\n", err)
		waitForKeypress()
		return
	}
	programDir = filepath.Dir(programDir)
	backupBaseDir := filepath.Join(programDir, backupDirName)

	// 检查备份目录是否存在
	if _, err := os.Stat(backupBaseDir); os.IsNotExist(err) {
		fmt.Println("\n❌ 未找到任何备份！")
		fmt.Println("   备份目录不存在: " + backupBaseDir)
		waitForKeypress()
		return
	}

	// 列出所有备份
	backups, err := listBackups(backupBaseDir)
	if err != nil || len(backups) == 0 {
		fmt.Println("\n❌ 未找到任何备份记录！")
		waitForKeypress()
		return
	}

	fmt.Printf("\n📂 找到 %d 个备份:\n\n", len(backups))
	for i, b := range backups {
		backupTypeLabel := "Antigravity"
		if b.record.BackupType == "continue" {
			backupTypeLabel = "Continue 扩展"
		}
		fmt.Printf("   %d. [%s] %s\n", i+1, backupTypeLabel, b.dirName)
		fmt.Printf("      时间: %s\n", b.record.Timestamp)
		fmt.Printf("      路径: %s\n", b.record.InstallPath)
		fmt.Printf("      文件: %d 个\n\n", len(b.record.Files))
	}

	// 选择要还原的备份
	fmt.Printf("请选择要还原的备份 (1-%d，0 取消): ", len(backups))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "0" || input == "" {
		fmt.Println("已取消操作")
		return
	}

	var choice int
	fmt.Sscanf(input, "%d", &choice)
	if choice < 1 || choice > len(backups) {
		fmt.Println("\n❌ 无效的选择")
		waitForKeypress()
		return
	}

	selectedBackup := backups[choice-1]

	// 确认还原
	fmt.Printf("\n⚠️  即将还原备份: %s\n", selectedBackup.dirName)
	fmt.Printf("   目标路径: %s\n", selectedBackup.record.InstallPath)
	fmt.Printf("   将还原 %d 个文件\n", len(selectedBackup.record.Files))
	fmt.Print("\n确认还原？(y/N): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		fmt.Println("已取消操作")
		return
	}

	// 执行还原
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("🔄 开始还原...")
	fmt.Println(strings.Repeat("─", 50))

	successCount := 0
	for originalPath, backupFileName := range selectedBackup.record.Files {
		backupFilePath := filepath.Join(selectedBackup.fullPath, backupFileName)

		fmt.Printf("\n📁 还原文件: %s\n", filepath.Base(originalPath))

		// 检查备份文件是否存在
		if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
			fmt.Printf("   ❌ 备份文件不存在: %s\n", backupFileName)
			continue
		}

		// 读取备份文件
		content, err := os.ReadFile(backupFilePath)
		if err != nil {
			fmt.Printf("   ❌ 读取备份失败: %v\n", err)
			continue
		}

		// 写入原始位置
		err = os.WriteFile(originalPath, content, 0644)
		if err != nil {
			fmt.Printf("   ❌ 还原失败: %v\n", err)
			continue
		}

		fmt.Printf("   ✓ 已还原: %s\n", originalPath)
		successCount++
	}

	// 显示结果
	fmt.Println("\n" + strings.Repeat("═", 50))
	if successCount == len(selectedBackup.record.Files) {
		fmt.Println("║         ✅ 全部还原完成！                        ║")
	} else {
		fmt.Printf("║  ⚠️ 还原完成 (%d/%d 成功)                         ║\n", successCount, len(selectedBackup.record.Files))
	}
	fmt.Println(strings.Repeat("═", 50))

	fmt.Println("\n💡 提示:")
	fmt.Println("   请完全关闭并重新打开 Antigravity 以应用还原")

	waitForKeypress()
}

// ========================================
// 备份列表
// ========================================

func showBackupList() {
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("📂 备份列表")
	fmt.Println(strings.Repeat("═", 50))

	// 获取备份目录
	programDir, err := os.Executable()
	if err != nil {
		fmt.Printf("\n❌ 获取程序目录失败: %v\n", err)
		waitForKeypress()
		return
	}
	programDir = filepath.Dir(programDir)
	backupBaseDir := filepath.Join(programDir, backupDirName)

	fmt.Printf("\n📁 备份目录: %s\n", backupBaseDir)

	// 检查备份目录是否存在
	if _, err := os.Stat(backupBaseDir); os.IsNotExist(err) {
		fmt.Println("\n📭 暂无备份")
		waitForKeypress()
		return
	}

	// 列出所有备份
	backups, err := listBackups(backupBaseDir)
	if err != nil || len(backups) == 0 {
		fmt.Println("\n📭 暂无备份记录")
		waitForKeypress()
		return
	}

	fmt.Printf("\n找到 %d 个备份:\n\n", len(backups))
	for i, b := range backups {
		backupTypeLabel := "Antigravity"
		if b.record.BackupType == "continue" {
			backupTypeLabel = "Continue"
		}
		fmt.Printf("   %d. 📦 [%s] %s\n", i+1, backupTypeLabel, b.dirName)
		fmt.Printf("      创建时间: %s\n", b.record.Timestamp)
		fmt.Printf("      安装路径: %s\n", b.record.InstallPath)
		fmt.Printf("      备份文件:\n")
		for origPath, backupName := range b.record.Files {
			fmt.Printf("         • %s -> %s\n", filepath.Base(origPath), backupName)
		}
		fmt.Println()
	}

	waitForKeypress()
}

// ========================================
// 辅助函数
// ========================================

type backupInfo struct {
	dirName  string
	fullPath string
	record   BackupRecord
}

func listBackups(backupBaseDir string) ([]backupInfo, error) {
	var backups []backupInfo

	entries, err := os.ReadDir(backupBaseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(backupBaseDir, entry.Name())
		recordPath := filepath.Join(dirPath, "backup_record.json")

		// 读取备份记录
		if _, err := os.Stat(recordPath); os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(recordPath)
		if err != nil {
			continue
		}

		var record BackupRecord
		if err := json.Unmarshal(content, &record); err != nil {
			continue
		}

		backups = append(backups, backupInfo{
			dirName:  entry.Name(),
			fullPath: dirPath,
			record:   record,
		})
	}

	// 按时间倒序排列
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].dirName > backups[j].dirName
	})

	return backups, nil
}

func getInstallPath(appName string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n请输入 %s 安装路径: ", appName)
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)
	// 去掉可能的引号
	path = strings.Trim(path, "\"'")
	return path
}

func validateAntigravityPath(path string) bool {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// 检查是否包含 resources/app 目录
	resourcesPath := filepath.Join(path, "resources", "app")
	if _, err := os.Stat(resourcesPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// findAntigravityInstallPath 自动检测 Antigravity 安装路径
func findAntigravityInstallPath() string {
	// 1. 优先从注册表查询
	registryPath := findAntigravityFromRegistry()
	if registryPath != "" && validateAntigravityPath(registryPath) {
		return registryPath
	}

	// 2. 获取用户目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// 3. 检查常见安装位置
	candidates := []string{
		// 用户目录安装 (最常见)
		filepath.Join(homeDir, "AppData", "Local", "Programs", "Antigravity"),
		filepath.Join(homeDir, "AppData", "Local", "Antigravity"),
		// 系统目录安装
		"C:\\Program Files\\Antigravity",
		"C:\\Program Files (x86)\\Antigravity",
		// 其他常见位置
		"D:\\Antigravity",
		"D:\\Program Files\\Antigravity",
		"E:\\Antigravity",
	}

	for _, path := range candidates {
		if validateAntigravityPath(path) {
			return path
		}
	}

	return ""
}

// findAntigravityFromRegistry 从 Windows 注册表查询 Antigravity 安装路径
func findAntigravityFromRegistry() string {
	// 注册表查询位置
	registryPaths := []string{
		// 用户安装的程序
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall`,
		// 系统安装的程序 (64位)
		`HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		// 系统安装的程序 (32位 on 64位系统)
		`HKLM\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	// 收集所有有效路径
	var validPaths []string

	for _, regPath := range registryPaths {
		// 使用 reg query 命令查询注册表
		cmd := exec.Command("reg", "query", regPath, "/s", "/f", "Antigravity", "/d")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		// 解析输出
		lines := strings.Split(string(output), "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)

			// 查找 InstallLocation
			if strings.Contains(line, "InstallLocation") && strings.Contains(line, "REG_SZ") {
				parts := strings.SplitN(line, "REG_SZ", 2)
				if len(parts) == 2 {
					path := cleanRegistryPath(parts[1])
					if path != "" && validateAntigravityPath(path) {
						validPaths = append(validPaths, path)
					}
				}
			}

			// 查找 DisplayIcon (通常指向 exe 文件)
			if strings.Contains(line, "DisplayIcon") && strings.Contains(line, "REG_SZ") {
				parts := strings.SplitN(line, "REG_SZ", 2)
				if len(parts) == 2 {
					iconPath := cleanRegistryPath(parts[1])
					// 移除可能的逗号和图标索引
					if idx := strings.Index(iconPath, ","); idx > 0 {
						iconPath = iconPath[:idx]
					}
					// 获取目录路径
					dir := filepath.Dir(iconPath)
					if dir != "" && validateAntigravityPath(dir) {
						validPaths = append(validPaths, dir)
					}
				}
			}
		}
	}

	// 从有效路径中选择最佳匹配
	// 优先选择路径最短的（通常是主程序而不是子工具）
	if len(validPaths) == 0 {
		return ""
	}

	bestPath := validPaths[0]
	for _, p := range validPaths[1:] {
		if len(p) < len(bestPath) {
			bestPath = p
		}
	}

	return bestPath
}

// cleanRegistryPath 清理注册表返回的路径
func cleanRegistryPath(path string) string {
	path = strings.TrimSpace(path)
	// 移除引号
	path = strings.Trim(path, "\"")
	// 移除尾部反斜杠
	path = strings.TrimSuffix(path, "\\")
	return path
}

func detectAntigravityFiles(installPath string) []FileInfo {
	var found []FileInfo
	for _, f := range targetFilesAntigravity {
		fullPath := filepath.Join(installPath, f.RelPath)
		if _, err := os.Stat(fullPath); err == nil {
			found = append(found, f)
		}
	}
	return found
}

func createBackupDir(backupType string) (string, error) {
	// 获取程序运行目录
	programDir, err := os.Executable()
	if err != nil {
		return "", err
	}
	programDir = filepath.Dir(programDir)

	// 创建备份根目录
	backupBaseDir := filepath.Join(programDir, backupDirName)
	if err := os.MkdirAll(backupBaseDir, 0755); err != nil {
		return "", err
	}

	// 创建以时间和类型命名的子目录
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupDir := filepath.Join(backupBaseDir, fmt.Sprintf("%s_%s", timestamp, backupType))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	return backupDir, nil
}

func createBackup(filePath string, backupDir string) (string, error) {
	// 使用原始文件名
	fileName := filepath.Base(filePath)

	backupPath := filepath.Join(backupDir, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(backupPath, content, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func saveBackupRecord(backupDir string, record BackupRecord) {
	recordPath := filepath.Join(backupDir, "backup_record.json")
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		fmt.Printf("   ⚠️ 保存备份记录失败: %v\n", err)
		return
	}
	if err := os.WriteFile(recordPath, data, 0644); err != nil {
		fmt.Printf("   ⚠️ 保存备份记录失败: %v\n", err)
	}
}

func removeProductJsonChecksums(installPath string) {
	productJsonPath := filepath.Join(installPath, "resources", "app", "product.json")

	if _, err := os.Stat(productJsonPath); os.IsNotExist(err) {
		fmt.Println("   ⚠️ 未找到 product.json，跳过")
		return
	}

	content, err := os.ReadFile(productJsonPath)
	if err != nil {
		fmt.Printf("   ❌ 读取 product.json 失败: %v\n", err)
		return
	}

	originalContent := string(content)
	lines := strings.Split(originalContent, "\n")

	// 需要删除的校验和关键词
	checksumKeys := []string{
		`"jetskiAgent/main.js"`,
		`"vs/workbench/workbench.desktop.main.js"`,
	}

	removedCount := 0
	var newLines []string

	for _, line := range lines {
		skip := false
		for _, key := range checksumKeys {
			if strings.Contains(line, key) {
				removedCount++
				fmt.Printf("   ✓ 移除校验和: %s\n", key)
				skip = true
				break
			}
		}
		if !skip {
			newLines = append(newLines, line)
		}
	}

	if removedCount > 0 {
		// 保存修改后的内容
		newContent := strings.Join(newLines, "\n")
		// 修复尾随逗号问题
		newContent = strings.ReplaceAll(newContent, ",\n}", "\n}")
		newContent = strings.ReplaceAll(newContent, ",\n]", "\n]")

		err = os.WriteFile(productJsonPath, []byte(newContent), 0644)
		if err != nil {
			fmt.Printf("   ❌ 保存 product.json 失败: %v\n", err)
		} else {
			fmt.Printf("   ✓ product.json 已更新 (移除 %d 个校验和)\n", removedCount)
		}
	} else {
		fmt.Println("   ✓ 校验和已移除过，无需重复处理")
	}
}

func waitForKeypress() {
	fmt.Println()
	fmt.Print("按回车键继续...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
