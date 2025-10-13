# Installing Go (Golang)

This guide explains how to install Go (Golang) on macOS, Linux, or Windows.

---

## 🧰 1. Check if Go is already installed

```bash
go version
```

If you get something like:
```
go version go1.23.3 darwin/amd64
```
then Go is already installed.  
If you see “command not found,” continue below.

---

## 🍏 macOS

### Option 1 — **Install via Homebrew** (recommended)
```bash
brew install go
```

Then verify:
```bash
go version
```

### Option 2 — **Manual install**
1. Download the latest `.pkg` from:  
   👉 [https://go.dev/dl/](https://go.dev/dl/)
2. Run the installer (it installs Go to `/usr/local/go`)
3. Add Go to your PATH (if not already):
   ```bash
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
   source ~/.zshrc
   ```

---

## 🐧 Linux

### Option 1 — **Using tarball (official method)**
```bash
wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
```

Add Go to your PATH:
```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Verify:
```bash
go version
```

### Option 2 — **Using your distro’s package manager**
(Ubuntu/Debian)
```bash
sudo apt install golang
```
Note: This may install an **older version**.

---

## 🪟 Windows

1. Download the Windows `.msi` installer from:  
   👉 [https://go.dev/dl/](https://go.dev/dl/)
2. Run the installer (it sets up Go in `C:\Go`)
3. Restart your terminal or PowerShell.
4. Verify installation:
   ```powershell
   go version
   ```

---

## ✅ 2. Test your Go installation

```bash
go version
go env GOPATH
```

You should see something like:
```
/Users/username/go
```

Create a test program:
```bash
mkdir -p ~/go/src/hello
cd ~/go/src/hello
nano main.go
```

Paste this:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

Run it:
```bash
go run main.go
```

You should see:
```
Hello, Go!
```

---

_This guide was generated on 2025-10-13_
