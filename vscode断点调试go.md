# vscode断点调试go

在vscode中除了正常安装go的插件外，需要编写launch.json

```json

{
    // 使用 IntelliSense 了解相关属性。 
    // 悬停以查看现有属性的描述。
    // 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "host": "127.0.0.1",
            "port": 2345,
            "program": "${fileDirname}",
            "env": {
                "GO111MODULE": "off",
            },
            "args": [],
            "showLog": true
        }
    ]
}
```

然后就可以正常的断点调试了