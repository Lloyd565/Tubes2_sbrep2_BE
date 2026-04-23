package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tubes2/backend/internal/parser"
	"tubes2/backend/internal/scraper"
	"tubes2/backend/internal/traversal"
)

const logDir = "logs"


type TraverseRequest struct {
	URL       string `json:"url"`
	HTML      string `json:"html"`
	Algorithm string `json:"algorithm"`
	Selector  string `json:"selector"`
	Limit     int    `json:"limit"`
}

type NodeJSON struct {
	Tag        string            `json:"tag"`
	ID         string            `json:"id,omitempty"`
	Classes    []string          `json:"classes,omitempty"`
	Depth      int               `json:"depth"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Children   []*NodeJSON       `json:"children,omitempty"`
	Visited    bool              `json:"visited"`
	Matched    bool              `json:"matched"`
}

type MatchJSON struct {
	Tag        string            `json:"tag"`
	ID         string            `json:"id,omitempty"`
	Classes    []string          `json:"classes,omitempty"`
	Depth      int               `json:"depth"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type StatsJSON struct {
	VisitedCount int     `json:"visited_count"`
	MatchCount   int     `json:"match_count"`
	MaxDepth     int     `json:"max_depth"`
	DurationMs   float64 `json:"duration_ms"`
}

type TraverseResponse struct {
	Tree         *NodeJSON   `json:"tree"`
	Matches      []MatchJSON `json:"matches"`
	TraversalLog []string    `json:"traversal_log"`
	Stats        StatsJSON   `json:"stats"`
	LogFile      string      `json:"log_file"`
}

func TraverseHandler(w http.ResponseWriter, r *http.Request) {
	var req TraverseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var htmlStr string
	if req.URL != "" {
		var err error
		htmlStr, err = scraper.Scrape(req.URL)
		if err != nil {
			http.Error(w, "scrape error: "+err.Error(), http.StatusBadGateway)
			return
		}
	} else if req.HTML != "" {
		htmlStr = req.HTML
	} else {
		http.Error(w, "url or html is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Selector) == "" {
		http.Error(w, "selector is required", http.StatusBadRequest)
		return
	}

	alg := strings.ToLower(strings.TrimSpace(req.Algorithm))
	if alg != "bfs" && alg != "dfs" {
		http.Error(w, "algorithm must be 'bfs' or 'dfs'", http.StatusBadRequest)
		return
	}

	root, err := parser.Parse(htmlStr)
	if err != nil {
		http.Error(w, "parse error: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tokens := parser.CombinatorParser(req.Selector)

	start := time.Now()
	var matches, visited []*parser.Node
	var log []string
	if alg == "bfs" {
		matches, visited, log = traversal.BFS(req.Limit, root, tokens)
	} else {
		matches, visited, log = traversal.DFS(req.Limit, root, tokens)
	}
	duration := time.Since(start)

	visitedSet := make(map[*parser.Node]bool, len(visited))
	for _, n := range visited {
		visitedSet[n] = true
	}
	matchedSet := make(map[*parser.Node]bool, len(matches))
	for _, n := range matches {
		matchedSet[n] = true
	}

	maxDepth := 0
	tree := serializeNode(root, visitedSet, matchedSet, &maxDepth)

	matchesJSON := make([]MatchJSON, len(matches))
	for i, n := range matches {
		matchesJSON[i] = MatchJSON{
			Tag:        n.Tag,
			ID:         n.ID,
			Classes:    n.Classes,
			Depth:      n.Depth,
			Attributes: n.Attributes,
		}
	}

	logFile := saveLog(alg, req.Selector, log, duration, len(visited), len(matches), maxDepth)

	resp := TraverseResponse{
		Tree:         tree,
		Matches:      matchesJSON,
		TraversalLog: log,
		LogFile:      logFile,
		Stats: StatsJSON{
			VisitedCount: len(visited),
			MatchCount:   len(matches),
			MaxDepth:     maxDepth,
			DurationMs:   float64(duration.Nanoseconds()) / 1e6,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func LogHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "file query parameter required", http.StatusBadRequest)
		return
	}
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "..") {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	path := filepath.Join(logDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "log not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Write(data)
}

// helper
func saveLog(alg, selector string, lines []string, duration time.Duration, visitedCount, matchCount, maxDepth int) string {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return ""
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("traversal_%s.log", timestamp)
	path := filepath.Join(logDir, filename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Algorithm : %s\n", strings.ToUpper(alg)))
	sb.WriteString(fmt.Sprintf("Selector  : %s\n", selector))
	sb.WriteString(fmt.Sprintf("Timestamp : %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Duration  : %.3f ms\n", float64(duration.Nanoseconds())/1e6))
	sb.WriteString(fmt.Sprintf("Visited   : %d nodes\n", visitedCount))
	sb.WriteString(fmt.Sprintf("Matched   : %d nodes\n", matchCount))
	sb.WriteString(fmt.Sprintf("Max Depth : %d\n", maxDepth))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	for _, line := range lines {
		sb.WriteString(line + "\n")
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return ""
	}
	return filename
}

func serializeNode(n *parser.Node, visited, matched map[*parser.Node]bool, maxDepth *int) *NodeJSON {
	if n == nil {
		return nil
	}
	if n.Depth > *maxDepth {
		*maxDepth = n.Depth
	}
	node := &NodeJSON{
		Tag:        n.Tag,
		ID:         n.ID,
		Classes:    n.Classes,
		Depth:      n.Depth,
		Attributes: n.Attributes,
		Visited:    visited[n],
		Matched:    matched[n],
	}
	for _, child := range n.Children {
		node.Children = append(node.Children, serializeNode(child, visited, matched, maxDepth))
	}
	return node
}