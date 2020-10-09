# Go DVPN web embedder

[![pipeline status](https://gitlab.com/mysteriumnetwork/go-dvpn-web/badges/master/pipeline.svg)](https://gitlab.com/mysteriumneam/go-dvpn-web/-/commits/master)

This library generates an embedded version of the [dvpn-web](https://github.com/mysteriumnetwork/dvpn-web).

To re-generate the `assets_vfsdata.go`, either `mage generate` or `go run mage.go generate`

## Local

To bundle webUI locally:

1)  Export path to WebUI root source directory:
    
    ```console
    $ export DVPN_WEB_LOCAL_PATH=/Users/user/dev/src/dvpn-web
    ```

2) Generate local dist.tar.gz by running yarn command in .../dvpn-web

    ```console
    $ yarn local_release
    ```
   
3) Bundle dist

    ```console
    mage GenerateLocal
   ```
   or
   ```console
   go run  mage.go GenerateLocal
   ```
   
4) In node source root replace module

    ```console
   $ go mod edit -replace github.com/mysteriumnetwork/go-dvpn-web=/Users/user/go/src/go-dvpn-web 
   ```
   
5) Finally, build and run node!