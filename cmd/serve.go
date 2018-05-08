package cmd

import (
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twhiston/factd/lib/common/logging"
	"github.com/twhiston/factd/lib/common/metrics"
	factd2 "github.com/twhiston/factd/lib/factd"
	"github.com/twhiston/factd/lib/formatter"
	"net/http"
	"net/http/pprof"
)

func getFactdRoute(f *factd2.Factd) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if !viper.GetBool("async") {
			f.Collect()
		}
		data, err := f.Format()
		logging.HandleError(err)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data.Bytes())
		logging.HandleError(err)
	}
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve fact data to an endpoint",
	Long: `Serve will send the data to an endpoint as json.
By default the endpoint is root, and the server is on port 8080
There is no authentication provided, it is strongly advised that you secure all routes.`,
	Run: func(cmd *cobra.Command, args []string) {

		factd := setupFactD()
		//Always use json formatting
		conf := factd.GetConfig()
		conf.Formatter = new(formatter.JSONFormatter)
		factd.SetConfig(conf)

		//Start the reporters
		if viper.GetBool("async") {
			factd.Collect()
			factd.RunReporters()
		}
		router := httprouter.New()
		router.GET("/", getFactdRoute(factd))

		router.GET("/healthz", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("ok"))
			logging.HandleError(err)
		})

		//Prometheus Setup
		if viper.GetBool("prometheus") {
			prometheus.MustRegister(metrics.PromErrorCount)
			prometheus.MustRegister(metrics.PromEnabledPlugins)
			metrics.PromEnabledPlugins.Set(float64(len(factd.GetConfig().Plugins)))
			router.Handler("GET", "/metrics", promhttp.Handler())
		}

		//Pprof
		if viper.GetBool("pprof") {
			router.HandlerFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
			router.Handler(http.MethodGet, "/debug/pprof/:item", http.DefaultServeMux)
		}

		err := http.ListenAndServe(":8080", router)
		factd.StopReporters()
		if err != nil {
			logging.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().Bool("async", true, "get fact data asynchronously instead of when hitting endpoint")
	serveCmd.PersistentFlags().Bool("prometheus", true, "if true publish a prometheus endpoint")
	serveCmd.PersistentFlags().Bool("pprof", true, "if true publish a pprof endpoint")

	logging.Fatal(viper.BindPFlag("async", serveCmd.PersistentFlags().Lookup("async")))
	logging.Fatal(viper.BindPFlag("prometheus", serveCmd.PersistentFlags().Lookup("prometheus")))
	logging.Fatal(viper.BindPFlag("pprof", serveCmd.PersistentFlags().Lookup("pprof")))
}
