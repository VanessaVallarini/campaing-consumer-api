build:
	docker-compose up
	rm -Rf dist
	cp -Rf src dist
	find dist -type f ! -name "*.png" -delete
	find src -type f -name "*.png" -delete