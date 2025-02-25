# WHAT
this little tool will mimics what yumdownloader does:
https://github.com/rpm-software-management/yum-utils/blob/master/yumdownloader.py

# WHY
i needed a tool to dl rpms at this point in time   
but installing yumdownloader might change a state of the system,  
which is smth not to be desired.

On the other hand installing go in one unzip and export away...

# HOW
everything is pretty simple.  

## build
to build it just do:
```
go mod tidy
go build
```
## use
it only takes one arg, name of the package...
```
./rpmdl openssl
``` 
example assumes you're running it from the dir you built it in.
moving binary to the path is bonus task for the reader....