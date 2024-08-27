build:
	docker build -t pavr1/auth:latest .
	docker push pavr1/people_project/auth:latest

install-helm-snbx:
	kubectl config set-context --current --namespace=snbx
	helm install auth auth --values ./auth/values-snbx.yaml --namespace snbx
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh snbx

uninstall-helm-snbx:
	helm uninstall auth -n snbx

install-helm-eng:
	make build
	kubectl config set-context --current --namespace=eng
	helm install auth auth --values ./auth/values-eng.yaml --namespace eng
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh eng