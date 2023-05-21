package worker

import (
	"time"

	"github.com/nmluci/gostellar/pkg/dto"
)

func (dw *deploymentWorker) NotifyError(cred *dto.DiscordWebhoookCred, msg string, reqID string, jobName string) {
	payload := &dto.DiscordWebhookMeta{
		Username: "Natsumi-chan",
		Embeds: []dto.DiscordEmbeds{
			dto.DiscordEmbeds{
				Title:       "Stellar CI/CD Error Report",
				Description: msg,
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
						Value: reqID,
					},
					dto.DiscordField{
						Name:  "Job Name",
						Value: jobName,
					},
				},
			},
		},
	}

	err := dw.goStellar.Notification.Discord.Notify(payload)
	if err != nil {
		dw.logger.Warnf("Notify err: %+v", err)
	}

	if cred != nil {
		err := dw.goStellar.Notification.Discord.NotifyWithCred(cred, payload)
		if err != nil {
			dw.logger.Warnf("Notify err: %+v", err)
		}
	}
}

func (dw *deploymentWorker) NotifyInfo(cred *dto.DiscordWebhoookCred, msg string, reqID string, jobName string, versionTag string) {
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
				Fields: []dto.DiscordField{
					dto.DiscordField{
						Name:  "Request ID",
						Value: reqID,
					},
					dto.DiscordField{
						Name:  "Job Name",
						Value: jobName,
					},
					dto.DiscordField{
						Name:  "Version Tag",
						Value: versionTag,
					},
				},
			},
		},
	}

	err := dw.goStellar.Notification.Discord.Notify(payload)
	if err != nil {
		dw.logger.Warnf("Notify err: %+v", err)
	}

	if cred != nil {
		err = dw.goStellar.Notification.Discord.NotifyWithCred(cred, payload)
		if err != nil {
			dw.logger.Warnf("Notify err: %+v", err)
		}
	}
}
