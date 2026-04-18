# 🌬️ WindDrop

> Instant file sharing from your terminal — local or global, no setup.

---

## ⚡ What is WindDrop?

WindDrop is a lightweight CLI tool to share files instantly.

* 📡 Share your files blazing fast .
* 🌍 Can generate public links (via Cloudflare tunnel)
* 🔐 Secure token-based access
* ⏳ Supports expiry & one-time downloads

👉 No accounts. No uploads. Just run and share.

---

## 🚀 Installation

```bash
git clone https://github.com/meetsoni555/winddrop.git
cd winddrop
chmod +x install.sh
./install.sh
```

---

## 🛠️ Usage

### 📤 Share a file (local network)

```bash
winddrop send file.zip
```

---

### 🌍 Share publicly (internet)

```bash
winddrop send file.zip --public
```

---

### ⏳ Share with expiry

```bash
winddrop send file.zip --expire 5m
```

---

### 🔥 One-time download

```bash
winddrop send file.zip --once
```

---

### 💀 Combined

```bash
winddrop send file.zip --public --once --expire 2m
```

---

## 🌐 How it works

```text
You run command → WindDrop starts server
→ Generates secure link
→ Receiver opens link in browser
→ File downloads instantly
```

---

## 📌 Example Output

```text
🌬️ WindDrop

File      : file.zip
Mode      : one-time
Expires   : one-time or 2m

Local Link  : http://192.168.x.x:8080/download?token=abc123
Public Link : https://xyz.trycloudflare.com/download?token=abc123
```

---

## 📡 Requirements

* Same WiFi for local sharing
* Internet required for `--public`

---

## 🚧 Status

Active development 🚀

---

## 📄 License

MIT
