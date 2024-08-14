package squeeze

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	squeezegoclient "github.com/dexpro-solutions-gmbh/squeeze-go-client"
	"github.com/spf13/cobra"
)

func newBenchmarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bench",
	}

	flags := cmd.Flags()

	sqzKey := flags.String("sqz-key", "", "Squeeze API key")
	sqzBasePath := flags.String("sqz-base-path", "http://squeeze.docker.localhost/api/v2", "Squeeze API base path")

	dataDir := flags.String("data-dir", "", "Directory containing files to be uploaded")
	delay := flags.Duration("delay", 5*time.Second, "How long should the tool wait before starting the benchmark?")
	pollInterval := flags.Duration("poll-interval", 5*time.Second, "How often should the tool check for completion of the benchmark?")
	timeout := flags.Duration("timeout", 30*time.Minute, "When should the benchmark be considered failed if it hasn't completed yet?")

	cmd.Run = func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		ctx, cancel := context.WithTimeout(ctx, *timeout)
		defer cancel()

		slog.Info("Setting up API client", "basePath", *sqzBasePath)
		time.Sleep(time.Second) // Give users some time to abort if wrong URL was set.

		client := squeezegoclient.NewClient(*sqzBasePath)
		client.ApiKey = *sqzKey

		documentCount := 0

		// Get initial count in validation step so we can later check if everything reached the validation
		// This also ensures the API credentials are valid.
		initialStepCount, err := getStepCount(client.Queue, "Validation")
		if err != nil {
			panic(err)
		}

		slog.Info("Initial Validation step Count", "count", initialStepCount)

		err = filepath.WalkDir(*dataDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			file, err := os.OpenFile(path, os.O_RDONLY, 0)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					slog.Error("Closing file failed", "err", err)
				}
			}(file)

			fileName := filepath.Base(path)

			_, resErr := client.Document.ProcessDocument(1, 0, "", nil, file, fileName)
			if resErr != nil {
				// Ignore the upload error and continue with next file - this allows the benchmark
				// to still be run. That is most often what we want since Squeeze may decide
				// to reject uploads of invalid files / file types.
				slog.Error("Failed to upload document", "file", fileName, "err", resErr)
				return nil
			}

			slog.Info("File uploaded", "file", fileName)

			documentCount += 1

			return nil
		})
		if err != nil {
			panic(err)
		}

		if documentCount == 0 {
			slog.Error("No files uploaded, stopping benchmark")
			return
		}

		slog.Info("All files uploaded", "count", documentCount)

		slog.Info("Waiting you can start the worker to process the document", "duration", *delay)

		time.Sleep(*delay)

		slog.Info("Timer started")
		start := time.Now()

		for {
			// Check if all documents have been processed
			count, err := getStepCount(client.Queue, "Validation")
			if err != nil {
				panic(err)
			}

			doneCount := count - initialStepCount

			if doneCount == documentCount {
				break
			}

			slog.Debug("Validation count", "done", doneCount)

			select {
			case <-ctx.Done():
				panic(ctx.Err())
			case <-time.After(*pollInterval):
				// Continue with next iteration
			}
		}

		elapsed := time.Since(start)
		slog.Info("All documents processed", "elapsed", elapsed)
	}

	return cmd
}

func getStepCount(client *squeezegoclient.QueueApi, stepName string) (int, error) {
	step, err := client.GetQueueStep(stepName)
	if err != nil {
		return 0, fmt.Errorf("failed to get step count: %w", err)
	}

	return step.Count, nil
}
