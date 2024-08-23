build:
	docker build -t person:1.0 .

apply-namespaces:
	kubectl delete namespace snbx
	kubectl create namespace snbx
	kubectl annotate namespace snbx app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=snbx
	kubectl label namespace snbx app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=snbx

install-helm-snbx:
	make build
	kubectl config set-context --current --namespace=snbx
	helm install people charts --values ./charts/values-snbx.yaml --namespace snbx
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh

install-helm-eng:
	helm install helloworld charts/hello_world --values ./charts/hello_world/values-eng.yaml --namespace eng

install-helm-stg:
	helm install helloworld charts/hello_world --values ./charts/hello_world/values-stg.yaml --namespace stg

install-helm-prod:
	helm install helloworld charts/hello_world --values ./charts/hello_world/values-prod.yaml --namespace prod

install-service-account:
	kubectl create clusterrolebinding people-charts-admin --clusterrole=cluster-admin --serviceaccount=snbx:people-charts
	kubectl create rolebinding people-charts-admin --role=admin --serviceaccount=snbx:people-charts --namespace=kubernetes-dashboard

generate-token-snbx:
	kubectl -n snbx create token people-charts

kubernetes-portforward:
	kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443