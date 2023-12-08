package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: scan <domain>")
        os.Exit(1)
    }

    domain := os.Args[1]
    outputDir := filepath.Join(".", domain)

    // Create the output directory
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        fmt.Println("Error creating output directory:", err)
        os.Exit(1)
    }

    // Subfinder
    subfinderCmd := exec.Command("subfinder", "-d", domain, "-silent")
    subfinderCmd.Dir = outputDir
    subfinderOutput, _ := subfinderCmd.Output()
    writeToFile(filepath.Join(outputDir, "Subfinder.txt"), subfinderOutput)

    // Assetfinder
    assetfinderCmd := exec.Command("assetfinder", "-subs-only", domain)
    assetfinderCmd.Dir = outputDir
    assetfinderOutput, _ := assetfinderCmd.Output()
    writeToFile(filepath.Join(outputDir, "Assetfinder.txt"), assetfinderOutput)

    // Amass
    amassCmd := exec.Command("amass", "enum", "-passive", "-d", domain)
    amassCmd.Dir = outputDir
    amassOutput, _ := amassCmd.Output()
    writeToFile(filepath.Join(outputDir, "amass.txt"), amassOutput)

    // Combine and deduplicate results
    combineAndDeduplicate(outputDir)

    // Httpx
    httpxCmd := exec.Command("httpx", "-silent")
    httpxCmd.Dir = outputDir
    httpxOutput, _ := httpxCmd.Output()
    writeToFile(filepath.Join(outputDir, "alive.txt"), httpxOutput)

    // Nuclei
    nucleiCmd := exec.Command("nuclei", "-es", "info,unknown", "-etags", "ssl,network")
    nucleiCmd.Dir = outputDir
    nucleiOutput, _ := nucleiCmd.Output()
    writeToFile(filepath.Join(outputDir, "nuclei.txt"), nucleiOutput)
}

func writeToFile(filename string, data []byte) {
    if err := os.WriteFile(filename, data, 0644); err != nil {
        fmt.Println("Error writing to file:", err)
        os.Exit(1)
    }
}

func combineAndDeduplicate(outputDir string) {
    files, err := filepath.Glob(filepath.Join(outputDir, "*.txt"))
    if err != nil {
        fmt.Println("Error listing files:", err)
        os.Exit(1)
    }

    var combinedData []byte
    for _, file := range files {
        data, err := os.ReadFile(file)
        if err != nil {
            fmt.Println("Error reading file:", err)
            os.Exit(1)
        }
        combinedData = append(combinedData, data...)
    }

    // Remove duplicates
    uniqueLines := removeDuplicates(strings.Split(string(combinedData), "\n"))

    // Write to a new file
    writeToFile(filepath.Join(outputDir, "combined.txt"), []byte(strings.Join(uniqueLines, "\n")))
}

func removeDuplicates(lines []string) []string {
    uniqueLines := make(map[string]struct{})
    for _, line := range lines {
        if line != "" {
            uniqueLines[line] = struct{}{}
        }
    }

    var result []string
    for line := range uniqueLines {
        result = append(result, line)
    }

    return result
}
