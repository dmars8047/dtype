package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	defaultTestDuration  = 60
	startingLineCount    = 3
	maxCharactersPerLine = 60
	startingRow          = 17
	startingCol          = 0
	welcomeMessage       = `       ________               
  ____/ /_  __/_  ______  ___ 
 / __  / / / / / / / __ \/ _ \
/ /_/ / / / / /_/ / /_/ /  __/
\__,_/ /_/  \__, / .___/\___/ 
	   /____/_/           
	 
Welcome to the Typing Test!

- You'll have {testTime} seconds to type as many words as you can. 
- The test timer will start as soon as you begin typing.
- Press any key to start the test.
- Press ESC to exit the test at any time.

Time remaining: {remainingTime} seconds

`
)

type CharacterStatus = uint8

const (
	NotSet = iota
	Correct
	Incorrect
)

type Prompt struct {
	Lines []Line
	// The index of the current line being typed
	LineIndex int
}

func NewPrompt() *Prompt {
	prompt := Prompt{
		Lines: []Line{},
	}

	// Generate the exact number of lines specified by promptLineCount
	for range startingLineCount {
		addNewLine(&prompt)
	}

	return &prompt
}

// Dynamically generate a new line.
func addNewLine(prompt *Prompt) {
	currentLine := Line{
		Words: []Word{},
	}

	var numCharactersInLine int

	prevIdx := -1

	// Keep adding words to the current line until we approach maxCharactersPerLine
	for {
		idx := -2

		for idx == -2 || idx == prevIdx {
			idx = rand.Intn(len(wordList))
		}

		prevIdx = idx

		// Get a random word from the wordList
		randomWord := MakeWord(wordList[idx])

		// Calculate how many characters this word would add (including space)
		numCharactersInLine += len(randomWord.Value)

		if len(currentLine.Words) > 0 {
			// Add a space for the word separator
			numCharactersInLine++
		}

		// Check if adding this word would exceed the max characters per line
		if numCharactersInLine > maxCharactersPerLine {
			// Line is full enough, stop adding words
			break
		}

		// Add the word to the current line
		currentLine.Words = append(currentLine.Words, randomWord)
	}

	// Add the completed line to the prompt
	prompt.Lines = append(prompt.Lines, currentLine)
}

// diffAndPrint compares the oldLines and newLines and updates only changed lines.
func diffAndPrint(oldLines, newLines []string) {
	maxLines := max(len(oldLines), len(newLines))

	// Hide the cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show the cursor again when done

	for i := range maxLines {
		oldLine := ""
		newLine := ""

		if i < len(oldLines) {
			oldLine = oldLines[i]
		}

		if i < len(newLines) {
			newLine = newLines[i]
		}

		if oldLine != newLine {
			// Move the cursor to the line (assuming terminal rows start at 1)
			fmt.Printf("\033[%d;1H", i+1)
			// Print the new line and clear the remainder of the line
			fmt.Printf("%s\033[K", newLine)
		}
	}
}

// Renamed RenderScreenBuffer replaces the old Print method to build and return the rendered screen buffer.
func (p *Prompt) RenderScreenBuffer(testTime, remainingSeconds int) []string {
	var sb strings.Builder

	welcome := strings.Replace(welcomeMessage, "{testTime}", fmt.Sprintf("%d", testTime), 1)
	welcome = strings.Replace(welcome, "{remainingTime}", fmt.Sprintf("%d", remainingSeconds), 1)

	sb.WriteString(welcome)

	for lineIdx, line := range p.Lines {
		for _, word := range line.Words {
			// Process each character in the word
			for charIdx, r := range word.Value {
				if charIdx < len(word.CharStatuses) {
					switch word.CharStatuses[charIdx] {
					case Correct:
						// Correctly typed characters are green
						sb.WriteString(colorGreen)
						sb.WriteRune(r)
						sb.WriteString(colorReset)
					case Incorrect:
						// If the expected character is a space, use dark red
						if r == ' ' {
							sb.WriteString(colorDarkRed)
							sb.WriteString("█")
						} else {
							sb.WriteString(colorRed)
							sb.WriteRune(r)
						}
						sb.WriteString(colorReset)
					case NotSet:
						// Characters that haven't been typed yet are yellow
						sb.WriteString(colorYellow)
						sb.WriteRune(r)
						sb.WriteString(colorReset)
					}
				} else {
					// Default for characters beyond current TypedStatus is yellow
					sb.WriteString(colorYellow)
					sb.WriteRune(r)
					sb.WriteString(colorReset)
				}
			}
		}

		// Add a newline after each line except the last one
		if lineIdx < len(p.Lines)-1 {
			sb.WriteString("\n")
		}
	}

	return strings.Split(sb.String(), "\n")
}

