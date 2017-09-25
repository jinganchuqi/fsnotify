#fsnotify
基于fsnotify的文件事件通知

```json

    {
        "reload":"taskkill -F /IM main.exe && start main.exe",
        "depth":20,
        "path": "F:\\www\\test",
        "isReturn":false,
        "except": [
          ".idea",
          ".git",
          "node_modules",
    	  "vendor"
        ],
        "callback": "rsync -aP  -vzrtopg --progress --exclude-from '/cygdrive/d/rsync/config/exclude.txt' --delete 172.16.1.117::test /cygdrive/z/app/test"
      }

```

| 参数名     |    解释 |   类型   |
| :--------:| :--------: | :------: |
| reload | 重新加载配置时的回调 | string |
| depth | 侦听目录深度 | int |
| path | 侦听目录 | string |
| isReturn | 是否有返回值 | bool |
| except | 排除目录或文件 | array |
| callback | 事件回调 | string |
