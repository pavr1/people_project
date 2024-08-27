Get the kubernetes certificate
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | base64 --decode
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.client-certificate-data}' | base64 --decode
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.client-key-data}' | base64 --decode


echo "${{ secrets.K8S_CA_CERT }}" | base64 -d > ca.crt
echo "${{ secrets.K8S_CLIENT_CERT }}" | base64 -d > client.crt
echo "${{ secrets.K8S_CLIENT_KEY }}" | base64 -d > client.key
echo "SECRET *************** ${{ secrets.K8S_CLIENT_CERT }}" 
kubectl config set-cluster docker-desktop --server=https://127.0.0.1:6443 --certificate-authority=ca.crt
kubectl config set-credentials docker-desktop --client-certificate=client.crt --client-key=client.key
kubectl config set-context docker-desktop --cluster=docker-desktop --user=docker-desktop
kubectl config use-context docker-desktop