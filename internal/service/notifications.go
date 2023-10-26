package service

import (
	"time"

	"github.com/nmluci/gostellar/pkg/dto"
)

func (s *service) NotifyError(msg string, reqID string, jobName string) {
	err := s.goStellar.Notification.Discord.Notify(&dto.DiscordWebhookMeta{
		Username: "Natsumi-chan",
		Embeds: []dto.DiscordEmbeds{
			{
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
					{
						Name:  "Request ID",
						Value: reqID,
					},
					{
						Name:  "Job Name",
						Value: jobName,
					},
				},
			},
		},
	})

	if err != nil {
		s.logger.Warn().Err(err).Send()
	}
}

func (s *service) NotifyInfo(msg string, reqID string, jobName string, versionTag string) {
	err := s.goStellar.Notification.Discord.Notify(&dto.DiscordWebhookMeta{
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
						Value: reqID,
					},
					{
						Name:  "Job Name",
						Value: jobName,
					},
					{
						Name:  "Version Tag",
						Value: versionTag,
					},
				},
			},
		},
	})

	if err != nil {
		s.logger.Warn().Err(err).Send()
	}
}
