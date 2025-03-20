# Dtype Typing Test

A simple command-line typing test written in Go. This tool challenges you to type words accurately and quickly while providing real-time feedback with ANSI colored output.

## Features

- **Real-Time Feedback:** Colored text indicates correct (green), incorrect (red/dark red), or pending (yellow) characters.
- **Countdown Timer:** Starts on the first keystroke and displays remaining time.
- **Performance Metrics:** Calculates words per minute (WPM) and the number of correctly typed words at the end of the test.

## Requirements

- Go 1.24.1 or later

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/dmars8047/dtype.git
   ```
2. Change into the project directory:
   ```bash
   cd dtype
   ```
3. Build the project:
   ```bash
   go build -o dtype
   ```

## Usage

To run a typing test with a custom duration (in seconds), use:

```bash
./dtype -t 60
```

Replace `60` with any test duration between 10 and 300 seconds.

## Acknowledgments

- [github.com/eiannone/keyboard](https://github.com/eiannone/keyboard) for handling keyboard inputs.

## License

This project is open source. Feel free to modify and distribute it.