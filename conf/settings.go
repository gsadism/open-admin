package conf

// todo conf 目录下的所有配置项请务在用户代码中直接导入使用, 会出现循环依赖问题

const (
	// SUPER_ADMIN_PASSWORD : 超级管理员用户密码, 为空则默认密码为 root
	// 使用 go run main.go -p 明文密码 进行修改
	SUPER_ADMIN_PASSWORD = "/6JHWdHYOxLpRVyY.qXbBrE5FNWxDMZPH/DqxgEMnQ0gP7k6"
	// SECRET_KEY : keep the secret key used in production secret!
	// 使用 go run main.go secret
	SECRET_KEY = "c39f2786e7307b5f4572e3ba65bb686bf3c21e95cb7bfe8006fde965edca5d9f"
	// SECRET_IV : 偏移量
	// 使用 go run main.go iv
	SECRET_IV = "cea46d2303ae06868a14836d2cbf4b57"
)

// 使用 go run main.go rsa 生成新的rsa秘钥
const (
	RSA_PUBLIC_KEY = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx+zhLgD7iYer7+c/G1Mx
0gUj+3tMPwtTOXB1SG1IsH2hoktv/qC6aaVoNSdB3j4i9UbWnxw1QXxXlYG+s0QP
1WaVI9XjqKLYxvA4M9gvq7Yf2mVnlKRqEFa9cDH1wo2YMO4KM5YZvl/MfHbGhzST
hFnHp8fzNwqkhndkIrYb6fInItRfD8gXwAVuc/jjr5Mh2tpeXg7D6aSKsg7dG6CR
oymXHLdvmv2SCFzfN2HyMsshM0BE2SbiPSeh8Ci/f+2WebFf/f2n4TZyB6g6nBWB
WQPfTO1JxaXYjyXu1sY7DwQ8YuB7SNnfBt+mXZmcrQfC8epMP4r7DGfAHWkSKMwf
jQIDAQAB
-----END RSA PUBLIC KEY-----
`
	RSA_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDH7OEuAPuJh6vv
5z8bUzHSBSP7e0w/C1M5cHVIbUiwfaGiS2/+oLpppWg1J0HePiL1RtafHDVBfFeV
gb6zRA/VZpUj1eOootjG8Dgz2C+rth/aZWeUpGoQVr1wMfXCjZgw7gozlhm+X8x8
dsaHNJOEWcenx/M3CqSGd2Qithvp8ici1F8PyBfABW5z+OOvkyHa2l5eDsPppIqy
Dt0boJGjKZcct2+a/ZIIXN83YfIyyyEzQETZJuI9J6HwKL9/7ZZ5sV/9/afhNnIH
qDqcFYFZA99M7UnFpdiPJe7WxjsPBDxi4HtI2d8G36ZdmZytB8Lx6kw/ivsMZ8Ad
aRIozB+NAgMBAAECggEAdiO1YOD1wuxoHsjQks4Tp56gB7BxeuOvV0bsF1N0284x
YF0hnh7OaTV0IczD8JdZQwWMN9OMaNkYTGGyFe7HSHxo7xcqg3FPlidHnuNGMDrc
EctQbRbHEkafhRgQuNnBMFDJuHgfWE714Cqz4BuifFgtUikmwMTKOid9BcIFNrV0
BLTvj2XpgftF5tO40MJqvaERTCn3K7Zi8EqOCXNAXSozL02sbM4gHWzfZUuV9b2C
sCl/dr6t+rsyc2B3DcN0CmMHAj1DQqAZ8dBrns75ggE27AHg9d5QC1myOruCyqFt
XSQEo7RFKgtGhCSH3nN3HVcHH4WDJQUasDsG5FWa4QKBgQDOQTPoNol6A/Z9xMlj
YITzpeyP8IafZ9tp7sA7L1FFpimH+EEvtWvfTJCt8i/saiP4JoyUObY5WTAxNDIW
GTqgJoAyuJJmzXe5AW2ZVcZ5NzkR+1kYBk6bZdvkWP7ndPjCOPKx8NOqdjeoC+8c
9j2VDokuTrff9n5xOeSUGsT9SQKBgQD4JOGOCs5aTzoB7hy5Nt+DuZE/lS3Zak1K
oOj9PSBU0teWjyri8fJcup1MtXdUqBmyC5H/06QpuVd5hbFh1EGOT9BRqoeLFm7b
6hl2faPmeNv57IwXGJeeMMswzw+/EoEWyvmTAiZCiY3EMPNxvWQta48K87xQwHAE
LRK1C0tkJQKBgEkfvspuQ3FZ5FMUqJ2uA5TqV6N0kBYF2TmGZurUzMiXr5PARUQn
2vXRDF5PRgRYYSZCbIbEq0FLyJU5PySsjydmXv8chvMhl3ek9Plu+q3REU+5wA2P
duELTNl5X+4pf2REkm7ZxgFiZjsS+BRZjN3CyIn9aJ1jo9J6xvpnPz3ZAoGBAPV6
fswm8oays77vEOlFWI6WSBNvf8ECgV3fmz3Sd6aSQgZGsXV7Rnjnf+8/4cIU3f6d
pZxAgblpChlE4frWoNepK3tuyNIWqirSBsiKeFbqnys6EEhU4sDyokAxxvEN3Dhx
Kb59tN3D+EPq7RpK8sxgCE0cW+JCzU3ETSpe9wElAoGAH5o54P9ZUa6LU1hwc/7K
htxiN5kbDcvPvneVsC3PrEfqQ3XU/V70kHh4dP4x9o8XcK3UTHTY58thy5kFQqVP
cI1eVUvGr7xouBKoi/wCxKzzjp0CMIqUWLm/oFLTxomZMQofeNGpZv7j4Js4YIzB
WJCVMUia4Y0Q056/AxaxkUI=
-----END RSA PRIVATE KEY-----
`
)
