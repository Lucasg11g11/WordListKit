package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	passwordMutex      sync.Mutex
	generatedPasswords = make(map[string]bool)
)

var theme = map[string]string{
	"Info":    "\033[36m",
	"Success": "\033[32m",
	"Error":   "\033[31m",
	"Reset":   "\033[0m",
}

// Get absolute paths for wordlist and dictionary
func getPaths() (string, string) {
	basePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wordlistFile := filepath.Join(basePath, "StormList.txt")
	dictionaryFile := filepath.Join(basePath, "dictionary.txt")
	return wordlistFile, dictionaryFile
}

// Load existing passwords to avoid duplicates
func loadExistingPasswords(wordlistFile string) {
	file, err := os.Open(wordlistFile)
	if err != nil {
		fmt.Println(theme["Info"], "[INFO] Wordlist not found, starting fresh.", theme["Reset"])
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if password != "" {
			generatedPasswords[password] = true
		}
	}
	fmt.Printf("%s[INFO] Loaded %d existing passwords from wordlist.%s\n", theme["Info"], len(generatedPasswords), theme["Reset"])
}

func loadDictionary(dictionaryFile string) ([]string, error) {
	file, err := os.Open(dictionaryFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	return words, nil
}

// Generate strong random password
func generateRandomPassword(minLen, maxLen int) string {
	length := rand.Intn(maxLen-minLen+1) + minLen
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}

// Generate password using a real word + 4-digit number
func generateRealPassword(words []string) string {
	word := words[rand.Intn(len(words))]
	suffix := fmt.Sprintf("%04d", rand.Intn(10000))
	return word + suffix
}

func loadNames(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names = append(names, name)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

// Gera uma senha usando um nome aleatório + número de 4 dígitos
func generateNamePassword(names []string) string {
	rand.Seed(time.Now().UnixNano())
	name := names[rand.Intn(len(names))]
	number := rand.Intn(10000) // 0 até 9999
	return fmt.Sprintf("%s%04d", name, number)
}

// Save password to file if it's new
func savePasswordToFile(wordlistFile, password string) {
	passwordMutex.Lock()
	defer passwordMutex.Unlock()

	if !generatedPasswords[password] {
		f, err := os.OpenFile(wordlistFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(theme["Error"], "[ERROR] Could not open wordlist:", err, theme["Reset"])
			return
		}
		defer f.Close()

		if _, err := f.WriteString(password + "\n"); err != nil {
			fmt.Println(theme["Error"], "[ERROR] Could not write password:", err, theme["Reset"])
			return
		}

		generatedPasswords[password] = true
		fmt.Println(theme["Success"], "[SUCCESS] Password saved:", password, theme["Reset"])
	}
}

func continuousPasswordGenerator(wordlistFile string, realWords, names []string, useReal, useNames bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var password string
		if useReal && len(realWords) > 0 {
			password = generateRealPassword(realWords)
		} else if useNames && len(names) > 0 {
			password = generateNamePassword(names)
		} else {
			password = generateRandomPassword(12, 18)
		}
		savePasswordToFile(wordlistFile, password)
		// Remove sleep pra acelerar
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	wordlistFile, dictionaryFile := getPaths()
	loadExistingPasswords(wordlistFile)

	realWords, err := loadDictionary(dictionaryFile)
	if err != nil {
		fmt.Println(theme["Erro"], "[ERROR] Failed to load dictionary words:", err, theme["Reset"])
		realWords = nil
	}

	names, err := loadNames("names.txt")
	if err != nil {
		fmt.Println(theme["Erro"], "[ERROR] Failed to load names:", err, theme["Reset"])
		names = nil
	}

	var mode string
	fmt.Println("Choose password generation mode:")
	fmt.Println("1 - random (strong random passwords)")
	fmt.Println("2 - real (passwords based on dictionary words)")
	fmt.Println("3 - names (passwords based on names)")
	fmt.Print("Enter 1, 2, or 3: ")
	fmt.Scanln(&mode)

	useReal := false
	useNames := false

	switch mode {
	case "1":
		fmt.Println(theme["Informação"], "[INFO] Random mode selected.", theme["Reset"])
	case "2":
		useReal = true
		fmt.Println(theme["Informação"], "[INFO] Dictionary word mode selected.", theme["Reset"])
	case "3":
		useNames = true
		fmt.Println(theme["Informação"], "[INFO] Names mode selected.", theme["Reset"])
	default:
		fmt.Println(theme["Erro"], "[ERROR] Invalid option, defaulting to random mode.", theme["Reset"])
	}

	const numThreads = 100
	var wg sync.WaitGroup

	fmt.Printf("%s[INFO] Starting password generation with %d goroutines...%s\n", theme["Informação"], numThreads, theme["Reset"])

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go continuousPasswordGenerator(wordlistFile, realWords, names, useReal, useNames, &wg)
	}

	wg.Wait()
}
