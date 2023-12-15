package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/nmluci/gostellar/pkg/dto"
)

type NotifyErrorParams struct {
	Message string
	ReqID   string
	JobName string
}

type NotifyInfoParams struct {
	Message         string
	ReqID           string
	JobName         string
	VersionTag      string
	CommitMessage   string
	CommitURL       string
	CommitTimestamp string
	CommitAuthor    string
	// CommitComitter  string
	BuildTime string
}

func (dw *deploymentWorker) NotifyError(cred *dto.DiscordWebhoookCred, params NotifyErrorParams) {
	payload := &dto.DiscordWebhookMeta{
		Username: "Natsumi-chan",
		Embeds: []dto.DiscordEmbeds{
			dto.DiscordEmbeds{
				Title:       "Stellar CI/CD Error Report",
				Description: params.Message,
				Color:       "13421823",
				Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
				Footer: dto.DiscordFooter{
					Text: "Stellar-CD Notification | Stellar-MS",
				},
				Author: dto.DiscordAuther{
					Name: "Stellar-CD by Natsumi-chan",
				},
				Fields: []dto.DiscordField{
					dto.DiscordField{
						Name:  "Request ID",
						Value: params.ReqID,
					},
					dto.DiscordField{
						Name:  "Job Name",
						Value: params.JobName,
					},
				},
			},
		},
	}

	err := dw.goStellar.Notification.Discord.Notify(payload)
	if err != nil {
		dw.logger.Warn().Err(err).Send()
	}

	if cred != nil {
		err := dw.goStellar.Notification.Discord.NotifyWithCred(cred, payload)
		if err != nil {
			dw.logger.Warn().Err(err).Send()
		}
	}

	msg := &dto.NtfyMessage{
		Topic:    "maid",
		Message:  fmt.Sprintf("Deployment failed for job: %s.\nReasons: %s\n\nRequestID: %s", params.JobName, params.Message, params.ReqID),
		Title:    "Deployment Error",
		Tags:     []string{"error"},
		Priority: 5,
	}

	err = dw.goStellar.Notification.Ntfy.Notify(context.Background(), msg)
	if err != nil {
		dw.logger.Warn().Err(err).Send()
	}
}

func (dw *deploymentWorker) NotifyInfo(cred *dto.DiscordWebhoookCred, params NotifyInfoParams) {
	fields := []dto.DiscordField{
		dto.DiscordField{
			Name:  "Request ID",
			Value: params.ReqID,
		},
		dto.DiscordField{
			Name:  "Job Name",
			Value: params.JobName,
		},
		dto.DiscordField{
			Name:  "Version Tag",
			Value: params.VersionTag,
		},
		dto.DiscordField{
			Name:  "Build Time",
			Value: params.BuildTime,
		},
	}

	if params.CommitMessage != "" {
		fields = append(fields, dto.DiscordField{
			Name:  "Commit Message",
			Value: params.CommitMessage,
		})
	}

	if params.CommitAuthor != "" {
		fields = append(fields, dto.DiscordField{
			Name:  "Commit Author",
			Value: params.CommitAuthor,
		})
	}

	if params.CommitTimestamp != "" {
		fields = append(fields, dto.DiscordField{
			Name:  "Commit Timestamp",
			Value: params.CommitTimestamp,
		})
	}

	payload := &dto.DiscordWebhookMeta{
		Username: "Natsumi-chan",
		Embeds: []dto.DiscordEmbeds{
			dto.DiscordEmbeds{
				Title:       "Stellar CI/CD Info Report",
				Description: "Deployment Success",
				Color:       "13421823",
				Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
				Footer: dto.DiscordFooter{
					Text: "Stellar-CD Notification | Stellar-MS",
				},
				Author: dto.DiscordAuther{
					Name: "Stellar-CD by Natsumi-chan",
				},
				Fields: fields,
				URL:    params.CommitURL,
			},
		},
	}

	err := dw.goStellar.Notification.Discord.Notify(payload)
	if err != nil {
		dw.logger.Warn().Err(err).Send()
	}

	if cred != nil {
		err = dw.goStellar.Notification.Discord.NotifyWithCred(cred, payload)
		if err != nil {
			dw.logger.Warn().Err(err).Send()
		}
	}

	msg := &dto.NtfyMessage{
		Topic: "maid",
		Message: fmt.Sprintf("Deployment success for job: %s.\nBuild Time: %s\nAuthor Name: %s\nCommit Message: %s\nVersion Tag: %s\n\nRequestID: %s",
			params.JobName,
			params.BuildTime,
			params.CommitAuthor,
			params.CommitMessage,
			params.VersionTag,
			params.ReqID,
		),
		Title:    "Deployment Success",
		Tags:     []string{"success"},
		Priority: 1,
	}

	err = dw.goStellar.Notification.Ntfy.Notify(context.Background(), msg)
	if err != nil {
		dw.logger.Warn().Err(err).Send()
	}
}
