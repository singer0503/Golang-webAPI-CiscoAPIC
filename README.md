# Golang-webAPI-CiscoAPIC process description

![https://www.notion.so/image/https%3A%2F%2Fs3-us-west-2.amazonaws.com%2Fsecure.notion-static.com%2Fbbac9b95-123d-4368-a36c-dd10525269c5%2F-Golang.png?table=block&id=ed011c58-704f-4d82-93b0-9905e90e1206&width=3040&userId=e8aa0888-ca7b-4216-869b-435a8115d8eb&cache=v2](https://www.notion.so/image/https%3A%2F%2Fs3-us-west-2.amazonaws.com%2Fsecure.notion-static.com%2Fbbac9b95-123d-4368-a36c-dd10525269c5%2F-Golang.png?table=block&id=ed011c58-704f-4d82-93b0-9905e90e1206&width=3040&userId=e8aa0888-ca7b-4216-869b-435a8115d8eb&cache=v2)

- To call APIC Web API, you need to open a firewall, "PRTG Probe — 443 port —→ APIC Server"
- Compile Golang source code into EXE File
- Establish Custom EXE/Script Advanced Sensor of PRTG for timing drive call EXE File
- EXE File uses Web API to log in to obtain a token, then use the token to obtain CPU or Disk usage data
- Finally, the usage data is presented in the JSON format required by PRTG

# Source code directory structure

```
.
├── CPU-go
│   └── mian.go
├── Disk-go
│   └── main.go
└── README.md
```

You can see that there are two `main.go` files in the directory, where `mian.go` is usually the entry point of the main program for Golang

## Source code:

`~/CPU-go/main.go` monitors CPU usage source code

`~/Disk-go/main.go` monitor Disk usage source code

The two programs obtain different data, so PRTG also needs to create two Custom EXE/Script Advanced Sensors for monitoring.

## Use Go to generate an EXE file

Switch to the CPU-go or Disk-go directory. When you execute go build, the main.go file will be automatically searched for as the entry point of the program code. Compilation will generate EXE File, which is the finished product!

`~/CPU-go/CPU-go.exe` to monitor CPU usage

`~/Disk-go/Disk-go.exe` to monitor Disk usage
