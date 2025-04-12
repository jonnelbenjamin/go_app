package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"  // For colored terminal output
	"github.com/spf13/cobra"  // CLI framework
)

// Global variables for command-line flags
var (
	uppercase   bool    // --uppercase flag
	name        string  // --name flag
	showWeather bool    // --weather flag
	showJoke    bool    // --joke flag
	showFact    bool    // --fact flag
	guessNumber bool    // --game flag
	city        string  // --city flag for weather location
)

// Weather struct to parse API response from weatherapi.com
type Weather struct {
	Location struct {
		Name string `json:"name"`  // City name
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`  // Temperature in Celsius
		Condition struct {
			Text string `json:"text"`  // Weather description
		} `json:"condition"`
	} `json:"current"`
}

// rootCmd defines the main CLI command and its behavior
var rootCmd = &cobra.Command{
	Use:   "greeter",  // Command name
	Short: "A supercharged CLI app with free API features",  // Short description
	Long: `A colorful CLI application that can:  // Detailed help text
- Greet you personally
- Tell you the weather (no API key needed)
- Share interesting facts
- Tell jokes
- Play number guessing games`,
	// The main execution function
	Run: func(cmd *cobra.Command, args []string) {
		// If no name provided via flag, prompt user
		if name == "" {
			color.Yellow("What's your name? ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
		}

		// Initialize color printers
		green := color.New(color.FgGreen).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		// Print personalized greeting
		fmt.Printf("%s Hello, %s!\n", green(">>"), name)

		// Handle uppercase flag
		if uppercase {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s YOUR NAME IN UPPERCASE: %s\n", red("!!"), strings.ToUpper(name))
		}

		// Execute features based on flags
		if showWeather {
			getWeather()
		}
		if showJoke {
			tellJoke()
		}
		if showFact {
			tellFact()
		}
		if guessNumber {
			playNumberGame()
		}

		// Show help if no features were selected
		if !uppercase && !showWeather && !showJoke && !showFact && !guessNumber {
			cyan("\nTry these commands:")
			cyan("  --weather\tGet current weather (use with --city)")
			cyan("  --joke\tHear a random joke")
			cyan("  --fact\tLearn an interesting fact")
			cyan("  --game\tPlay a number guessing game")
		}
	},
}

// init defines all CLI flags and their configurations
func init() {
	// Personalization flags
	rootCmd.Flags().BoolVarP(&uppercase, "uppercase", "u", false, "Shout your name")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Provide name directly")
	
	// Feature flags
	rootCmd.Flags().BoolVar(&showWeather, "weather", false, "Show current weather")
	rootCmd.Flags().StringVar(&city, "city", "London", "City for weather forecast")
	rootCmd.Flags().BoolVar(&showJoke, "joke", false, "Tell a random joke")
	rootCmd.Flags().BoolVar(&showFact, "fact", false, "Share an interesting fact")
	rootCmd.Flags().BoolVar(&guessNumber, "game", false, "Play number guessing game")
}

// Execute runs the root command and handles errors
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

// ----------------------------
// Feature Implementations
// ----------------------------

// getWeather fetches and displays current weather using weatherapi.com
func getWeather() {
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	// Build API URL with city parameter
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?", city)
	
	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		color.Red("Failed to get weather: %v", err)
		return
	}
	defer resp.Body.Close()  // Ensure response body is closed

	// Parse JSON response
	var weather Weather
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		color.Red("Failed to decode weather: %v", err)
		return
	}

	// Display formatted weather information
	fmt.Printf("%s Current weather in %s: %s%.1f°C (%s)%s\n", 
		yellow("☀️"), 
		blue(weather.Location.Name),
		blue(""),
		weather.Current.TempC,
		weather.Current.Condition.Text,
		blue(""))
}

// tellJoke fetches a random joke from an API or uses local fallback
func tellJoke() {
	// Define structure to parse API response
	type Joke struct {
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	}

	// Try to get joke from API
	resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		color.Red("Failed to get joke: %v", err)
		return
	}
	defer resp.Body.Close()

	var joke Joke
	if err := json.NewDecoder(resp.Body).Decode(&joke); err != nil {
		// Fallback to local jokes if API fails
		localJokes := []string{
			"Why don't scientists trust atoms? Because they make up everything!",
			"Parallel lines have so much in common... it's a shame they'll never meet.",
		}
		rand.Seed(time.Now().UnixNano())
		color.Cyan("\n🎭 Joke: %s", localJokes[rand.Intn(len(localJokes))])
		return
	}

	// Display joke with setup and punchline
	color.Cyan("\n🎭 %s", joke.Setup)
	color.Cyan("   %s", joke.Punchline)
}

// tellFact fetches a random fact from an API or uses local fallback
func tellFact() {
	// Try to get fact from API
	resp, err := http.Get("https://uselessfacts.jsph.pl/api/v2/facts/random?language=en")
	if err != nil {
		color.Red("Failed to get fact: %v", err)
		return
	}
	defer resp.Body.Close()

	// Parse API response
	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Fallback to local facts if API fails
		localFacts := []string{
			"The Eiffel Tower can be 15 cm taller during summer due to thermal expansion.",
			"Bananas are berries, but strawberries aren't.",
		}
		rand.Seed(time.Now().UnixNano())
		color.Magenta("\n🧠 Fact: %s", localFacts[rand.Intn(len(localFacts))])
		return
	}

	// Display the fact
	color.Magenta("\n🧠 Fact: %s", result.Text)
}

// playNumberGame implements an interactive number guessing game
func playNumberGame() {
	// Initialize game with random target number (1-100)
	rand.Seed(time.Now().UnixNano())
	target := rand.Intn(100) + 1
	attempts := 0

	color.Blue("\n🔢 I'm thinking of a number between 1 and 100. Can you guess it?")

	reader := bufio.NewReader(os.Stdin)
	for {
		attempts++
		color.Yellow("Your guess: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Convert input to number
		guess, err := strconv.Atoi(input)
		if err != nil {
			color.Red("Please enter a valid number!")
			continue
		}

		// Check guess against target
		switch {
		case guess < target:
			color.Cyan("Too low!")
		case guess > target:
			color.Cyan("Too high!")
		default:
			color.Green("🎉 Correct! You guessed it in %d attempts!", attempts)
			return
		}
	}
}