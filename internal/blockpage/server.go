package blockpage

import (
	"fmt"
	"html/template"
	"net/http"
)

type Server struct {
	Port int
}

var pageTpl = template.Must(template.New("blocked").Parse(`
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Blocked</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background: #111827;
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      margin: 0;
    }
    .box {
      max-width: 680px;
      padding: 32px;
      border-radius: 16px;
      background: #1f2937;
      box-shadow: 0 10px 30px rgba(0,0,0,0.3);
    }
    h1 { margin-top: 0; }
    .host {
      color: #93c5fd;
      font-weight: bold;
    }
    p {
      line-height: 1.6;
      color: #d1d5db;
    }
  </style>
</head>
<body>
  <div class="box">
    <h1>Website Blocked</h1>
    <p>Access to <span class="host">{{.Host}}</span> is blocked on this Mac.</p>
    <p>Take a small pause. Leave this page and return to your planned work.</p>
  </div>
</body>
</html>
`))

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if host == "" {
			host = "this site"
		}
		_ = pageTpl.Execute(w, map[string]string{
			"Host": host,
		})
	})

	addr := fmt.Sprintf("127.0.0.1:%d", s.Port)
	return http.ListenAndServe(addr, mux)
}
