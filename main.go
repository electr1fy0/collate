package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func ImagesToPDFHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(50 << 20)
	files := r.MultipartForm.File["files"]

	if len(files) == 0 {
		http.Error(w, "no files", http.StatusBadRequest)
		return
	}

	var readers []io.Reader
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, "failed to open file", http.StatusBadRequest)
			return
		}

		defer f.Close()
		readers = append(readers, f)
	}

	var pdfBuf bytes.Buffer

	err := api.ImportImages(nil, &pdfBuf, readers, nil, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=images.pdf")
	w.Write(pdfBuf.Bytes())
}

func main() {
	r := chi.NewRouter()
	r.Post("/convert", ImagesToPDFHandler)
	http.ListenAndServe(":8080", r)
}
