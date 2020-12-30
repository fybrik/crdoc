all: clean
	mkdir -p output/
	go run main.go --strict false --output-dir output/ --config-dir config/ --file config/example.yaml

clean:
	rm -rf output/