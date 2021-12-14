build: md
	zip -vr application.zip application.go public/

md:
	pandoc public/index.md  -s --css="https://cdn.jsdelivr.net/gh/kognise/water.css@latest/dist/light.min.css" -f markdown -t html  -o public/index.html

serve: md
	go run application.go

.PHONY: build md serve
