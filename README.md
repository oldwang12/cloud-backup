# Introduce

Cloud-backup can help you backup file to github.

# Require

* [Register a github token](https://github.com/settings/tokens/new)
# How to use

**Run on docker**

save as `docker-compose.yaml`, modify your args, use command `docker compose up -d`.

```dockerfile
version: '3'
services:
  myservice:
    container_name: cloud-backup
    image: oldwang6/cloud-backup:latest
    command:
    - /root/cloud-backup
    # Register a github token: https://github.com/settings/tokens/new.
    - --token=<Github Token>
    # Your github username.
    - --owner=<username>
    # Your backup repo.
    - --repo=<Your Repo>
    # Your backup repo branch.
    - --branch=<Your Branch>
    # Your backup filename.
    - -f=<Filename>
    volumes:
    - type: bind
      source: /path/backup_file
      # Just support path of /root/filename,you can move your backup_file to /root/.
      target: /root/backup_file
```

**Run on kubernetes**

save as `deployment.yaml`, modify your args, use command `kubectl apply -f deployment.yaml`.

```sh
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloud-backup
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cloud-backup
  template:
    metadata:
      labels:
        app: cloud-backup
    spec:
      containers:
        - name: myservice
          image: oldwang6/cloud-backup:latest
          command:
            - /root/cloud-backup
            - --token=<Github Token>
            - --owner=<username>
            - --repo=<Your Repo>
            - --branch=<Your Branch>
            - --filename=<Filename>
          volumeMounts:
            - name: backup-volume
              mountPath: /root/backup_file
      volumes:
        - name: backup-volume
          hostPath:
            path: /path/backup_file
```