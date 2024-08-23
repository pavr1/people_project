build:
	docker build -t auth:1.0 .

install-helm-snbx:
	make build
	kubectl config set-context --current --namespace=snbx
	helm install auth auth --values ./auth/values-snbx.yaml --namespace snbx
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh