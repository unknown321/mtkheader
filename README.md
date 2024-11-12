mtkheader
=========

Patches MediaTek `header` file with size of `content`.

```
Usage: ./mtkheader-linux-amd64 -header [FILE] -content [FILE]
   or: ./mtkheader-linux-amd64 [FILE]

Patches MediaTek `header` file with size of `content`.
Print header info by providing path to header file without any flags

  -content string
    	content file
  -header string
    	mtk header
```

Example:

```shell
$ ./mtkheader-linux-amd64 test/initrd_header
File:   test/initrd_header
Type:   ROOTFS
Length: 3340176

$ echo -n "length: 9" > content
$ ./mtkheader-linux-amd64 -header test/initrd_header -content content
$ ./mtkheader-linux-amd64 test/initrd_header
File:   test/initrd_header
Type:   ROOTFS
Length: 9
```