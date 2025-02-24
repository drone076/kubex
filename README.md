![Build Status](https://github.com/drone076/kubex/actions/workflows/go.yml/badge.svg)

# Kubex

`kubex` is a lightweight CLI tool for managing Kubernetes contexts. It simplifies switching between Kubernetes contexts and viewing the current context. Think of it as a stripped-down version of `kubectx`, designed to be simple, fast, and cross-platform.

## Features
- List Contexts: View all available Kubernetes contexts.
- Switch Contexts: Switch to a specific context by name.
- Show Current Context: Display the currently active context.
- Autocompletion: Supports shell autocompletion for commands and context names (Bash, Zsh, Fish).
- Cross-Platform: Works on Linux, macOS, and optionally Windows.

---

## Installation

### 1. Quick Installation via "curl-to-bash"
You can install `kubex` with a single command using the following "curl-to-bash" script:
```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/kubex/main/install.sh | bash
```
> **Note**: Always inspect the script before running it:
> ```bash
> curl -fsSL https://raw.githubusercontent.com/yourusername/kubex/main/install.sh | less
> ```

The script will automatically detect your operating system and architecture, download the appropriate binary, and install it to `/usr/local/bin`.

---

### 2. Download Precompiled Binaries
Download the appropriate binary for your operating system:

- Linux:
  ```bash
  wget https://example.com/kubex-linux -O kubex
  chmod +x kubex
  sudo mv kubex /usr/local/bin/
  ```

- macOS:
  ```bash
  curl -L https://example.com/kubex-macos -o kubex
  chmod +x kubex
  sudo mv kubex /usr/local/bin/
  ```

> Replace `https://example.com/` with the actual URL where you host the binaries.

---

### 3. Build from Source
If you prefer to build `kubex` from source, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/kubex.git
   cd kubex
   ```

2. Build the binary:
   ```bash
   go build -o kubex
   ```

3. Move the binary to a directory in your `PATH`:
   ```bash
   sudo mv kubex /usr/local/bin/
   ```

---

## Usage

### List Contexts
View all available Kubernetes contexts:
```bash
kubex list
```

Example output:
```
Available contexts:
- minikube
- gke_my-cluster (active)
- docker-desktop
```

---

### Switch Context
Switch to a specific context:
```bash
kubex use <context-name>
```

Example:
```bash
kubex use minikube
```

Output:
```
Switched to context: minikube
```

---

### Show Current Context
Display the currently active context:
```bash
kubex current
```

Example output:
```
Current context: minikube
```

---

## Autocompletion

`kubex` supports shell autocompletion for Bash, Zsh, and Fish. Follow the instructions below to enable it.

### Bash
Add the following line to your `~/.bashrc` or `~/.bash_profile`:
```bash
source <(kubex completion bash)
```

Reload the shell:
```bash
source ~/.bashrc
```

---

### Zsh
Add the following line to your `~/.zshrc`:
```bash
source <(kubex completion zsh)
```

Reload the shell:
```bash
source ~/.zshrc
```

---

### Fish
Add the following line to your Fish configuration file (`~/.config/fish/config.fish`):
```bash
kubex completion fish | source
```

---

## Building from Source

To build `kubex` for multiple platforms, use the following commands:

### Build for Linux
```bash
GOOS=linux GOARCH=amd64 go build -o kubex-linux
```

### Build for macOS
```bash
GOOS=darwin GOARCH=amd64 go build -o kubex-macos
```

### Automated Build Script
Save the following script as `build.sh` and run it to generate binaries for both platforms:
```bash
#!/bin/bash

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o kubex-linux

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o kubex-macos

echo "Binaries built: kubex-linux, kubex-macos"
```

Make it executable and run:
```bash
chmod +x build.sh
./build.sh
```

---

## Contributing

We welcome contributions! If you find a bug or have an idea for a new feature, please open an issue or submit a pull request.

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m "Add some feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

---

## License

This project is licensed under the [GNU General Public License](LICENSE).

---

## Acknowledgments

- Inspired by [`kubectx`](https://github.com/ahmetb/kubectx).
- Built using the Go programming language and the [`cobra`](https://github.com/spf13/cobra) library.

---

## Support

If you encounter any issues or have questions, feel free to open an issue on GitHub or reach out to the maintainers.

Happy Kubing! 🚀

