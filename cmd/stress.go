package cmd

import (
	"fmt"
	"net/http"

	"github.com/IsJordanBraz/go-stress-test/internal/entity"
	"github.com/spf13/cobra"
)

func NewStressCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stress",
		Short: "A brief description of your command",
		Long:  "A brief description of your command",
		RunE:  RunStressTest(),
	}
}

func init() {
	stressCmd := NewStressCommand()
	rootCmd.AddCommand(stressCmd)

	stressCmd.PersistentFlags().StringP("url", "u", "", "URL do serviço a ser testado.")
	stressCmd.PersistentFlags().Int64P("requests", "r", 1, "Número total de requests.")
	stressCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Número de chamadas simultâneas.")
	stressCmd.MarkFlagRequired("url")
}

func reportGeneration(report *entity.Report) {
	report.ValidateTotalTime()
	fmt.Println("Tempo Total: ", report.TotalTime)
	fmt.Printf("Total de %d Requests\n", report.RequestsCount)
	for requests := range report.RequestsStatus {
		fmt.Printf("%s: %d Requests\n", requests, report.RequestsStatus[requests])
	}
}

func RunStressTest() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt64("requests")
		workers, _ := cmd.Flags().GetInt64("concurrency")

		report := entity.NewReport()

		data := make(chan int)

		for i := 0; i < int(workers); i++ {
			go worker(data, url, report)
		}

		for i := 0; i < int(requests); i++ {
			data <- i
		}

		reportGeneration(report)

		return nil
	}
}

func worker(data chan int, url string, report *entity.Report) {
	for range data {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		res.Body.Close()

		report.Mutex.Lock()

		report.RequestsCount++
		mapStatus := make(map[string]int)
		mapStatus[res.Status] = res.StatusCode

		_, found := report.RequestsStatus[res.Status]
		if found {
			report.RequestsStatus[res.Status]++
		} else {
			report.RequestsStatus[res.Status] = 1
		}

		report.Mutex.Unlock()
	}
}
