package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "error parsing multiform", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	if len(files) == 0 {
		http.Error(w, "no files", http.StatusBadRequest)
		return
	}

	readers := make([]io.Reader, 0, len(files))
	openFiles := make([]io.Closer, 0, len(files))

	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, "failed to open file", http.StatusBadRequest)
			return
		}

		openFiles = append(openFiles, f)
		readers = append(readers, f)
	}

	defer func() {
		for _, f := range openFiles {
			f.Close()
		}
	}()

	var pdfBuf bytes.Buffer
	err := api.ImportImages(nil, &pdfBuf, readers, nil, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=images.pdf")
	w.Write(pdfBuf.Bytes())
}
