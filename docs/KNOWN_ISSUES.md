# Known Issues

## gRPC Transport Configuration

### Issue
当使用 gRPC 传输层时，配置界面只有一个 `path` 字段（实际是 ServiceName），缺少独立的 `Host` 头设置。

### 问题详情
- **当前行为**: NekoRay 将服务器地址直接用作 HTTP Host 头
- **代码位置**: `rpc/gRPC.cpp:143` - `url_base = "http://" + url_;`
- **影响**: 当服务器使用 SNI/TLS 且 Host 与服务器地址不同时，TLS 验证会失败

### 预期配置
gRPC 传输应该分离以下字段：
1. **服务器地址**: 实际连接的 IP/域名
2. **Host 头**: HTTP/2 请求中的 `:authority` 伪头部（用于 SNI）
3. **ServiceName**: gRPC 服务路径（如 `GunService`）

### 临时解决方案
1. **使用不验证 TLS 的服务器**（不推荐）
2. **确保服务器地址与 TLS 证书的 SAN 一致**
3. **等待后续版本修复**

### 修复计划
需要修改以下文件：
- `rpc/gRPC.cpp`: 添加独立的 Host 头设置
- `ui/edit/dialog_edit_profile.ui`: 添加 Host 输入框
- `fmt/Bean2Link.cpp` 和 `Link2Bean.cpp`: 序列化/反序列化新字段

---

## Geosite.db Not Found (已修复)

### Issue
点击连接时显示 "Geosite.db Not Found" 错误。

### 原因
Windows 发布包中未包含地理位置数据库文件：
- `geosite.db`
- `geoip.db`
- `geosite.dat`
- `geoip.dat`

### 修复状态
✅ **已修复** - 从下一个 release 开始，这些文件会自动打包进 Windows 发布版。

### 手动修复（旧版本）
下载以下文件到 NekoRay 安装目录：
```bash
# geoip.dat
https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat

# geosite.dat
https://github.com/v2fly/domain-list-community/releases/latest/download/dlc.dat

# geoip.db (sing-box)
https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db

# geosite.db (sing-box)
https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db
```
