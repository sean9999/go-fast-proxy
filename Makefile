IMAGE=crazyhorsecoding/gcp-goproxy
BUILDER=gcr.io/buildpacks/builder:v1 
BRANCH := $$(git branch --show-current)
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
REF := $$(git describe --dirty --tags --always)

tidy:
	go mod tidy

ref:
	echo $(REF)

info:
	echo image:$(IMAGE)	semver:$(SEMVER) branch:$(BRANCH) ref:$(REF) builder:$(BUILDER)

build:
	pack build --builder=$(BUILDER) $(IMAGE) -t $(IMAGE):$(SEMVER) -t $(IMAGE):$(BRANCH) -t $(IMAGE):$(REF)

push:
	docker push -a $(IMAGE)

