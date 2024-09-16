// Package stream contains SSE stream handling logic.
package stream

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/r3labs/sse/v2"
	"github.com/szkiba/xk6-top/internal/digest"
)

// Subscribe subscribes to SSE channel and returns a tea command for message event generation.
func Subscribe(ctx context.Context, url string, sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		client := sse.NewClient(url)

		client.OnConnect(func(_ *sse.Client) {
			sub <- &digest.Event{Type: digest.EventTypeConnect}
		})

		client.OnDisconnect(func(_ *sse.Client) {
			sub <- &digest.Event{Type: digest.EventTypeDisconnect}
		})

		parser := newParser()

		return client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
			event, perr := parser.parse(msg)
			if perr != nil {
				sub <- perr

				return
			}

			sub <- event
		})
	}
}
