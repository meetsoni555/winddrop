# WindDrop

Instant file sharing from your terminal over LAN.

---

<img width="1546" height="470" alt="cover" src="https://github.com/user-attachments/assets/a8241864-d61a-4950-9d68-e9828bcdfcd8" />

---

## What is WindDrop?

WindDrop is a lightweight CLI tool for fast file transfer on a local network.

* Share files and folders instantly
* No upload, no accounts, no setup
* Token-based secure access
* Supports expiry and one-time downloads

Run a command, get a link, open it on another device.

---

## Installation

```bash
git clone https://github.com/meetsoni555/winddrop.git
cd winddrop
chmod +x install.sh
./install.sh
```

---

## Usage

### Share a file

```bash
winddrop send file.zip
```

---

### Share a folder

```bash
winddrop send ~/Pictures
```

Folders are automatically compressed into a zip before sending.

---

### Share multiple files / folders

```bash
winddrop send file1.mp3 file2.jpg ~/Documents
```

All inputs are bundled into a single archive.

---

### Expiring link

```bash
winddrop send file.zip --expire 5m
```

The link becomes invalid after the specified duration.

---

### One-time download

```bash
winddrop send file.zip --once
```

The link works only once, then the server shuts down.

---

### Combined

```bash
winddrop send ~/Pictures --once --expire 2m
```

Whichever happens first:

* file is downloaded once
* time expires

---

## How it works

```text
You run command
→ WindDrop starts a local HTTP server
→ Generates a secure download link
→ Receiver opens the link in a browser
→ File downloads directly from your machine
```

---

## Example Output

```text
WindDrop

Mode      : Multi-file
Items     : 1
Archive   : winddrop_files.zip

Mode      : one-time
Expires   : 2m0s

Link : http://192.168.x.x:8080/download?token=abc123
```

---

## Requirements

* Devices must be on the same network (WiFi/LAN)
* Port 8080 should be available

---

## Limitations

* Works only on local network
* No resume support for interrupted downloads
* Large files depend on network speed

---

## Status

Stable LAN version
