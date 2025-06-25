package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"os"
)

var serverPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		handler := func(ctx *fasthttp.RequestCtx) {
			uuid := uuid.New().String()
			switch string(ctx.Path()) {
			case "/version":
				log.Info().
					Str("request_id", uuid).
					Str("method", string(ctx.Method())).
					Str("path", string(ctx.Path())).
					Str("remote_addr", ctx.RemoteAddr().String()).
					Msg("Check version request")
				ctx.Response.Header.SetContentType("application/json")
				ctx.Response.Header.Set("X-Request-ID", uuid)
				_, err := fmt.Fprintf(ctx, `{"version": "%s", "commit": "%s", "date": "%s", "requestID": "%s"}`, Version, Commit, BuildDate, uuid)
				if err != nil {
					return
				}
			default:
				log.Info().
					Str("request_id", uuid).
					Str("method", string(ctx.Method())).
					Str("path", string(ctx.Path())).
					Str("remote_addr", ctx.RemoteAddr().String()).
					Msg("Incoming request")
				ctx.Response.Header.Set("X-Request-ID", uuid)
				_, err := fmt.Fprintf(ctx, `{"message:" "FastHTTP welcomes you, traveller!", "requestID": "%s"}`, uuid)
				if err != nil {
					return
				}
			}
		}
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s port", addr)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
}
