# i3 Workspaces

This is a small utility that runs commands on new workspace when they're created.

## Example

Setting workspace layout:

config.toml:

```toml
["workspace name"]
commands = [
    "layout tabbed"
]
```

Exec line in i3 config:

```
exec ~/bin/i3-workspaces ~/.config/i3-workspaces/config.toml
```
