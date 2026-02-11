# Detect Build System

Your task is to identify the build tools and dependency management systems used in this repository. Use the file system tools to scan for build configuration files.

Look for:
- go.mod, go.sum (Go)
- package.json, package-lock.json (Node.js)
- Cargo.toml, Cargo.lock (Rust)
- pom.xml, build.gradle (Java)
- requirements.txt, pyproject.toml (Python)
- Makefile, CMakeLists.txt (C/C++)

Store the detected build system information for use in later phases.
