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
	// Global flags
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "gantral-verify",
		Short: "Offline verification tool for Gantral Commitment Artifacts",
		Long: `A standalone tool for auditors to cryptographically verify 
the integrity and chain-of-custody of Gantral artifacts 
without requiring access to the operational database.`,
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose human-readable output")

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

			// Parse first to get details for verbose output
			var art models.CommitmentArtifact
			_ = json.Unmarshal(data, &art)

			result, err := verifier.VerifyArtifact(data)
			if err != nil {
				fmt.Printf("❌ ERROR: Logic failure: %v\n", err)
				os.Exit(2)
			}

			// Verbose Output
			if verbose {
				printArtifactSummary(art, result.Valid)
			}

			// Final Result
			if result.Valid {
				if !verbose {
					fmt.Printf("✅ VALID | ID: %s | Hash: %s\n", result.ArtifactID, result.CalculatedHash)
				}
				os.Exit(0)
			} else {
				if !verbose {
					fmt.Printf("❌ INVALID | ID: %s | Error: %s\n", result.ArtifactID, result.Error)
				}
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
					var art models.CommitmentArtifact
					_ = json.Unmarshal(data, &art)
					artifacts = append(artifacts, art)

					if verbose {
						fmt.Printf("  [+] Loaded Valid Artifact: %s\n", art.ArtifactID[:8])
					}
				} else {
					fmt.Printf("❌ INVALID INDIVIDUAL ARTIFACT: %s\n", f.Name())
					os.Exit(1)
				}
			}

			if len(artifacts) == 0 {
				fmt.Println("⚠️  No valid artifacts found.")
				os.Exit(0)
			}

			// Sort by Timestamp
			sort.Slice(artifacts, func(i, j int) bool {
				return artifacts[i].Timestamp < artifacts[j].Timestamp
			})

			// Verify Chain
			chainRes := verifier.VerifyChain(artifacts)

			if verbose {
				fmt.Println("\n--- CHAIN VERIFICATION SUMMARY ---")
				fmt.Printf("[✓] Total Blocks: %d\n", len(artifacts))
				fmt.Printf("[✓] Start Time:   %s\n", artifacts[0].Timestamp)
				fmt.Printf("[✓] End Time:     %s\n", artifacts[len(artifacts)-1].Timestamp)
			}

			if chainRes.Valid {
				if verbose {
					fmt.Println("RESULT: ADMISSIBLE")
				}
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

func printArtifactSummary(art models.CommitmentArtifact, valid bool) {
	fmt.Println("\n--- ARTIFACT VERIFICATION SUMMARY ---")

	icon := "[✓]"
	if !valid {
		icon = "[❌]"
	}

	fmt.Printf("%s Artifact ID:  %s\n", icon, art.ArtifactID)
	fmt.Printf("%s Signer:       %s\n", icon, art.HumanActorID)
	fmt.Printf("%s Decision:     %s\n", icon, art.AuthorityState)

	// Truncate hash for display
	hashDisp := art.ContextHash
	if len(hashDisp) > 10 {
		hashDisp = hashDisp[:10] + "..."
	}
	fmt.Printf("%s Evidence:     Verified (Hash: %s)\n", icon, hashDisp)

	if valid {
		fmt.Println("RESULT: ADMISSIBLE")
	} else {
		fmt.Println("RESULT: INADMISSIBLE")
	}
	fmt.Println("-------------------------------------")
}