type Word struct {
	// The value of the word
	Value string
	// A slice of booleans indicating the correctness of each character typed
	CharStatuses []CharacterStatus
	// The index of the current character being typed
	CharIndex int
}

func MakeWord(value string) Word {
	return Word{
		Value:        fmt.Sprintf("%s ", value),
		CharStatuses: make([]CharacterStatus, len(value)+1),
	}
}

// Line represents a line of text in the typing test.
type Line struct {
	// The index of the line in the prompt
	Words []Word

	// The index of the current word being typed
	WordIndex int
}

// Available words to use in the typing test
var wordList = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
	"person", "into", "year", "your", "good", "some", "could", "them", "see", "other",
	"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
	"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
	"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	"those", "such", "through", "between", "both", "few", "while", "might", "place",
	"long", "need", "same", "right", "still", "own", "last", "never", "under",
	"double", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven",
	"triple", "baseball", "video", "computer", "global", "save", "widget", "cell", "stream", "input",
	"output", "love", "peace", "world", "hello", "tangle", "hero", "villain",
	"goodbye", "goodnight", "morning", "night", "week", "month", "super",
	"wifi", "game", "click", "groove", "raspberry", "put", "post", "delete", "tiny", "core",
	"pasta", "pizza", "taco", "cheese", "water", "honey", "iron", "diamond", "gold", "silver",
	"dungeon", "dragon", "sword", "shield", "armor", "helmet", "boots", "gloves", "ring", "amulet", "castle",
	"eagle", "whale", "sushi", "ghost", "zombie", "rush", "crash", "slice", "bakery", "coffee",
	"cocktail", "beer", "bear", "lion", "tiger", "attack", "decay", "envelope", "guitar", "piano", "teeth", "bite",
	"scratch", "stomp", "raptor", "trigger", "fire",
}

type printRequest struct {
	oldLines []string
	newLines []string
}

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[91m"
	colorDarkRed = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
)

