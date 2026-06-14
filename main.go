package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kamuridesu/readeck-manager/internal/config"
	"github.com/kamuridesu/readeck-manager/internal/manager"
	"github.com/kamuridesu/readeck-manager/internal/utils"
)

func Check[T any](v T, err error) T {
	CheckErr(err)
	return v
}

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	utils.LoadDotenv()
}

func printUsage() {
	fmt.Printf("Usage: %s <command> [options]", os.Args[0])
	fmt.Println("\nCommands:")
	fmt.Println("  tag                  Tag all untagged bookmarks using AI")
	fmt.Println("  clean [options]      Delete all bookmarks currently tagged with 'DELETE'")
	fmt.Println("  broken               List all broken bookmarks")
	fmt.Println("  help                 Show this help message")
	fmt.Println("\nOptions for 'clean':")
	fmt.Println("  --dry-run, -d        Show what would be deleted without performing actual deletion")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := Check(config.Load())
	mg := manager.New(cfg)
	ctx := context.Background()

	command := os.Args[1]

	switch command {
	case "tag":
		fmt.Println("Starting AI tagging process...")
		CheckErr(mg.TagUntaggedLinks(ctx))
		fmt.Println("Tagging complete.")

	case "clean":
		cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
		var dryRun bool
		cleanCmd.BoolVar(&dryRun, "dry-run", false, "Show what would be deleted without performing actual deletion")
		cleanCmd.BoolVar(&dryRun, "d", false, "Show what would be deleted without performing actual deletion (shorthand)")

		CheckErr(cleanCmd.Parse(os.Args[2:]))

		if dryRun {
			fmt.Println("Cleaning up bookmarks marked for deletion (DRY RUN)...")
		} else {
			fmt.Println("Cleaning up bookmarks marked for deletion...")
		}

		CheckErr(mg.DeletedTaggedWithDelete(ctx, dryRun))
		fmt.Println("Cleanup complete.")

	case "broken":
		fmt.Println("Finding broken bookmarks...")
		CheckErr(mg.GetBrokenBookmarks(ctx))
		fmt.Println("Search complete.")

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}
