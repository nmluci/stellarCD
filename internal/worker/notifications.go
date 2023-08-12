package worker

import (
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
			{
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
					{
						Name:  "Request ID",
						Value: params.ReqID,
					},
					{
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
}

func (dw *deploymentWorker) NotifyInfo(cred *dto.DiscordWebhoookCred, params NotifyInfoParams) {
	payload := &dto.DiscordWebhookMeta{
		Username: "Natsumi-chan",
		Embeds: []dto.DiscordEmbeds{
			{
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
				Fields: []dto.DiscordField{
					{
						Name:  "Request ID",
						Value: params.ReqID,
					},
					{
						Name:  "Job Name",
						Value: params.JobName,
					},
					{
						Name:  "Version Tag",
						Value: params.VersionTag,
					},
					{
						Name:  "Commit Message",
						Value: params.CommitMessage,
					},
					{
						Name:  "Commit Author",
						Value: params.CommitAuthor,
					},
					{
						Name:  "Commit Timestamp",
						Value: params.CommitTimestamp,
					},
					{
						Name:  "Build Time",
						Value: params.BuildTime,
					},
				},
				URL: params.CommitURL,
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
}
