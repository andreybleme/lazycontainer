# lazycontainer

A terminal UI to manage Apple Containers without stress. Written in Go with [Bubbletea](https://github.com/charmbracelet/bubbletea) 🧋

![lazycontainerdemo](https://github.com/user-attachments/assets/71220800-46e3-4932-a0c0-9e4fe55ff99b)

## Requirements

- [Apple containers](https://github.com/apple/container) CLI **0.1.0**

## Install

### Homebrew
```
brew tap andreybleme/lazycontainer
brew install lazycontainer
```

## Usage

Start the terminal UI:

```
$ lazycontainer
```

Press `key-up` ⬆️ / `key-down` ⬇️ to navigate across containers.

Press `tab` to switch between containers and images.

Press `enter` to select a resource (container or image) and see its details.

Press `q` or `ctrl+c` to exit

## Features

This is an alpha release, so you may find bugs and missing features. Currently, these are the supported features:

- viewing the state of containers
- inspecting the details of a container
- vieweing the state of images
- inspecting the details of an image

## Running 

1. **Clone the repository:**

   ```bash
   git clone https://github.com/andreybleme/lazycontainer
   cd lazycontainer
   ```

2. **Install dependencies:**

   Run the following command to install the necessary dependencies:

   ```bash
   go mod tidy
   ```

3. **Run the application:**

   You can run the application using the following command:

   ```bash
   go run cmd/main.go
   ```

## Contributing

Contributions are welcome! Feel free to submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
