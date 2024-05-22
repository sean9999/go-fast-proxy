IMAGE=crazyhorsecoding/gcp-goproxy
BUILDER=gcr.io/buildpacks/builder:v1 
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

tidy:
	go mod tidy

ref:
	echo $(REF)

info:
	echo image:$(IMAGE) branch:$(BRANCH) ref:$(REF) builder:$(BUILDER)

build:
	pack build --builder=$(BUILDER) $(IMAGE) -t $(IMAGE):$(BRANCH) -t $(IMAGE):$(REF)

push:
	docker push -a $(IMAGE)