func main() {
	testDurationSeconds := defaultTestDuration

	// Get the -t flag value if it is provided and valid
	flag.IntVar(&testDurationSeconds, "t", defaultTestDuration, "The duration of the typing test in seconds. Must be between 1 and 300.")

	flag.Parse()

	if testDurationSeconds < 10 {
		fmt.Println("The test duration must be at least 10 second.")
		return
	}

	if testDurationSeconds > 300 {
		fmt.Println("The test duration must be less than 300 seconds.")
		return
	}

	clearScreen()

	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	// Ensure the cursor is shown again when the program exits
	defer fmt.Print("\033[?25h")

	// Present the prompt (unmodified text)
	prompt := NewPrompt()

	// Create a channel to signal termination.
	endChan := make(chan struct{})
	printChan := make(chan printRequest, 1)

	var wg sync.WaitGroup
	// Declare timerOnce to ensure timer starts only on first key stroke.
	var timerOnce sync.Once

	// Initially show the prompt with the cursor at the start of the first word.
	oldScreen := prompt.RenderScreenBuffer(testDurationSeconds, testDurationSeconds)

	renderMutex := sync.Mutex{}
	remainingSeconds := testDurationSeconds
	escaped := false

	for _, line := range oldScreen {
		fmt.Println(line)
	}

	wg.Add(1)
	go func() {
		defer keyboard.Close()
		defer wg.Done()
		defer close(endChan)

		for {
			renderMutex.Lock()
			if remainingSeconds <= 0 {
				renderMutex.Unlock()
				return
			}
			renderMutex.Unlock()

			// Get the current line and word being typed
			currentLine := &prompt.Lines[prompt.LineIndex]
			currentWord := &currentLine.Words[currentLine.WordIndex]

			typedChar, key, err := keyboard.GetSingleKey()

			if err != nil {
				return
			}

			renderMutex.Lock()

			if remainingSeconds <= 0 {
				renderMutex.Unlock()
				return
			}

			// Lazy-start timer on first keystroke.
			timerOnce.Do(func() {
				wg.Add(1)
				go func() {
					defer wg.Done()
					timer := time.NewTimer(time.Duration(testDurationSeconds) * time.Second)
					select {
					case <-endChan:
						return
					case <-timer.C:
						keyboard.Close()
						return
					}
				}()

				wg.Add(1)
				// Use ANSI codes to save and restore the cursor position.
				go func() {
					defer wg.Done()
					ticker := time.NewTicker(1 * time.Second)
					defer ticker.Stop()

					for {
						if remainingSeconds <= 0 {
							return
						}
						select {
						case <-endChan:
							return
						case <-ticker.C:
							renderMutex.Lock()
							remainingSeconds--
							newScreen := prompt.RenderScreenBuffer(testDurationSeconds, remainingSeconds)
							printChan <- printRequest{
								oldLines: oldScreen,
								newLines: newScreen,
							}
							oldScreen = newScreen
							renderMutex.Unlock()
						}
					}
				}()
			})

			// Process backspace and regular typing
			if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
				lastLine := len(oldScreen)
				fmt.Printf("\033[%d;1H", lastLine+2)
				fmt.Printf("Exiting...\n\n")
				escaped = true
				renderMutex.Unlock()
				return
			}

			if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
				// If they are in the middle of the word
				if currentWord.CharIndex > 0 {
					// Move the cursor back one character
					currentWord.CharIndex--
					currentWord.CharStatuses[currentWord.CharIndex] = NotSet
				} else if currentLine.WordIndex > 0 && currentWord.CharIndex == 0 {
					// Move to the previous word
					currentLine.WordIndex--
					currentWord = &currentLine.Words[currentLine.WordIndex]
					// Reset the status of the last character in the previous word
					currentWord.CharStatuses[len(currentWord.CharStatuses)-1] = NotSet
				} else {
					// Move to the previous line
					prompt.LineIndex--
					currentLine = &prompt.Lines[prompt.LineIndex]
					currentWord = &currentLine.Words[currentLine.WordIndex]
					currentWord.CharStatuses[len(currentWord.CharStatuses)-1] = NotSet
				}
			} else {
				// Compare the character with the expected character
				var status CharacterStatus
				expectedChar := currentWord.Value[currentWord.CharIndex]

				if typedChar == rune(expectedChar) || (key == keyboard.KeySpace && currentWord.Value[currentWord.CharIndex] == ' ') {
					status = Correct
				} else {
					status = Incorrect
				}

				// Set the status for the current character
				currentWord.CharStatuses[currentWord.CharIndex] = status

				// When the current word is completed...
				if (len(currentWord.Value) - 1) == currentWord.CharIndex {
					// If this is the last word in the current line...
					if currentLine.WordIndex == len(currentLine.Words)-1 {
						// And if we are at the last generated line, generate a new one.
						if prompt.LineIndex+3 >= len(prompt.Lines) {
							addNewLine(prompt)
						}
						prompt.LineIndex++
					} else {
						currentLine.WordIndex++
					}
				} else {
					currentWord.CharIndex++
				}
			}

			// Redraw the screen with updated state
			newScreen := prompt.RenderScreenBuffer(testDurationSeconds, remainingSeconds)

			// Send the old and new screen to the print channel
			printChan <- printRequest{
				oldLines: oldScreen,
				newLines: newScreen,
			}

			oldScreen = newScreen
			renderMutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-endChan:
				return
			case request := <-printChan:
				diffAndPrint(request.oldLines, request.newLines)
				// Update the cursor position
				row, col := computeCursorPosition(prompt, startingRow, startingCol)
				fmt.Printf("\033[%d;%dH", row, col)
			}
		}
	}()

	// Update the cursor position
	row, col := computeCursorPosition(prompt, startingRow, startingCol)
	fmt.Printf("\033[%d;%dH", row, col)

	wg.Wait()

	if escaped {
		return
	}

	// Print "Time's up!" below the prompt on a new row.
	lastLine := len(oldScreen)
	fmt.Printf("\033[%d;1H", lastLine+2)
	fmt.Printf("Time's up!\n\n")

	// Calculate the number of correct words typed
	correctWords := 0

	for _, line := range prompt.Lines {
		for _, word := range line.Words {
			correct := true
			for _, status := range word.CharStatuses {
				if status != Correct {
					correct = false
					break
				}
			}
			if correct {
				correctWords++
			}
		}
	}

	fmt.Printf("You typed %d words correctly.\n", correctWords)
	wpm := float64(correctWords) / (float64(testDurationSeconds) / 60.0)
	fmt.Printf("That's approximately %d words per minute!\n\n", int(math.Ceil(wpm)))
}

func computeCursorPosition(prompt *Prompt, startRow, startCol int) (int, int) {
	// Each line adds one to the row.
	row := startRow + prompt.LineIndex
	col := startCol

	currentLine := prompt.Lines[prompt.LineIndex]

	// For each complete word in the current line add the word length and a space.
	for i := 0; i < currentLine.WordIndex && i < len(currentLine.Words); i++ {
		col += len(currentLine.Words[i].Value) // add one for the space
	}

	// If the user is still in the middle of a word,
	// add the character offset within that word.
	if currentLine.WordIndex < len(currentLine.Words) {
		col += currentLine.Words[currentLine.WordIndex].CharIndex
	}

	return row, col + 1 // it should be on the next char to be typed
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
