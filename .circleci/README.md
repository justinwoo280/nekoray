# CircleCI 配置说明

## 流量优化策略

### 1. 缓存机制
- **Go modules**: 缓存 `~/go/pkg/mod`，避免每次下载依赖
- **缓存 key**: 基于 `go.mod` checksum，依赖变化时才更新

### 2. 浅克隆
```bash
git clone --depth 1  # 只拉取最新提交，节省 90%+ 流量
```

### 3. Workspace 共享
- 在 jobs 间传递构建产物，避免重复构建
- 比 artifacts 更高效（不经过外部存储）

### 4. 条件构建
- **普通提交**: 只构建 Go core（2-3分钟）
- **Release 分支**: 完整构建（15-20分钟）
- **手动触发**: API 或 Web UI 触发完整构建

## 使用方式

### 快速验证（自动）
推送到任意分支 → 仅构建 Go core

### 完整构建
1. **方法1**: 推送到 `release/*` 分支
2. **方法2**: 打 tag `v1.0.0`
3. **方法3**: CircleCI Web UI 手动触发

## 预计资源消耗

| 构建类型 | 时间 | Credits |
|---------|------|---------|
| Go core only | 2-3分钟 | ~5 credits |
| 完整构建 | 15-20分钟 | ~40 credits |

**免费额度**: 6000 credits/月 ≈ 150次完整构建
