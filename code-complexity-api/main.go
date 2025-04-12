package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        handler(w, r)
    }
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
    var req AnalyzeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    fmt.Println("Received code:", req.Code)
    fmt.Println("Language:", req.Language)
	cleanCode := removeComments(req.Code)
    res := analyzeCode(cleanCode, req.Language)

    json.NewEncoder(w).Encode(res)
}

func removeComments(code string) string {
	// Remove multiline comments
	multiLine := regexp.MustCompile(`(?s)/\*.*?\*/`)
	code = multiLine.ReplaceAllString(code, "")

	// Remove single line comments
	singleLine := regexp.MustCompile(`(?m)//.*$`)
	code = singleLine.ReplaceAllString(code, "")

	return code
}


func analyzeCode(code string, language string) AnalyzeResponse {
    switch language {
    case "cpp", "java", "csharp":
        return analyzeCStyleCode(code)
    case "python":
        return analyzePythonCode(code)
    case "javascript", "golang":
        return analyzeJsOrGoCode(code)
    default:
        return AnalyzeResponse{TimeComplexity: "Unknown", SpaceComplexity: "Unknown"}
    }
}

func analyzeCStyleCode(code string) AnalyzeResponse {
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

	return AnalyzeResponse{
		TimeComplexity:  time,
		SpaceComplexity: "O(1)",
	}
}

func analyzePythonCode(code string) AnalyzeResponse {
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

	return AnalyzeResponse{
		TimeComplexity:  time,
		SpaceComplexity: "O(1)",
	}
}

func analyzeJsOrGoCode(code string) AnalyzeResponse {
	loopRegex := regexp.MustCompile(`(?m)for\s*\(|while\s*\(`)
	loopCount := len(loopRegex.FindAllString(code, -1))

	time := "O(1)"
	if loopCount == 1 {
		time = "O(n)"
	} else if loopCount > 1 {
		time = fmt.Sprintf("O(n^%d)", loopCount)
	}

	return AnalyzeResponse{
		TimeComplexity:  time,
		SpaceComplexity: "O(1)",
	}
}

func main() {
    http.HandleFunc("/analyze", withCORS(analyzeHandler))

    fmt.Println("🚀 Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
