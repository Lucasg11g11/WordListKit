# 🌩️ StormList - Wordlist Generation Toolkit

**StormList** is a powerful wordlist generation tool created for ethical hacking, red team simulations, and brute-force testing. Built with Go, it automatically creates high-quality password lists using smart generation techniques — and it's FAST ⚡.

---

## 📦 Features

- ✅ **Almost 17 million unique passwords generated** 💣  
- ✅ **Random strong passwords** (12–18 characters)  
- ✅ **Real-word based passwords** (e.g., `summer2024`)  
- ✅ **Name-based passwords** (e.g., `Lucas0987`)  
- ✅ **Avoids duplicate passwords**  
- ✅ **Multi-threaded generation (100 goroutines)**  
- ✅ **Cross-platform: works on Linux, macOS, and Windows**

---

## 📁 Files Included

- `StormList.txt` → Your evolving wordlist with over **17 million** unique entries  
- `dictionary.txt` → A list of real words used to make realistic passwords  
- `names.txt` → Optional file for generating passwords from common names  
- `main.go` → The Go source code of the StormList generator  
- `README.md` → You're reading it 😎  

---

## 🚀 How to Run

1. Use one of the two executables
2. Make sure you have **Go** installed. Then:

```bash
go run main.go
```

Choose a generation mode:

```
1 - random (strong random passwords)  
2 - real (based on dictionary words)  
3 - names (based on name list)
```

Passwords will be generated and saved into `StormList.txt` automatically.

---

## 🔧 Customize It

Want to change how passwords are built? Edit the functions:

- `generateRandomPassword(minLen, maxLen)`  
- `generateRealPassword(words)`  
- `generateNamePassword(names)`

You can also tweak the number of threads (`numThreads`) for faster generation.

---

## 🛡️ Warning

This tool is intended for **ethical hacking only**. Always get permission before using any wordlist in penetration tests.

---

## 🧠 Credits

Created by **XploitShade** – student, cybersecurity enthusiast, and creator of StormPulse & StormRecon.
