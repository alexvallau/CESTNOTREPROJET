package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/decoder"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	fmt.Println("Serveur démarré sur : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
	<title>Générateur de Dashboard</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 2rem; }
		.container { max-width: 500px; margin: 0 auto; }
		form { background: #f5f5f5; padding: 2rem; border-radius: 8px; }
		.input-group { margin-bottom: 1rem; }
		label { display: block; margin-bottom: 0.5rem; }
		input { width: 100%; padding: 0.5rem; border: 1px solid #ddd; }
		button { background: #4CAF50; color: white; padding: 0.5rem 1rem; border: none; cursor: pointer; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Configuration du Dashboard</h1>
		<form action="/generate" method="post">
			<div class="input-group">
				<label for="numDevices">Nombre d'appareils:</label>
				<input type="number" id="numDevices" name="numDevices" min="1" required>
			</div>
			<button type="submit">Générer</button>
		</form>
	</div>
</body>
</html>`

	t, _ := template.New("webpage").Parse(tmpl)
	t.Execute(w, nil)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	numDevices, _ := strconv.Atoi(r.FormValue("numDevices"))
	if numDevices < 1 {
		http.Error(w, "Nombre d'appareils invalide", http.StatusBadRequest)
		return
	}

	// Génération du YAML
	yamlContent := generateYAML(numDevices)

	// Déploiement Grafana
	if err := deployDashboard(yamlContent); err != nil {
		http.Error(w, "Erreur de déploiement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `
		<h1>Dashboard généré avec succès!</h1>
		<p>%d appareils configurés.</p>
		<a href="/">Retour</a>
	`, numDevices)
}

func generateYAML(numDevices int) string {
	var yamlBuilder strings.Builder

	yamlBuilder.WriteString(`title: "Dashboard IoT Dynamique"
editable: true
tags: [auto-generated]
auto_refresh: 5s

rows:
  - name: Mesures
    panels:
`)

	for i := 1; i <= numDevices; i++ {
		deviceName := fmt.Sprintf("device_%d", i)
		//handlerName := fmt.Sprintf("{{ handler %d }}", i)

		panel := fmt.Sprintf(`      - timeseries:
          title: "%s"
          height: 400px
          datasource: InfluxDB
          targets:
            - influxdb:
                query: >
                  from(bucket: "iot-platform")
                  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
                  |> filter(fn: (r) => r["_measurement"] == "%s")
                  |> filter(fn: (r) => r["_field"] == "consigne" or r["_field"] == "setpoint" or r["_field"] == "temperature")
                  |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
                  |> yield(name: "mean")
`, deviceName, deviceName)

		yamlBuilder.WriteString(panel)
	}

	// Debug
	_ = os.WriteFile("debug_dashboard.yaml", []byte(yamlBuilder.String()), 0644)

	return yamlBuilder.String()
}

func deployDashboard(yamlContent string) error {
	ctx := context.Background()
	client := grabana.NewClient(
		&http.Client{},
		"http://192.168.141.122:3000",
		grabana.WithAPIToken("eyJrIjoid1U5NjdGUDNIZ0xCdDN3YlNoa05jWWtuTEZQTk01UmEiLCJuIjoiQVBJIEpTT04iLCJpZCI6MX0="),
	)

	dashboard, err := decoder.UnmarshalYAML(bytes.NewBufferString(yamlContent))
	if err != nil {
		return err
	}

	folder, err := client.FindOrCreateFolder(ctx, "Dashboards Dynamiques")
	if err != nil {
		return err
	}

	_, err = client.UpsertDashboard(ctx, folder, dashboard)
	return err
}
