generate:
	docker buildx bake generate

build:
	docker buildx bake build

lint:
	docker buildx bake lint

unit:
	docker buildx bake unit
