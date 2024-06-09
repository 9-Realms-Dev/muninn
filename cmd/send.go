package cmd

import (
	"fmt"
	"strings"

	"github.com/9-Realms-Dev/muninn-client/internal/util"
	munnin "github.com/9-Realms-Dev/muninn-core"
	"github.com/9-Realms-Dev/muninn-core/formats"
	"github.com/spf13/cobra"
)

var (
	filePath string
	req      string
	sendCmd  = &cobra.Command{
		Use:   "send",
		Short: "send specified http request",
		Long:  "",
		RunE:  sendCommand,
	}
)

func init() {
	sendCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "path/to/file", "http file location")
	sendCmd.PersistentFlags().StringVarP(&req, "request", "r", "test endpoing 1", "this would be the name following ### in httop file")

	rootCmd.AddCommand(sendCmd)
}

func sendCommand(cmd *cobra.Command, args []string) error {
	util.Logger.Debug("sending request")
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	requests, err := munnin.ReadHttpFile(file)
	if err != nil {
		return err
	}

	request, err := cmd.Flags().GetString("request")
	if err != nil {
		return err
	}

	if request != "" {
		for _, req := range requests {
			if strings.EqualFold(request, req.Title) {
				response, err := munnin.SendHttpRequest(req)
				if err != nil {
					return err
				}

				util.Logger.Debug(response)

				json, err := formats.FormatJSONResponse(response.Response)
				if err != nil {
					util.Logger.Fatal(err.Error())
				}

				fmt.Println(json.CliRender())
				return nil
			}
		}
	} else {
		responses, err := munnin.SendHttpRequests(requests)
		if err != nil {
			return err
		}

		util.Logger.Info(responses)
	}

	return nil
}
