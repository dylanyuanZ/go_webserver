Analyze the current changes and generate a professional commit message with Chinese description.

1.  **Context Gathering**: Check `git status`, `git diff`, and `git diff --cached` to understand all pending changes.
2.  **Style Alignment**: Use conventional commits format: `<type>(<scope>): <中文描述>`.
    - Types: `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `chore`
    - Scope: the module/package affected (e.g., `handler`, `service`, `repository`, `config`)
    - Description: **use Chinese** for the description part
3.  **Drafting**: Create concise **Chinese** commit message(s). The `type` and `scope` remain in English, but the `description` must be in Chinese. Include a short header and bullet points for complex changes.
4.  **Output**: Display the commit message(s) in code block(s). Provide the corresponding `git commit` command(s) for convenience.
