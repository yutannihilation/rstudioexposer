# rstudioexposer

(This won't work anymore...)

Expose RStudio Server to public.

## Usage

1. Launch RStudio with rstudioexposer

```sh
docker run -d -p 8787:8787 -p 80:80 yutannihilation/tidyverse-open
```

2. Access to http://localhost/ (not 8787) on web browsers.

3. You can log in to RStudio Server automatically!
