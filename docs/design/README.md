# 设计文档

每个需求在 `docs/design/<feature>/` 下建一个目录：

```
docs/design/<feature>/
├── spec.md     # 需求 + 设计方案（由 workflow-requirements-clarification / system-design 生成）
├── tasks.md    # 实施任务清单（由 workflow-code-generation 生成）
└── feedback_driven_report.md  # 调试类任务的迭代记录（按需）
```

`<feature>` 用小写连字符命名，例如 `user-register-api`。

历史设计文档保留在此目录，作为后续维护的上下文依据。
