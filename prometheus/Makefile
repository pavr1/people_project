build:
	docker build -t prometheus:1.0 .

install-helm-snbx:
	make build
	kubectl config set-context --current --namespace=snbx
	helm install prometheus prometheus --values ./prometheus/values-snbx.yaml --namespace snbx
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh snbx

install-helm-eng:
	make build
	kubectl config set-context --current --namespace=eng
	helm install prometheus prometheus --values ./prometheus/values-eng.yaml --namespace eng
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh eng