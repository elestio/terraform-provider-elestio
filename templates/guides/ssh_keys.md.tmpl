---
page_title: "SSH Keys"
---

# SSH Keys with Elestio and Terraform

## Summary

1. [What is a SSH Key?](#what-is-a-ssh-key)
2. [How to generate a valid SSH key](#how-to-generate-a-valid-ssh-key)
3. [How to use SSH keys with Elestio and Terraform](#how-to-use-ssh-keys-with-elestio-in-terraform)

## What is a SSH Key?

SSH keys are used for secure connections across a network. They come in pairs, so you have a **public key** and a **private key**.

The standard ssh2 [file format](http://www.openssh.org/txt/draft-ietf-secsh-publickeyfile-02.txt) looks like this:

```
---- BEGIN SSH2 PUBLIC KEY ----
Comment: "1024-bit RSA, converted from OpenSSH by galb@test1"
AAAAB3NzaC1yc2EAAAABIwAAAIEA1on8gxCGJJWSRT4uOrR13mUaUk0hRf4RzxSZ1zRbYY
Fw8pfGesIFoEuVth4HKyF8k1y4mRUnYHP1XNMNMJl1JcEArC2asV8sHf6zSPVffozZ5TT4
SfsUu/iKy9lUcCfXzwre4WWZSXXcPff+EHtWshahu3WzBdnGxm5Xoi89zcE=
---- END SSH2 PUBLIC KEY ----
```

The purpose of using SSH keys is to simplify access to the remote server. Once you have keys on both the remote service and your local host, the service will only grant access to someone who has the matching private key. This means that you don't have to remember and enter passwords every time you want to access the server.

## How to generate a valid SSH key

On Elestio the above format will not work because Elestio services uses OpenSSH and requires the public key to be in OpenSSH format.
Here is an example of the OpenSSH format using DSA encryption (usually they are all in one line):

```
ssh-dss AAAAB3NzaC1kc3MAAACBAIpAwMFJGHmoQ91HoUGS1WL1GRg2K4hTgxXcJqszIJOrya+8vYY
1YiuazPYkOAOhVaSAofNQ754BKelaIERAWARNCFvf72AtVaa8wwNNveuRF6rEbxLbzzPKk0l6/7K0ZY
GAZIOapipBXoFV+nqS95VXvIgY73RNCWesXCOU2f2NAAAAFQDjCACwCNIwp7Jqc+4rxF7zQGkjoQAAA
IBULNkxCd0Y3z90DAdmhvhQar62QGp4XEl6hM+bLShLkD3MFNGYhELo5MVBd12KKJi+srzp6ohYMLbi
beUEHhvKLV3RnIzFaocCu5JCn2rybJqeW4QrOmN2ofGDZs0wx9LyI8F1vyFMtGv+uWzaI2Uye8Ri5Qq
bnNg/LBRPdZRCxAAAAIBRHttgRQv1+AAYDDduT/GJHeOVugIMENPhTbIp5a/sfXcJi5W4FVZzpLtGy7
Q4we16aGv4Wy4dMdaPHIAJtNeRviw10WZbWZHTJ6x30M2/vxrOSuM/KFKSM5SssVrYmorXG+ATgiO/v
7iBZAZRZXcqsbYMBWYVXEIO/utzkU0HRQ== username@test.com
```

In the example above you will note that the key starts with "ssh-dss". This is because this key was generated using DSA as opposed to RSA. A number of vendors in the SSH arena [have argued](http://the.earth.li/~sgtatham/putty/0.55/htmldoc/Chapter8.html#S8.2.10) that users should not employ DSA encryption because

```
DSA has an intrinsic weakness which makes it very easy to create a signature
which contains enough information to give away the private key! This would 
allow an attacker to pretend to be you for any number of future sessions. 
```

For this reason, **Elestio only accepts RSA public keys**.
A valid SSH public key accepted by Elestio, using RSA encryption, using OpenSSH format, will start with "ssh-rsa".

```
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC0IEndxgICS/gUPkbCRO5tgEuTZOOamLpkIh3vwUD
SI016HMjIFQZzzgF+l2Q90MHxFcPiPP1fKCt4YSp4HOOMA3FsZerxnG/ay73WadY38BpJLsb+hx7STo
7LWfCNdCkYPtlSb3fFKpBBI+q2EG1tKddFRtlSI1+mDPIfzA1m7w== username@test.com
```

### Unix/Linux and Mac OS X
1. Open a terminal window.
2. Enter `ssh-keygen -t rsa` and press enter.
3. Look in your ~/.ssh directory (or where you saved the output). You'll `find id_XXX` (private key) and `id_XXX.pub` (public key).

### Windows
1. Download and use [PuTTYgen](https://www.puttygen.com/)
2. Make sure you choose the RSA2 key format
3. Save the private key and the public key that are generated

### Convert existing SSH2 key to OpenSSH format
If you already have a SSH2 key, you can convert it to OpenSSH format using the command `ssh-keygen -i -f ssh2.pub`

## How to use SSH keys in Terraform

If you want to login to your Elestio services to do any work, you will need to submit your public key in OpenSSH format via the `ssh_public_keys` attribute in your resource.

```tf
resource "elestio_postgresql" "postgres" {
  ssh_public_keys = [
    {
      username = "admin"
      key_data = file("~/.ssh/id_rsa.pub")
    }
  ]
}
```

You can now login to the server using the private key and the resource `cname` (use `terraform show` command to retrieve it):

```sh
ssh -o -i ~/.ssh/id_rsa root@database-u525.vm.elestio.app
```

You can also use Terraform provisioners to perform actions from your configuration:

```tf
resource "elestio_postgresql" "postgres" {
  project_id    = "project_id"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"

  # Execute remote commands on the service
  # https://www.terraform.io/docs/language/resources/provisioners/remote-exec.html
  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      host        = self.ipv4
      private_key = file("~/.ssh/id_rsa") # Specify the matching private key
    }

    inline = [
      "cd /opt/app", # We can now run commands on the service
      "docker exec -it postgres psql -U ${self.database_admin.user} -c 'CREATE DATABASE production;'"
    ]
  }
}
```

There is many provisioners available in Terraform, you can find the full list in the [terraform official documentation](https://www.terraform.io/docs/language/resources/provisioners/index.html).

Here is a configuration that combines two provisioners to replace the default `docker-compose.yml` file with our own and restart the container:

```tf
resource "elestio_postgresql" "postgres" {
  project_id    = "project_id"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
  ssh_public_keys = [
    {
      username = "admin"
      key_data = file("~/.ssh/id_rsa.pub") # Submit the public key
    }
  ]

  # 1. Specify the SSH connection config
  connection {
    type        = "ssh"
    user        = "root"
    host        = self.ipv4
    private_key = file("~/.ssh/id_rsa")
  }

  # 2. Replace the default docker-compose.yml file with our own
  provisioner "file" {
    source      = "${path.module}/files/docker-compose.yml" # The pathfile on your local machine
    destination = "/opt/app/docker-compose.yml"             # The destination path on the service
  }

  # 3. Restart the container
  provisioner "remote-exec" {
    inline = [
      "cd /opt/app",
      "docker-compose down",
      "docker-compose up -d"
    ]
  }
}
```