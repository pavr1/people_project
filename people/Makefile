build:
	docker build -t pavr1/person:1.0 .
	docker push pavr1/person:latest

apply-namespaces:
	kubectl create namespace snbx
	kubectl annotate namespace snbx app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=snbx
	kubectl label namespace snbx app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=snbx
	kubectl create namespace eng
	kubectl annotate namespace eng app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=eng
	kubectl label namespace eng app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=eng
	kubectl create namespace stg
	kubectl annotate namespace stg app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=stg
	kubectl label namespace stg app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=stg
	kubectl create namespace prod
	kubectl annotate namespace prod app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=prod
	kubectl label namespace prod app.kubernetes.io/managed-by=Helm meta.helm.sh/release-name=people meta.helm.sh/release-namespace=prod

install-helm-snbx:
	make build
	kubectl config set-context --current --namespace=snbx
	helm install people charts --values ./charts/values-snbx.yaml --namespace snbx
	kubectl apply -f ./charts/3rd-parties/mongo/mongo-snbx.yaml
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh snbx

uninstall-helm-snbx:
	helm uninstall people -n snbx

install-helm-eng:
	make build
	kubectl config set-context --current --namespace=eng
	helm install people charts --values ./charts/values-eng.yaml --namespace eng
	kubectl apply -f ./charts/3rd-parties/mongo/mongo-eng.yaml
	chmod +x ./scripts/vars.sh
	sh ./scripts/vars.sh eng

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