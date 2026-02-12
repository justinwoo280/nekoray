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

## Geoip/Geosite 配置迁移到 Rule-Set (已修复)

### Issue
旧版本使用 sing-box 1.8.0 之前的 geoip/geosite 配置方式，这些字段在 sing-box 1.12.0 中已完全移除。

### 原因
- **sing-box 1.8.0**: 废弃了规则中的 `geoip` 和 `geosite` 字段
- **sing-box 1.12.0**: 完全移除了这些字段的支持
- **新方式**: 使用 `rule_set` 配置远程规则集

### 修复状态
✅ **已修复** - ConfigBuilder 已更新为使用 rule_set 方式：
- `geoip:cn` → 自动生成 rule-set `geoip-cn` (远程 .srs 文件)
- `geosite:cn` → 自动生成 rule-set `geosite-cn` (远程 .srs 文件)

### 配置变化示例

**旧方式（已废弃）：**
```json
{
  "route": {
    "geoip": {"path": "geoip.db"},
    "geosite": {"path": "geosite.db"},
    "rules": [
      {"geosite": ["cn"], "outbound": "direct"},
      {"geoip": ["cn"], "outbound": "direct"}
    ]
  }
}
```

**新方式（当前实现）：**
```json
{
  "route": {
    "rule_set": [
      {
        "tag": "geosite-cn",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite-cn.srs",
        "download_detour": "proxy"
      },
      {
        "tag": "geoip-cn",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip-cn.srs",
        "download_detour": "proxy"
      }
    ],
    "rules": [
      {"rule_set": ["geosite-cn"], "outbound": "direct"},
      {"rule_set": ["geoip-cn"], "outbound": "direct"}
    ]
  }
}
```

### 优点
- ✅ **自动更新**: rule-set 会自动从远程下载最新数据
- ✅ **更小体积**: 不需要打包大型 .db 文件到发布版
- ✅ **更灵活**: 支持自定义规则集和版本控制

---

## 其他待迁移的 sing-box 废弃配置 (未来需要修复)

### 1. Inbound Domain Strategy (计划在 sing-box 1.13.0 移除)

**当前使用位置**: ConfigBuilder.cpp lines 414, 435

**废弃配置**:
```json
{
  "inbounds": [{
    "type": "mixed",
    "domain_strategy": "prefer_ipv4"
  }]
}
```

**新配置** (使用 route action):
```json
{
  "inbounds": [{
    "type": "mixed",
    "tag": "in"
  }],
  "route": {
    "rules": [{
      "inbound": "in",
      "action": "resolve",
      "strategy": "prefer_ipv4"
    }]
  }
}
```

### 2. Outbound Domain Strategy (计划在 sing-box 1.14.0 移除)

**当前使用位置**: ConfigBuilder.cpp line 349

**废弃配置**:
```json
{
  "outbounds": [{
    "type": "socks",
    "domain_strategy": "prefer_ipv4"
  }]
}
```

**新配置** (使用 domain_resolver):
```json
{
  "dns": {
    "servers": [{
      "type": "local",
      "tag": "local"
    }]
  },
  "outbounds": [{
    "type": "socks",
    "domain_resolver": {
      "server": "local",
      "strategy": "prefer_ipv4"
    }
  }]
}
```

### 3. Sniff 配置 (计划在 sing-box 1.13.0 移除)

**当前使用位置**: ConfigBuilder.cpp lines 403-404, 432-433

**废弃配置**:
```json
{
  "inbounds": [{
    "type": "mixed",
    "sniff": true,
    "sniff_override_destination": true
  }]
}
```

**新配置** (使用 route action):
```json
{
  "inbounds": [{
    "type": "mixed",
    "tag": "in"
  }],
  "route": {
    "rules": [{
      "inbound": "in",
      "action": "sniff"
    }]
  }
}
```

### 注意事项
这些配置项目前仍然可用，但将在未来版本中移除。建议在 sing-box 1.13.0 发布前进行迁移。

**参考文档**: https://sing-box.sagernet.org/migration/
