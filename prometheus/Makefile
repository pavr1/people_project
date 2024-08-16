build:
	docker build -t prometheus:1.0 .

install-helm-snbx:
	make build
	kubectl config set-context --current --namespace=snbx
	helm install prometheus prometheus --values ./prometheus/values-snbx.yaml --namespace snbx