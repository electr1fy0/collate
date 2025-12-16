# collate
A minimal images to pdf converting API.


# Usage
```bash
curl -X POST https://images-to-pdf.ayushpandey.xyz -F "files=@a.jpg" -F "files=@b.webp" -F "files=@c.webp" -F "files=@d.png" --output final.pdf
```

# Notes
- Uses "github.com/pdfcpu/pdfcpu/pkg/api" for the conversion
