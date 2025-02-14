SSE Server with weak.Pointer

## Overview

This project demonstrates an efficient way to manage Server-Sent Events (SSE) connections in Go using weak.Pointer. The approach ensures that inactive clients are automatically garbage collected, reducing memory leaks and optimizing resource usage.

## Features

- Efficient SSE connection management using weak.Pointer

- Automatic cleanup of inactive connections

- Periodic message broadcasting to active clients

- CORS support for cross-origin communication

- Minimal memory overhead by allowing GC to handle expired references

## Installation

Clone the repository:

```bash
git clone https://github.com/9ssi7/goweaksse.git
cd goweaksse
```

Run the server:

```bash
go run main.go
```

## License

This project is licensed under the MIT License.
