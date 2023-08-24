# Introduce

Cloud-backup can help you backup file to github.

# Require

* [Register a github token](https://github.com/settings/tokens/new)
# How to use

**Use docker compose**

save as `docker-compose.yaml`

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