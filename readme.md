# Generate code gen

## Install package
```
go get -u github.com/smallnest/gen
go install github.com/smallnest/gen@latest 
```


## Code gen
```
gen --sqltype=sqlite3 \
        --connstr "../backend/db.sqlite3" \ 
        --database main  \
        --json \
        --gorm \
        --guregu \
        --rest \
        --out ./example \
        --module example.com/rest/example \
        --mod \
        --server \
        --makefile \
        --json-fmt=snake \
        --generate-dao \
        --generate-proj \
        --overwrite
```