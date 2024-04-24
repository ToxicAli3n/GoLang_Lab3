package Lang

import (
	Painter "github.com/roman-mazur/architecture-lab-3/painter"
	"io"
	"log"
	"net/http"
	"strings"
)

// HttpHandler конструює обробник HTTP запитів, який дані з запиту віддає у Parser, а потім відправляє отриманий список
// операцій у Painter.Loop.
func HttpHandler(loop *Painter.Loop, p *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			in = strings.NewReader(r.URL.Query().Get("Cmd"))
		}

		cmds, err := p.Parse(in)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, cmd := range cmds {
			loop.Post(cmd)
		}
		rw.WriteHeader(http.StatusOK)
	})
}
