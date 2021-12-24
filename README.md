# Go DVPN web embedder

[![pipeline status](https://gitlab.com/mysteriumnetwork/go-dvpn-web/badges/master/pipeline.svg)](https://gitlab.com/mysteriumneam/go-dvpn-web/-/commits/master)

This library generates an embedded version of the [dvpn-web](https://github.com/mysteriumnetwork/dvpn-web).

To re-generate the `assets_vfsdata.go`, either `mage generate` or `go run mage.go generate`

## Local

To bundle webUI locally:

1)  Export version you would like to bundle (makes sure it is available in [dpv-web releases](https://github.com/mysteriumnetwork/dvpn-web/releases))
    
    ```console
    $ export GIT_TAG_VERSION=1.2.1
    ```

2) Bundle dist

    ```console
   $ mage Generate
   ```
   or
   ```console
   $ go run mage.go Generate
   ```
   
3) In node source root replace module

    ```console
   $ go mod edit -replace github.com/mysteriumnetwork/go-dvpn-web=/Users/user/go/src/go-dvpn-web 
   ```
   
4) Finally, build and run node!