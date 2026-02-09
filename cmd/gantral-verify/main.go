package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Rainminds/gantral/pkg/models"
	"github.com/Rainminds/gantral/pkg/verifier"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gantral-verify",
		Short: "Offline verification tool for Gantral Commitment Artifacts",
		Long: `A standalone tool for auditors to cryptographically verify 
the integrity and chain-of-custody of Gantral artifacts 
without requiring access to the operational database.`,
	}

	// Subcommand: verify file
	var fileCmd = &cobra.Command{
		Use:   "file [path]",
		Short: "Verify a single artifact file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("❌ ERROR: Failed to read file: %v\n", err)
				os.Exit(2)
			}

			result, err := verifier.VerifyArtifact(data)
			if err != nil {
				fmt.Printf("❌ ERROR: Logic failure: %v\n", err)
				os.Exit(2)
			}

			// JSON Output Support (Defaulting to plain text for now, but structured)
			if result.Valid {
				fmt.Printf("✅ VALID | ID: %s | Hash: %s\n", result.ArtifactID, result.CalculatedHash)
				os.Exit(0)
			} else {
				fmt.Printf("❌ INVALID | ID: %s | Error: %s\n", result.ArtifactID, result.Error)
				os.Exit(1)
			}
		},
	}

	// Subcommand: verify chain
	var chainCmd = &cobra.Command{
		Use:   "chain [directory]",
		Short: "Verify a chain of artifacts in a directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			files, err := os.ReadDir(dir)
			if err != nil {
				fmt.Printf("❌ ERROR: Failed to read directory: %v\n", err)
				os.Exit(2)
			}

			var artifacts []models.CommitmentArtifact
			for _, f := range files {
				if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
					continue
				}

				path := filepath.Join(dir, f.Name())
				data, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("⚠️  WARNING: Skipping unreadable file %s: %v\n", f.Name(), err)
					continue
				}

				// Verify individual integrity first
				res, _ := verifier.VerifyArtifact(data)
				if res != nil && res.Valid {
					// We need to parse it again to get the struct, or VerifyArtifact could return it.
					// For simplicity, strict parse here.
					var art models.CommitmentArtifact
					json.Unmarshal(data, &art)
					artifacts = append(artifacts, art)
				} else {
					fmt.Printf("❌ INVALID INDIVIDUAL ARTIFACT: %s\n", f.Name())
					os.Exit(1)
				}
			}

			if len(artifacts) == 0 {
				fmt.Println("⚠️  No valid artifacts found.")
				os.Exit(0)
			}

			// Sort by Timestamp (or chain logic).
			// Since artifacts don't have strictly monotonic IDs, we sort by timestamp to reconstruct sequence.
			// Ideally we follow the linked list from HEAD, but for directory scanning, sorting helps.
			sort.Slice(artifacts, func(i, j int) bool {
				return artifacts[i].Timestamp < artifacts[j].Timestamp
			})

			// Verify Chain
			chainRes := verifier.VerifyChain(artifacts)
			if chainRes.Valid {
				fmt.Printf("✅ CHAIN VALID | Count: %d\n", len(artifacts))
				os.Exit(0)
			} else {
				fmt.Printf("❌ CHAIN BROKEN | Index: %d | Reason: %s\n", chainRes.BrokenIndex, chainRes.BrokenReason)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(fileCmd)
	rootCmd.AddCommand(chainCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
