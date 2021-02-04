deepfindexe
===========

Detect executable in the mail attachments.

Have scan recursive archive like `zip->rar->xx.exe`.
Detect executable by extension and by mime.

It support passworded zip and rar.

```bash
./bin/Darwin/deepfindexe -h
Usage:
  deepfindexe [OPTIONS] file

Application Options:
  -v, --verbose  Show verbose debug information
  -e, --ext=     Executable extensions (default:
                 ade|adp|asd|bas|bat|cab|chm|cmd|com|cpl|crt|dll|exe|hlp|hta|inf|ins|isp|jse|jar|lib|lnk|mdb|mde|mdz|msc|msi|msp|mst|ole|ocx|p-
                 cd|pif|reg|scr|sct|shs|shb|sys|url|vbe|vbs|vxd|wsc|wsf|wsh)

Help Options:
  -h, --help     Show this help message

```

Example using with `exim`.
 
```bash

# check mime
acl_smtp_mime = acl_check_mime
...

begin acl

acl_check_mime:
  
    # exim_exe_sender_whitelist contain emails of permited sender. Example: *@gmail.com 
    deny 	!senders 	= wildlsearch;/etc/exim4/exim_exe_sender_whitelist
  	        decode 		= default
  	        set acl_m0  = ${run{/etc/exim4/scan/deepfindexe $mime_decoded_filename}}
  	        condition   = ${if eq{$runrc}{0}{false}{true}}
  	        message 	= "Executable not allowed ${mime_filename}: ${value}."

    accept
```

Used packages
-------------
  
1. `github.com/mholt/archiver` - multi archiver, detect files by extension
2. `github.com/gabriel-vasile/mimetype` - detect file type by magic code.

From `archiver` use waker and decoder, from `mimetype` use archives and executable detection.

References
----------

1. [Exiscan Filename Blocking](https://github.com/Exim/exim/wiki/ExiscanFilenameBlocking)
2. [Запрет передачи *.exe файлов через Exim](https://forum.lissyara.su/mta-mail-transfer-agent-f20/zapret-peredachi-exe-fajlov-cherez-exim-t3360.html)
3. [Exim - фильтрация писем с вложениями в zip](https://forum.lissyara.su/mta-mail-transfer-agent-f20/exim-fil-traciya-pisem-s-vlojeniyami-v-zip-t43423.html)
