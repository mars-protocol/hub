<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>{{ .Title }}</title>
    <style>
      body {
        margin: 0;
      }
      .swagger-ui .info {
        margin: 35px 0 !important;
      }
    </style>
    <link rel="stylesheet" type="text/css" href="//unpkg.com/swagger-ui-dist@4.1.3/swagger-ui.css" />
    <link rel="icon" type="image/png" href="//unpkg.com/swagger-ui-dist@4.1.3/favicon-32x32.png" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="//unpkg.com/swagger-ui-dist@4.1.3/swagger-ui-bundle.js"></script>
    <script src="//unpkg.com/swagger-ui-dist@4.1.3/swagger-ui-standalone-preset.js"></script>
    <script>
      window.onload = function() {
        window.ui = SwaggerUIBundle({
          url: {{ .URL }},
          dom_id: "#swagger-ui",
          deepLinking: true,
          presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
          ],
          plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
          ],
          layout: "StandaloneLayout"
        });
      }
    </script>
  </body>
</html>
