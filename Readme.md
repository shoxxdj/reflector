# Reflector 

A binary to test parameter reflection in urls. 

## Setup 

```
go install github.com/shoxxdj/reflector@latest
```

## Usage

```
            __ _           _             
           / _| |         | |            
  _ __ ___| |_| | ___  ___| |_ ___  _ __ 
 | '__/ _ \  _| |/ _ \/ __| __/ _ \| '__|
 | | |  __/ | | |  __/ (__| || (_) | |   
 |_|  \___|_| |_|\___|\___|\__\___/|_|                                            
                                         
Reflector : a binary to find xss based on url params reflections. v:0.1
        -c: Color mode
        -m: Mode of fuzzing
                sniper= one param a time
                bram= all parameters same time (used to limit requests)
                anything else : Same pattern in all the parameters
        -u: URL to test
        -v: Verbose mode
```
