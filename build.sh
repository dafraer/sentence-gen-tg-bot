#Pass version using command line argument
command sudo docker build --no-cache  --platform linux/amd64 -t dafraer/sen-gen-bot:$1 .
command  docker push dafraer/sen-gen-bot:$1