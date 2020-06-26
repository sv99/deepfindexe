deepfindexe
===========

Detect executable in the archive. Support recursive archive like `zip->rar->xx.exe`.

Use for detect executable in the received mail files.
 
Example using with `exim`.
 
```bash
acl_check_mime:
 
deny message   = We do not accept attachments like: $mime_filename
     condition = ${if match{$mime_filename}{\N\.(bat|js|pif|cd|com|exe|lnk|reg|vbs|jse|msi|ocx|dll|sys|cab)$\N}}
 
deny message = Unwanted file extension ($found_extension)
     demime  = bat:com:lnk:pif:scr:vbs:ade:adep:asd:chm:cmd:cpl:crt:dll:hlp:hta:inf:isp:jse:ocx:pcd:reg:url
```
 
script based version
--------------------
 
From [Exiscan Filename Blocking](https://github.com/Exim/exim/wiki/ExiscanFilenameBlocking)
 
```bash
#free_arqexec contain emails of permited sender. Example: *@gmail.com
drop !senders     = wildlsearch;/etc/exim4/lst/fre_arqexec
    demime         = zip:rar:arj:tar:tgz:gz:bz2
    set acl_m9     = ${run{/etc/exim4/exim.checkpkt.sh ${lc:$found_extension} $message_exim_id}}
    message        = This message contains an unwanted binary Attachment in .${uc:$found_extension} file.
    condition      = ${if eq {$runrc}{0}{false}{true}}
```
 
`exim.checkpkt.sh` - detect archive type by extension and not scan deep level.
 
zip->zip->XX.exe - not scaned!!
 
Detect archive type by magic
----------------------------
  
Based on:

1. `github.com/mholt/archiver` - multi archiver, detect files by extension
2. `github.com/gabriel-vasile/mimetype` - detect file type by magic code.
