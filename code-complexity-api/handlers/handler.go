package handlers

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type AnalyzeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type AnalyzeResponse struct {
	TimeComplexity  string `json:"timeComplexity"`
	SpaceComplexity string `json:"spaceComplexity"`
}

// WithCORS is a wrapper to add CORS headers
func WithCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

// AnalyzeHandler handles the request to analyze the code
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Received code:", req.Code)
	fmt.Println("Language:", req.Language)
	cleanCode := removeComments(req.Code)
	
	// res := analyzeCode(cleanCode, req.Language) //earlier call was direclty using main thread.
	timeChan := make(chan string)
	spaceChan := make(chan string)

	// Start goroutines
	go func() {
		timeChan <- analyzeTimeComplexity(cleanCode, req.Language)
	}()

	go func() {
		spaceChan <- analyzeSpaceComplexity(cleanCode, req.Language)
	}()

	// Code will wait till we get result from it.
	timeResult := <-timeChan
	spaceResult := <-spaceChan

	// Build response
	res := AnalyzeResponse{
		TimeComplexity:  timeResult,
		SpaceComplexity: spaceResult,
	}
	// Respond with time and space complexity
	json.NewEncoder(w).Encode(res)
}

// removeComments removes comments from code
func removeComments(code string) string {
	// Remove multiline comments
	multiLine := regexp.MustCompile(`(?s)/\*.*?\*/`)
	code = multiLine.ReplaceAllString(code, "")

	// Remove single line comments
	singleLine := regexp.MustCompile(`(?m)//.*$`)
	code = singleLine.ReplaceAllString(code, "")

	return code
}

// analyzeCode analyzes the code based on the language
func analyzeSpaceComplexity(code string, language string) string {
	return "Space complexity is currently not available."
	switch language {
	case "cpp", "java", "csharp":
		return analyzeCStyleCode(code)
	case "python":
		return analyzePythonCode(code)
	case "javascript", "golang":
		return analyzeJsOrGoCode(code)
	default:
		return "Unknown"
	}
}

func analyzeTimeComplexity(code string, language string) string {
	switch language {
	case "cpp", "java", "csharp":
		return analyzeCStyleCode(code)
	case "python":
		return analyzePythonCode(code)
	case "javascript", "golang":
		return analyzeJsOrGoCode(code)
	default:
		return "Unknown"
	}
}

func analyzeCStyleCode(code string) string {
	loopRegex := regexp.MustCompile(`(?m)for\s*\(|while\s*\(`)
	loopCount := len(loopRegex.FindAllString(code, -1))

	functionNameRegex := regexp.MustCompile(`(?m)(\w+)\s*\(.*\)\s*\{`)
	functionMatch := functionNameRegex.FindStringSubmatch(code)
	recursive := false
	if len(functionMatch) > 1 {
		funcName := functionMatch[1]
		if strings.Count(code, funcName+"(") > 1 {
			recursive = true
		}
	}

	time := "O(1)"
	if recursive {
		time = "O(n)"
	} else if loopCount == 1 {
		time = "O(n)"
	} else if loopCount > 1 {
		time = fmt.Sprintf("O(n^%d)", loopCount)
	}

	return time
}

func analyzePythonCode(code string) string {
	loopRegex := regexp.MustCompile(`(?m)for\s+\w+\s+in|while\s+`)
	loopCount := len(loopRegex.FindAllString(code, -1))
	recursive := strings.Contains(code, "def") && strings.Contains(code, "():")

	time := "O(1)"
	if recursive {
		time = "O(n)"
	} else if loopCount == 1 {
		time = "O(n)"
	} else if loopCount > 1 {
		time = fmt.Sprintf("O(n^%d)", loopCount)
	}

	return time
}

func analyzeJsOrGoCode(code string) string {
	loopRegex := regexp.MustCompile(`(?m)for\s*\(|while\s*\(`)
	loopCount := len(loopRegex.FindAllString(code, -1))

	time := "O(1)"
	if loopCount == 1 {
		time = "O(n)"
	} else if loopCount > 1 {
		time = fmt.Sprintf("O(n^%d)", loopCount)
	}

	return time
}
