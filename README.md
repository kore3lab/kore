# Kore

## KoreCtl

```
$ korectl version
$ korectl --help
```

```
$ korectl profile list                  # profile 목록 조회
$ korectl profile dump <PROFILE-NAME>   # profile 덤프
```

```
$ korectl manifest generate                 # manifest 프린트  (default)
$ korectl manifest generate <PROFILE-NAME>
$ korectl manifest generate dashboard       # dashboard manifest 프린트
$ korectl manifest generate dashboard-aio   # dashboard (all-in-one) manifest 프린트
```

```
$ korectl install                 # 설치 (default)
$ korectl install <PROFILE-NAME>
$ korectl install dashboard       # dashboard 설치
$ korectl install dashboard-aio   # dashboard (all-in-one) 설치
```

```
$ korectl uninstall                 # 삭제 (default)
$ korectl uninstall <PROFILE-NAME> 
$ korectl uninstall dashboard       # dashboard 삭제
$ korectl uninstall dashboard-aio   # dashboard (all-in-one) 삭제
```


```
$ korectl operator init     # kore-operator 설치
$ korectl operator remove   # kore-operator 제거
```


```
$ korectl manifest generate dashboard --set metrics-server.enabled=true   # TODO
```

## Development

* KoreOperator

```
$ kubectl  apply -f manifests/charts/base/crds/install.kore3lab.io_koreoperators.yaml


$ go run operator/main.go

$ kubectl apply -f -<<EOF
apiVersion: install.kore3lab.io/v1alpha1
kind: KoreOperator 
metadata:
  name: default
  namespace: kore-system
spec:
  components:
    base:
      enabled: false
    dashboard:
      enabled: true
EOF

$ kubectl apply -f manifests/profiles/default.yaml
$ kubectl delete -f manifests/profiles/default.yaml
```




* initialize project

```
operator-sdk init --domain kore3lab.io --repo github.com/kore3lab/operator
go mod tidy
operator-sdk create api --group install --version v1alpha1 --kind KoreOperator --resource --controller
make install
```
