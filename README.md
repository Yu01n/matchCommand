# matchCommand
Match the Linux/Windows operating system command and return the help information of the corresponding option (call go with python)

# Usage
This `main.go` uses python to call, before calling, you need to use the following command to compile into a .so file
```
 go build -buildmode=c-shared -o matchCommand.so main.go
```
At this time, a ```.so``` file and a ```.h``` file will be generated, and then executed using ```test.py``` in test
