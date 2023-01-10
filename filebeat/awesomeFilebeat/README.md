# awesomeFilebeat - Filebeat的日志收集配置动态更改

## Filebeat 实时加载机制

通过监视配置更改的路径(Glob)。当Glob找到的文件发生变化时，将根据配置文件中的更改启动和停止新的 prospetor。

```yaml
filebeat.config.prospectors:
  enabled: true
  path: usr/share/filebeat/filebeat.yml
  reload.enabled: true
  reload.period: 10s
```

- path: 一个Glob，用于定义要检查更改的文件
- reload.enabled: 设置true为时，启用动态配置重新加载
- reload.period: 指定检查文件更改的频率。不要将其设置period为小于1，因为文件的修改时间通常以秒为单位存储。将period小于1 设置为将导致不必要的开销。

## awesomeFilebeat

项目主要工作流程如下：

1. 在程序中指定需要动态更改的 Filebeat 配置文件路径 `watchYamlFilePath`

2. etcd 的 Watch 机制可以实时地监控到 etcd 中增量的数据更新。当发生变化时, awesomeFilebeat 将更改内容写入 filebeat.yml 的 filebeat.inputs 模块
3. Filebeat 启用动态配置重新加载，收集更新后的日志源。