IMAGE=crazyhorsecoding/gcp-goproxy
BUILDER=gcr.io/buildpacks/builder:v1 
BRANCH := $$(git branch --show-current)
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

tidy:
	go mod tidy

pack:
	pack build --builder=$(BUILDER) $(IMAGE) -t $(IMAGE):$(SEMVER) -t $(IMAGE):$(BRANCH)

push:
	docker push -a crazyhorsecoding/gcp-goproxy
