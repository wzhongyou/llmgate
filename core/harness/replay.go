package harness

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/wzhongyou/llmgate/core"
)

// ReplayResult holds the comparison between original and replayed response.
type ReplayResult struct {
	Original Record             `json:"original"`
	Replayed *core.ChatResponse `json:"replayed,omitempty"`
	Error    string             `json:"error,omitempty"`
	Latency  float64            `json:"latency_ms"`
	Match    ReplayMatch        `json:"match"`
}

// ReplayMatch summarizes what matched/changed between original and replay.
type ReplayMatch struct {
	FinishReasonMatch bool    `json:"finish_reason_match"`
	ToolCallsMatch    bool    `json:"tool_calls_match"`
	TokenDeltaPct     float64 `json:"token_delta_pct"` // (new-old)/old * 100
}

// LoadRecords reads a JSONL file and returns up to limit records.
func LoadRecords(path string, limit int) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []Record
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		var rec Record
		if json.Unmarshal(scanner.Bytes(), &rec) == nil {
			records = append(records, rec)
			if limit > 0 && len(records) >= limit {
				break
			}
		}
	}
	return records, scanner.Err()
}

// Replay sends recorded requests to a target provider and returns comparison results.
func Replay(ctx context.Context, engine *core.Engine, provider string, records []Record) []ReplayResult {
	results := make([]ReplayResult, 0, len(records))
	for _, rec := range records {
		if rec.Request == nil {
			continue
		}
		start := time.Now()
		resp, err := engine.ChatWithProvider(ctx, rec.Request, provider)
		latency := float64(time.Since(start).Microseconds()) / 1000.0

		result := ReplayResult{
			Original: rec,
			Replayed: resp,
			Latency:  latency,
		}
		if err != nil {
			result.Error = err.Error()
		}
		if resp != nil && rec.Response != nil {
			result.Match = compareResponses(rec.Response, resp)
		}
		results = append(results, result)
	}
	return results
}

func compareResponses(orig, replay *core.ChatResponse) ReplayMatch {
	m := ReplayMatch{
		FinishReasonMatch: orig.FinishReason == replay.FinishReason,
		ToolCallsMatch:    len(orig.ToolCalls) == len(replay.ToolCalls),
	}
	if orig.Usage.OutputTokens > 0 {
		m.TokenDeltaPct = float64(replay.Usage.OutputTokens-orig.Usage.OutputTokens) / float64(orig.Usage.OutputTokens) * 100
	}
	return m
}
