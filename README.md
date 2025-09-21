# 🎯 Probocis
a mini-Git re-implementation from scratch in Go
## 📚 Motivation

Git is a powerful but complex version control system.  
Understanding Git’s internals — how it stores objects, manages the index, and builds commits — is often opaque for beginners.  

This project aims to:  
- Explore Git’s **plumbing commands** and underlying mechanisms  
- Learn how **blobs, trees, commits, and the index** are structured  
- Reimplement Git from scratch in Go for **educational purposes**  
- Provide a modular foundation to experiment with new Git features  

---

## ✨ Feature (in progress)

- Reimplement core Git plumbing commands:  
  - `init`  
  - `hash-object`  
  - `update-index` (add/remove files)  
  - `write-tree`  
  - `commit-tree`  
- Build minimal porcelain commands on top:  
  - `add` (via `update-index`)  
  - `rm --cached`  
  - `commit`
  - `clone`
- Modular design for objects, index, and refs  
- Supports content-addressable storage for blobs  
- Index entries are sorted by path (as in real Git)  

---

## 🛠️ Installation
(To Be Continued)
